package handlers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"opencw/internal/common"
	"opencw/internal/models"
	"opencw/internal/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ForumHandler struct {
	DB *gorm.DB
}

const (
	forumDefaultPage   = 1
	forumDefaultLimit  = 20
	forumMaxLimit      = 100
	forumMaxPage       = 1_000_000
	forumCursorVersion = 1
)

type forumCursor struct {
	Version   int       `json:"v"`
	Scope     string    `json:"scope"`
	ParentID  uuid.UUID `json:"parent_id"`
	IsPinned  bool      `json:"is_pinned"`
	Timestamp time.Time `json:"timestamp"`
	ItemID    uuid.UUID `json:"item_id"`
}

func parseForumPagination(c *gin.Context, scope string, parentID uuid.UUID) (int, int, *forumCursor, error) {
	page := forumDefaultPage
	limit := forumDefaultLimit
	cursorValue, hasCursor := c.GetQuery("cursor")
	_, hasPage := c.GetQuery("page")

	if hasCursor {
		if hasPage || cursorValue == "" {
			return 0, 0, nil, errors.New("cursor cannot be combined with page and must not be empty")
		}

		if limitStr, hasLimit := c.GetQuery("limit"); hasLimit {
			parsed, err := strconv.Atoi(limitStr)
			if err != nil || parsed <= 0 {
				return 0, 0, nil, errors.New("limit must be a positive integer")
			}
			limit = min(forumMaxLimit, parsed)
		}
		if cursorValue == "first" {
			return 0, limit, &forumCursor{Version: forumCursorVersion, Scope: scope, ParentID: parentID}, nil
		}

		cursor, err := decodeForumCursor(cursorValue, scope, parentID)
		if err != nil {
			return 0, 0, nil, err
		}
		return 0, limit, &cursor, nil
	}

	if pageStr := c.Query("page"); pageStr != "" {
		if parsed, err := strconv.Atoi(pageStr); err == nil && parsed > 0 && parsed <= forumMaxPage {
			page = parsed
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
			limit = min(forumMaxLimit, parsed)
		}
	}

	return page, limit, nil, nil
}

func encodeForumCursor(cursor forumCursor) (string, error) {
	payload, err := json.Marshal(cursor)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(payload), nil
}

func decodeForumCursor(value, scope string, parentID uuid.UUID) (forumCursor, error) {
	payload, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return forumCursor{}, errors.New("invalid cursor")
	}

	var cursor forumCursor
	if err := json.Unmarshal(payload, &cursor); err != nil || cursor.Version != forumCursorVersion || cursor.Scope != scope || cursor.ParentID != parentID || cursor.ItemID == uuid.Nil || cursor.Timestamp.IsZero() {
		return forumCursor{}, errors.New("invalid cursor")
	}
	return cursor, nil
}

