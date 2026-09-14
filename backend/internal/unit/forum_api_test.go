package unit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"opencw/internal/configs"
	"opencw/internal/databases"
	"opencw/internal/models"
	"opencw/internal/server"
	"opencw/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type forumAPITest struct {
	db     *gorm.DB
	router *gin.Engine
	prefix string
}

func TestForumAPI(t *testing.T) {
	if os.Getenv("FORUM_API_TESTS") != "1" {
		t.Skip("set FORUM_API_TESTS=1 to run against PostgreSQL")
	}

	_ = godotenv.Load(".env")
	configureForumTestConfig(t)
	db := openForumTestDB(t)
	if err := db.AutoMigrate(
		&models.User{}, &models.EmailOTP{}, &models.RefreshToken{},
		&models.CWSettings{}, &models.PageSettings{}, &models.Progress{},
		&models.ForumCategory{}, &models.ForumPost{}, &models.ForumThread{},
	); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}

	previousDB := databases.DB
	databases.DB = db
	t.Cleanup(func() {
		databases.DB = previousDB
		sqlDB, err := db.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	})

	gin.SetMode(gin.TestMode)
	router := gin.New()
	server.RouterV1Setup(router)

	t.Run("public forum reads", func(t *testing.T) {
		test := forumAPITest{db: db, router: router, prefix: "forum-test-" + uuid.NewString()}
		category, thread, post := test.seedForum(t)

		response := test.request(t, http.MethodGet, "/v1/forum/categories", "", nil)
		assertStatus(t, response, http.StatusOK)
		assertJSONListContainsID(t, response, category.ID, "data")

		response = test.request(t, http.MethodGet, "/v1/forum/categories/"+category.ID.String()+"/threads?page=1&limit=1", "", nil)
		assertStatus(t, response, http.StatusOK)
		var threads struct {
			Data  []models.ForumThread `json:"data"`
			Page  int                  `json:"page"`
			Limit int                  `json:"limit"`
			Total int                  `json:"total"`
		}
		decodeJSON(t, response, &threads)
		if threads.Page != 1 || threads.Limit != 1 || threads.Total != 1 || len(threads.Data) != 1 || threads.Data[0].ID != thread.ID {
			t.Fatalf("unexpected thread listing: %+v", threads)
		}

		response = test.request(t, http.MethodGet, "/v1/forum/threads/"+thread.ID.String(), "", nil)
		assertStatus(t, response, http.StatusOK)
		assertJSONID(t, response, thread.ID, "data")

		response = test.request(t, http.MethodGet, "/v1/forum/threads/"+thread.ID.String()+"/posts", "", nil)
		assertStatus(t, response, http.StatusOK)
		assertJSONListContainsID(t, response, post.ID, "data")

		response = test.request(t, http.MethodGet, "/v1/forum/categories/not-a-uuid/threads", "", nil)
		assertStatus(t, response, http.StatusBadRequest)

		response = test.request(t, http.MethodGet, "/v1/forum/threads/00000000-0000-0000-0000-000000000000", "", nil)
		assertStatus(t, response, http.StatusNotFound)
	})

	t.Run("cursor pagination", func(t *testing.T) {
		test := forumAPITest{db: db, router: router, prefix: "forum-cursor-" + uuid.NewString()}
		category, _, _ := test.seedForum(t)
		baseTime := time.Now().UTC().Add(-time.Hour)
		for index := 0; index < 3; index++ {
			thread := models.ForumThread{
				CategoryID: category.ID,
				AuthorID:   uuid.New(),
				Title:      fmt.Sprintf("%s-thread-%d", test.prefix, index),
				UpdatedAt:  baseTime.Add(time.Duration(index) * time.Minute),
			}
			if err := test.db.Create(&thread).Error; err != nil {
				t.Fatalf("create cursor thread: %v", err)
			}
			t.Cleanup(func() { test.db.Unscoped().Delete(&models.ForumThread{}, "id = ?", thread.ID) })
		}

		response := test.request(t, http.MethodGet, "/v1/forum/categories/"+category.ID.String()+"/threads?cursor=first&limit=2", "", nil)
		assertStatus(t, response, http.StatusOK)
		var firstPage struct {
			Data       []models.ForumThread `json:"data"`
			HasMore    bool                 `json:"has_more"`
			NextCursor string               `json:"next_cursor"`
		}
		decodeJSON(t, response, &firstPage)
		if len(firstPage.Data) != 2 || !firstPage.HasMore || firstPage.NextCursor == "" {
			t.Fatalf("unexpected first thread cursor page: %+v", firstPage)
		}

		response = test.request(t, http.MethodGet, "/v1/forum/categories/"+category.ID.String()+"/threads?cursor="+firstPage.NextCursor+"&limit=2", "", nil)
		assertStatus(t, response, http.StatusOK)
		var secondPage struct {
			Data       []models.ForumThread `json:"data"`
			HasMore    bool                 `json:"has_more"`
			NextCursor string               `json:"next_cursor"`
		}
		decodeJSON(t, response, &secondPage)
		if len(secondPage.Data) != 2 || secondPage.HasMore || secondPage.NextCursor != "" {
			t.Fatalf("unexpected terminal thread cursor page: %+v", secondPage)
		}
		if firstPage.Data[0].ID == secondPage.Data[0].ID || firstPage.Data[1].ID == secondPage.Data[1].ID {
			t.Fatal("cursor pages contain duplicate records")
		}

		response = test.request(t, http.MethodGet, "/v1/forum/categories/"+category.ID.String()+"/threads?cursor=not-a-cursor", "", nil)
		assertStatus(t, response, http.StatusBadRequest)
		response = test.request(t, http.MethodGet, "/v1/forum/categories/"+category.ID.String()+"/threads?cursor=first&page=1", "", nil)
		assertStatus(t, response, http.StatusBadRequest)

		thread := models.ForumThread{CategoryID: category.ID, AuthorID: uuid.New(), Title: test.prefix + "post-thread"}
		if err := test.db.Create(&thread).Error; err != nil {
			t.Fatalf("create post cursor thread: %v", err)
		}
		t.Cleanup(func() {
			test.db.Unscoped().Delete(&models.ForumPost{}, "thread_id = ?", thread.ID)
			test.db.Unscoped().Delete(&models.ForumThread{}, "id = ?", thread.ID)
		})
		for index := 0; index < 3; index++ {
			post := models.ForumPost{ThreadID: thread.ID, AuthorID: uuid.New(), Body: fmt.Sprintf("%s-post-%d", test.prefix, index), CreatedAt: baseTime.Add(time.Duration(index) * time.Minute)}
			if err := test.db.Create(&post).Error; err != nil {
				t.Fatalf("create cursor post: %v", err)
			}
		}
		response = test.request(t, http.MethodGet, "/v1/forum/threads/"+thread.ID.String()+"/posts?cursor=first&limit=2", "", nil)
		assertStatus(t, response, http.StatusOK)
		var postPage struct {
			Data       []models.ForumPost `json:"data"`
			HasMore    bool               `json:"has_more"`
			NextCursor string             `json:"next_cursor"`
		}
		decodeJSON(t, response, &postPage)
		if len(postPage.Data) != 2 || !postPage.HasMore || postPage.NextCursor == "" {
			t.Fatalf("unexpected post cursor page: %+v", postPage)
		}
	})

	t.Run("authenticated thread and post creation", func(t *testing.T) {
		test := forumAPITest{db: db, router: router, prefix: "forum-test-" + uuid.NewString()}
		user, category, thread, post := test.seedForumWithUser(t)
		token, err := utils.GenerateAccessToken(user.ID, time.Now().UTC())
		if err != nil {
			t.Fatalf("generate access token: %v", err)
		}

		body := map[string]any{"category_id": category.ID, "title": "Created through HTTP", "body": "Opening post"}
		response := test.request(t, http.MethodPost, "/v1/forum/threads", token, body)
		assertStatus(t, response, http.StatusCreated)
		var created struct {
			Data struct {
				Thread    models.ForumThread `json:"thread"`
				FirstPost models.ForumPost   `json:"first_post"`
			} `json:"data"`
		}
		decodeJSON(t, response, &created)
		if created.Data.Thread.AuthorID != user.ID || created.Data.FirstPost.AuthorID != user.ID || created.Data.FirstPost.ThreadID != created.Data.Thread.ID {
			t.Fatalf("created records have wrong ownership or relationship: %+v", created.Data)
		}

		body = map[string]any{"body": "A reply", "parent_id": post.ID}
		response = test.request(t, http.MethodPost, "/v1/forum/threads/"+thread.ID.String()+"/posts", token, body)
		assertStatus(t, response, http.StatusCreated)

		response = test.request(t, http.MethodPost, "/v1/forum/threads/"+thread.ID.String()+"/posts", "", map[string]any{"body": "unauthorized"})
		assertStatus(t, response, http.StatusUnauthorized)

		response = test.request(t, http.MethodPost, "/v1/forum/threads/"+thread.ID.String()+"/posts", token, map[string]any{"body": ""})
		assertStatus(t, response, http.StatusBadRequest)
	})
}

