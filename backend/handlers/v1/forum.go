package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"opencw/common"
	"opencw/models"
	"opencw/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ForumHandler struct {
	DB *gorm.DB
}

const (
	forumDefaultPage  = 1
	forumDefaultLimit = 20
	forumMaxLimit     = 100
)

func parseForumPagination(c *gin.Context) (int, int) {
	page := forumDefaultPage
	limit := forumDefaultLimit

	if pageStr := c.Query("page"); pageStr != "" {
		if parsed, err := strconv.Atoi(pageStr); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
			if parsed > forumMaxLimit {
				limit = forumMaxLimit
			} else {
				limit = parsed
			}
		}
	}

	return page, limit
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

	page, limit := parseForumPagination(c)
	offset := (page - 1) * limit

	var total int64
	if err := h.DB.Model(&models.ForumThread{}).Where("category_id = ?", categoryID).Count(&total).Error; err != nil {
		slog.Error("Failed to count forum threads", "category_id", categoryID, "err", err)
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(common.ErrorCodeForumQueryFailed, "Failed to query forum threads"))
		return
	}

	var threads []models.ForumThread
	if err := h.DB.Where("category_id = ?", categoryID).
		Order("is_pinned DESC").
		Order("updated_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&threads).Error; err != nil {
		slog.Error("Failed to query forum threads", "category_id", categoryID, "err", err)
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(common.ErrorCodeForumQueryFailed, "Failed to query forum threads"))
		return
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

	c.JSON(http.StatusOK, gin.H{
		"data":  responses,
		"page":  page,
		"limit": limit,
		"total": total,
	})
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

	page, limit := parseForumPagination(c)
	offset := (page - 1) * limit

	var total int64
	if err := h.DB.Model(&models.ForumPost{}).Where("thread_id = ?", threadID).Count(&total).Error; err != nil {
		slog.Error("Failed to count forum posts", "thread_id", threadID, "err", err)
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(common.ErrorCodeForumQueryFailed, "Failed to query forum posts"))
		return
	}

	var posts []models.ForumPost
	if err := h.DB.Where("thread_id = ?", threadID).
		Order("created_at ASC").
		Offset(offset).
		Limit(limit).
		Find(&posts).Error; err != nil {
		slog.Error("Failed to query forum posts", "thread_id", threadID, "err", err)
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(common.ErrorCodeForumQueryFailed, "Failed to query forum posts"))
		return
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

	c.JSON(http.StatusOK, gin.H{
		"data":  responses,
		"page":  page,
		"limit": limit,
		"total": total,
	})
}

func (h ForumHandler) CreateThread(c *gin.Context) {
	user := utils.MustGetUser(c)

	var input common.CreateForumThreadInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(common.ErrorCodeInvalidRequestBody, "Invalid request body"))
		return
	}

	_ = user
	_ = input

	// Learning TODO: implement a transaction that creates the thread and its first post atomically.
	c.JSON(http.StatusNotImplemented, common.MessageResponse{Message: "TODO: implement thread creation"})

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

	_ = user
	_ = threadID
	_ = input

	// Learning TODO: enforce locked-thread and parent-post validation, then create post.
	c.JSON(http.StatusNotImplemented, common.MessageResponse{Message: "TODO: implement post creation"})
}