func (h ForumHandler) GetCategories(c *gin.Context) {
	var categories []models.ForumCategory
	if err := h.DB.Order("name ASC").Find(&categories).Error; err != nil {
		slog.Error("Failed to query forum categories", "err", err)
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(common.ErrorCodeForumQueryFailed, "Failed to query forum categories"))
		return
	}

	responses := make([]common.ForumCategoryResponse, 0, len(categories))
	for _, category := range categories {
		responses = append(responses, common.ForumCategoryResponse{
			ID:          category.ID,
			Name:        category.Name,
			Description: category.Description,
			CreatedAt:   category.CreatedAt,
			UpdatedAt:   category.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": responses})
}

func (h ForumHandler) GetThreadsByCategory(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("categoryID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(common.ErrorCodeInvalidRequestBody, "Invalid category id"))
		return
	}

	var category models.ForumCategory
	if err := h.DB.Select("id").Take(&category, "id = ?", categoryID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, common.NewErrorResponse(common.ErrorCodeForumCategoryNotFound, "Forum category not found"))
			return
		}

		slog.Error("Failed to query forum category", "category_id", categoryID, "err", err)
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(common.ErrorCodeForumQueryFailed, "Failed to query forum category"))
		return
	}

	page, limit, cursor, err := parseForumPagination(c, "threads", categoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(common.ErrorCodeInvalidRequestBody, err.Error()))
		return
	}

	var threads []models.ForumThread
	query := h.DB.Where("category_id = ?", categoryID)
	if cursor != nil && cursor.ItemID != uuid.Nil {
		query = query.Where("is_pinned < ? OR (is_pinned = ? AND (updated_at < ? OR (updated_at = ? AND id < ?)))", cursor.IsPinned, cursor.IsPinned, cursor.Timestamp, cursor.Timestamp, cursor.ItemID)
	} else {
		query = query.Offset((page - 1) * limit)
	}
	if err := query.Order("is_pinned DESC").Order("updated_at DESC").Order("id DESC").Limit(limit + 1).Find(&threads).Error; err != nil {
		slog.Error("Failed to query forum threads", "category_id", categoryID, "err", err)
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(common.ErrorCodeForumQueryFailed, "Failed to query forum threads"))
		return
	}
	hasMore := len(threads) > limit
	if hasMore {
		threads = threads[:limit]
	}

	responses := make([]common.ForumThreadResponse, 0, len(threads))
	for _, thread := range threads {
		responses = append(responses, common.ForumThreadResponse{
			ID:         thread.ID,
			CategoryID: thread.CategoryID,
			AuthorID:   thread.AuthorID,
			Title:      thread.Title,
			IsPinned:   thread.IsPinned,
			IsLocked:   thread.IsLocked,
			CreatedAt:  thread.CreatedAt,
			UpdatedAt:  thread.UpdatedAt,
		})
	}

	if cursor != nil {
		response := gin.H{"data": responses, "limit": limit, "has_more": hasMore, "next_cursor": ""}
		if hasMore {
			last := threads[len(threads)-1]
			response["next_cursor"], err = encodeForumCursor(forumCursor{forumCursorVersion, "threads", categoryID, last.IsPinned, last.UpdatedAt, last.ID})
			if err != nil {
				c.JSON(http.StatusInternalServerError, common.NewErrorResponse(common.ErrorCodeForumQueryFailed, "Failed to encode forum cursor"))
				return
			}
		}
		c.JSON(http.StatusOK, response)
		return
	}

	var total int64
	if err := h.DB.Model(&models.ForumThread{}).Where("category_id = ?", categoryID).Count(&total).Error; err != nil {
		slog.Error("Failed to count forum threads", "category_id", categoryID, "err", err)
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(common.ErrorCodeForumQueryFailed, "Failed to query forum threads"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": responses, "page": page, "limit": limit, "total": total})
}

func (h ForumHandler) GetPostsByThread(c *gin.Context) {
	threadID, err := uuid.Parse(c.Param("threadID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(common.ErrorCodeInvalidRequestBody, "Invalid thread id"))
		return
	}

	var thread models.ForumThread
	if err := h.DB.Select("id").Take(&thread, "id = ?", threadID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, common.NewErrorResponse(common.ErrorCodeForumThreadNotFound, "Forum thread not found"))
			return
		}

		slog.Error("Failed to query forum thread", "thread_id", threadID, "err", err)
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(common.ErrorCodeForumQueryFailed, "Failed to query forum thread"))
		return
	}

	page, limit, cursor, err := parseForumPagination(c, "posts", threadID)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(common.ErrorCodeInvalidRequestBody, err.Error()))
		return
	}

	var posts []models.ForumPost
	query := h.DB.Where("thread_id = ?", threadID)
	if cursor != nil && cursor.ItemID != uuid.Nil {
		query = query.Where("created_at > ? OR (created_at = ? AND id > ?)", cursor.Timestamp, cursor.Timestamp, cursor.ItemID)
	} else {
		query = query.Offset((page - 1) * limit)
	}
	if err := query.Order("created_at ASC").Order("id ASC").Limit(limit + 1).Find(&posts).Error; err != nil {
		slog.Error("Failed to query forum posts", "thread_id", threadID, "err", err)
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(common.ErrorCodeForumQueryFailed, "Failed to query forum posts"))
		return
	}
	hasMore := len(posts) > limit
	if hasMore {
		posts = posts[:limit]
	}

	responses := make([]common.ForumPostResponse, 0, len(posts))
	for _, post := range posts {
		responses = append(responses, common.ForumPostResponse{
			ID:        post.ID,
			ThreadID:  post.ThreadID,
			AuthorID:  post.AuthorID,
			Body:      post.Body,
			ParentID:  post.ParentID,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		})
	}

	if cursor != nil {
		response := gin.H{"data": responses, "limit": limit, "has_more": hasMore, "next_cursor": ""}
		if hasMore {
			last := posts[len(posts)-1]
			response["next_cursor"], err = encodeForumCursor(forumCursor{forumCursorVersion, "posts", threadID, false, last.CreatedAt, last.ID})
			if err != nil {
				c.JSON(http.StatusInternalServerError, common.NewErrorResponse(common.ErrorCodeForumQueryFailed, "Failed to encode forum cursor"))
				return
			}
		}
		c.JSON(http.StatusOK, response)
		return
	}

	var total int64
	if err := h.DB.Model(&models.ForumPost{}).Where("thread_id = ?", threadID).Count(&total).Error; err != nil {
		slog.Error("Failed to count forum posts", "thread_id", threadID, "err", err)
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(common.ErrorCodeForumQueryFailed, "Failed to query forum posts"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": responses, "page": page, "limit": limit, "total": total})
}

func (h ForumHandler) GetThread(c *gin.Context) {
	threadID, err := uuid.Parse(c.Param("threadID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(common.ErrorCodeInvalidRequestBody, "Invalid thread id"))
		return
	}

	var thread models.ForumThread
	if err := h.DB.First(&thread, "id = ?", threadID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, common.NewErrorResponse(common.ErrorCodeForumThreadNotFound, "Forum thread not found"))
			return
		}

		slog.Error("Failed to query forum thread", "thread_id", threadID, "err", err)
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(common.ErrorCodeForumQueryFailed, "Failed to query forum thread"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": forumThreadResponse(thread)})
}

func (h ForumHandler) CreateThread(c *gin.Context) {
	user := utils.MustGetUser(c)

	var input common.CreateForumThreadInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(common.ErrorCodeInvalidRequestBody, "Invalid request body"))
		return
	}

	var category models.ForumCategory
	if err := h.DB.Select("id").First(&category, "id = ?", input.CategoryID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, common.NewErrorResponse(common.ErrorCodeForumCategoryNotFound, "Forum category not found"))
			return
		}

		slog.Error("Failed to query forum category", "category_id", input.CategoryID, "err", err)
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(common.ErrorCodeForumQueryFailed, "Failed to query forum category"))
		return
	}

	var thread models.ForumThread
	var post models.ForumPost
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		thread = models.ForumThread{
			CategoryID: input.CategoryID,
			AuthorID:   user.ID,
			Title:      input.Title,
		}
		if err := tx.Create(&thread).Error; err != nil {
			return err
		}

		post = models.ForumPost{
			ThreadID: thread.ID,
			AuthorID: user.ID,
			Body:     input.Body,
		}
		return tx.Create(&post).Error
	})
	if err != nil {
		slog.Error("Failed to create forum thread", "category_id", input.CategoryID, "author_id", user.ID, "err", err)
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(common.ErrorCodeForumCreateFailed, "Failed to create forum thread"))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": common.ForumThreadCreatedResponse{
		Thread:    forumThreadResponse(thread),
		FirstPost: forumPostResponse(post),
	}})

}