func configureForumTestConfig(t *testing.T) {
	t.Helper()
	if err := os.Setenv("GIN_MODE", "test"); err != nil {
		t.Fatal(err)
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		t.Fatal("JWT_SECRET must be set")
	}
	configs.App.JWTSecret = []byte(secret)
}

func openForumTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		t.Fatalf("open PostgreSQL test database: %v", err)
	}
	return db
}

func (test forumAPITest) seedForum(t *testing.T) (models.ForumCategory, models.ForumThread, models.ForumPost) {
	t.Helper()
	category := models.ForumCategory{Name: test.prefix + "category", Description: "Forum API test category"}
	if err := test.db.Create(&category).Error; err != nil {
		t.Fatalf("create category: %v", err)
	}
	authorID := uuid.New()
	thread := models.ForumThread{CategoryID: category.ID, AuthorID: authorID, Title: test.prefix + "thread"}
	if err := test.db.Create(&thread).Error; err != nil {
		t.Fatalf("create thread: %v", err)
	}
	post := models.ForumPost{ThreadID: thread.ID, AuthorID: authorID, Body: "seed post"}
	if err := test.db.Create(&post).Error; err != nil {
		t.Fatalf("create post: %v", err)
	}
	t.Cleanup(func() {
		test.db.Unscoped().Delete(&models.ForumPost{}, "thread_id = ?", thread.ID)
		test.db.Unscoped().Delete(&models.ForumThread{}, "id = ?", thread.ID)
		test.db.Unscoped().Delete(&models.ForumCategory{}, "id = ?", category.ID)
	})
	return category, thread, post
}

