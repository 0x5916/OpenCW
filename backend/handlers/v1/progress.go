package handlers

import (
	"log/slog"
	"net/http"
	"opencw/utils"

	"opencw/handlers/v1/common"
	"opencw/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProgressHandler struct {
	DB *gorm.DB
}

func (h ProgressHandler) GetAllProgress(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	var progresses []utils.ProgressResponse
	if err := h.DB.Model(&models.Progress{}).Find(&progresses, "user_id = ?", user.ID).Error; err != nil {
		slog.Error("Failed to query progress", "user_id", user.ID, "err", err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: "failed to query progress"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": progresses})
}

func (h ProgressHandler) AddProgress(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	var input common.ProgressInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: "Invalid request body"})
		return
	}

	progress := models.Progress{
		UserID:   user.ID,
		Lesson:   input.Lesson,
		CharWPM:  input.CharWPM,
		EffWPM:   input.EffWPM,
		Accuracy: *input.Accuracy,
	}

	if err := h.DB.Create(&progress).Error; err != nil {
		slog.Error("Failed to create progress", "user_id", user.ID, "lesson", progress.Lesson, "err", err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: "Failed to create progress"})
		return
	}
	c.JSON(http.StatusCreated, utils.MessageResponse{Message: "Progress Created"})
}