func (h ForumHandler) CreatePost(c *gin.Context) {
	user := utils.MustGetUser(c)

	threadID, err := uuid.Parse(c.Param("threadID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(common.ErrorCodeInvalidRequestBody, "Invalid thread id"))
		return
	}

	var input common.CreateForumPostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(common.ErrorCodeInvalidRequestBody, "Invalid request body"))
		return
	}

	var thread models.ForumThread
	if err := h.DB.First(&thread, "id = ?", threadID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, common.NewErrorResponse(common.ErrorCodeForumThreadNotFound, "Forum thread not found"))
			return
		}

		slog.Error("Failed to query forum thread", "thread_id", threadID, "err", err)
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(common.ErrorCodeForumQueryFailed, "Failed to query forum thread"))
		return
	}
	if thread.IsLocked {
		c.JSON(http.StatusConflict, common.NewErrorResponse(common.ErrorCodeForumThreadLocked, "Forum thread is locked"))
		return
	}

	var post models.ForumPost
	err = h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&thread, "id = ?", threadID).Error; err != nil {
			return err
		}
		if thread.IsLocked {
			return gorm.ErrInvalidData
		}

		if input.ParentID != nil {
			var parent models.ForumPost
			if err := tx.Select("id", "thread_id").First(&parent, "id = ? AND thread_id = ?", *input.ParentID, threadID).Error; err != nil {
				return err
			}
		}

		post = models.ForumPost{
			ThreadID: threadID,
			AuthorID: user.ID,
			Body:     input.Body,
			ParentID: input.ParentID,
		}
		if err := tx.Create(&post).Error; err != nil {
			return err
		}
		return tx.Model(&thread).Update("updated_at", time.Now().UTC()).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrInvalidData) {
			c.JSON(http.StatusConflict, common.NewErrorResponse(common.ErrorCodeForumThreadLocked, "Forum thread is locked"))
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) && input.ParentID != nil {
			c.JSON(http.StatusBadRequest, common.NewErrorResponse(common.ErrorCodeForumParentPostInvalid, "Parent post does not belong to this thread"))
			return
		}

		slog.Error("Failed to create forum post", "thread_id", threadID, "author_id", user.ID, "err", err)
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(common.ErrorCodeForumCreateFailed, "Failed to create forum post"))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": forumPostResponse(post)})
}

func forumThreadResponse(thread models.ForumThread) common.ForumThreadResponse {
	return common.ForumThreadResponse{
		ID:         thread.ID,
		CategoryID: thread.CategoryID,
		AuthorID:   thread.AuthorID,
		Title:      thread.Title,
		IsPinned:   thread.IsPinned,
		IsLocked:   thread.IsLocked,
		CreatedAt:  thread.CreatedAt,
		UpdatedAt:  thread.UpdatedAt,
	}
}

func forumPostResponse(post models.ForumPost) common.ForumPostResponse {
	return common.ForumPostResponse{
		ID:        post.ID,
		ThreadID:  post.ThreadID,
		AuthorID:  post.AuthorID,
		Body:      post.Body,
		ParentID:  post.ParentID,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}
}