func (test forumAPITest) seedForumWithUser(t *testing.T) (models.User, models.ForumCategory, models.ForumThread, models.ForumPost) {
	t.Helper()
	user := models.User{Username: test.prefix + "user", Email: test.prefix + "@example.com", Password: "test-password"}
	if err := test.db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	category := models.ForumCategory{Name: test.prefix + "category", Description: "Forum API test category"}
	if err := test.db.Create(&category).Error; err != nil {
		t.Fatalf("create category: %v", err)
	}
	thread := models.ForumThread{CategoryID: category.ID, AuthorID: user.ID, Title: test.prefix + "thread"}
	if err := test.db.Create(&thread).Error; err != nil {
		t.Fatalf("create thread: %v", err)
	}
	post := models.ForumPost{ThreadID: thread.ID, AuthorID: user.ID, Body: "seed post"}
	if err := test.db.Create(&post).Error; err != nil {
		t.Fatalf("create post: %v", err)
	}
	t.Cleanup(func() {
		test.db.Unscoped().Delete(&models.ForumPost{}, "thread_id = ?", thread.ID)
		test.db.Unscoped().Delete(&models.ForumThread{}, "id = ?", thread.ID)
		test.db.Unscoped().Delete(&models.ForumCategory{}, "id = ?", category.ID)
		test.db.Unscoped().Delete(&models.User{}, "id = ?", user.ID)
	})
	return user, category, thread, post
}

func (test forumAPITest) request(t *testing.T, method, path, token string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var requestBody *bytes.Reader
	if body == nil {
		requestBody = bytes.NewReader(nil)
	} else {
		encoded, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("encode request body: %v", err)
		}
		requestBody = bytes.NewReader(encoded)
	}
	request := httptest.NewRequest(method, path, requestBody)
	request.Header.Set("Content-Type", "application/json")
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}
	response := httptest.NewRecorder()
	test.router.ServeHTTP(response, request)
	return response
}

func assertStatus(t *testing.T, response *httptest.ResponseRecorder, expected int) {
	t.Helper()
	if response.Code != expected {
		t.Fatalf("expected HTTP %d, got %d: %s", expected, response.Code, strings.TrimSpace(response.Body.String()))
	}
}

func decodeJSON(t *testing.T, response *httptest.ResponseRecorder, target any) {
	t.Helper()
	if err := json.Unmarshal(response.Body.Bytes(), target); err != nil {
		t.Fatalf("decode response JSON: %v; body=%s", err, response.Body.String())
	}
}

func assertJSONID(t *testing.T, response *httptest.ResponseRecorder, expected uuid.UUID, field string) {
	t.Helper()
	var payload map[string]json.RawMessage
	decodeJSON(t, response, &payload)
	var data map[string]any
	if err := json.Unmarshal(payload[field], &data); err != nil {
		t.Fatalf("decode %s: %v", field, err)
	}
	if data["id"] != expected.String() {
		t.Fatalf("expected %s id %s, got %v", field, expected, data["id"])
	}
}

func assertJSONListContainsID(t *testing.T, response *httptest.ResponseRecorder, expected uuid.UUID, field string) {
	t.Helper()
	var payload map[string]json.RawMessage
	decodeJSON(t, response, &payload)
	var data []map[string]any
	if err := json.Unmarshal(payload[field], &data); err != nil {
		t.Fatalf("decode %s: %v", field, err)
	}
	for _, item := range data {
		if item["id"] == expected.String() {
			return
		}
	}
	t.Fatalf("expected %s to contain id %s", field, expected)
}
