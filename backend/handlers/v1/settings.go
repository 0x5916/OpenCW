package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"opencw/handlers/v1/common"
	"opencw/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SettingsHandler struct {
	DB *gorm.DB
}

func (h SettingsHandler) GetAllSettings(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	response := struct {
		CWSettings   common.CWSettingsResponse   `json:"cw_settings"`
		PageSettings common.PageSettingsResponse `json:"page_settings"`
	}{}

	if err := h.DB.Model(&models.CWSettings{}).Take(&response.CWSettings, "user_id = ?", user.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			defaultCWSettings := models.GetDefaultCWSettings()
			response.CWSettings = common.CWSettingsResponse{
				CharWPM:    defaultCWSettings.CharWPM,
				EffWPM:     defaultCWSettings.EffWPM,
				Freq:       defaultCWSettings.Freq,
				StartDelay: defaultCWSettings.StartDelay,
			}
		} else {
			slog.Error("Failed to get settings", "user_id", user.ID, "err", err)
			c.JSON(http.StatusInternalServerError, common.ErrorResponse{Error: "Failed to get settings"})
			return
		}
	}

	if err := h.DB.Model(&models.PageSettings{}).Take(&response.PageSettings, "user_id = ?", user.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			defaultPageSettings := models.GetDefaultPageSettings()
			response.PageSettings = common.PageSettingsResponse{
				Theme:     defaultPageSettings.Theme,
				Lang:      defaultPageSettings.Lang,
				CurLesson: defaultPageSettings.CurLesson,
			}
		} else {
			slog.Error("Failed to get settings", "user_id", user.ID, "err", err)
			c.JSON(http.StatusInternalServerError, common.ErrorResponse{Error: "Failed to get settings"})
			return
		}
	}

	c.JSON(http.StatusOK, response)
}

func (h SettingsHandler) GetCWSettings(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	var settings common.CWSettingsResponse
	if err := h.DB.Model(&models.CWSettings{}).Take(&settings, "user_id = ?", user.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, common.FromCwSettingsModel(models.GetDefaultCWSettings()))
			return
		}
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Error: "Failed to get settings"})
		slog.Error("Failed to get settings", "user_id", user.ID, "err", err)
		return
	}

	c.JSON(http.StatusOK, settings)
}

func (h SettingsHandler) GetPageSettings(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	var settings common.PageSettingsResponse
	if err := h.DB.Model(&models.PageSettings{}).Take(&settings, "user_id = ?", user.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, common.FromPageSettingsModel(models.GetDefaultPageSettings()))
			return
		}
		slog.Error("Failed to get settings", "user_id", user.ID, "err", err)
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Error: "Failed to get settings"})
		return
	}

	c.JSON(http.StatusOK, settings)
}

func (h SettingsHandler) UpdateCWSettings(c *gin.Context) {
	var input common.CWSettingsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: "Invalid request body"})
		return
	}

	user := c.MustGet("user").(*models.User)
	settings := models.CWSettings{UserID: user.ID}

	if err := h.DB.Where(&settings).Assign(&input).FirstOrCreate(&settings).Error; err != nil {
		slog.Error("Failed to update settings", "user_id", user.ID, "err", err)
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Error: "Failed to update settings"})
		return
	}

	c.JSON(http.StatusOK, common.MessageResponse{Message: "Settings updated"})
}

func (h SettingsHandler) UpdatePageSettings(c *gin.Context) {
	var input common.PageSettingsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: "Invalid request body"})
		return
	}

	user := c.MustGet("user").(*models.User)
	settings := models.PageSettings{UserID: user.ID}

	if err := h.DB.Where(&settings).Assign(&input).FirstOrCreate(&settings).Error; err != nil {
		slog.Error("Failed to update settings", "user_id", user.ID, "err", err)
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Error: "Failed to update settings"})
		return
	}

	c.JSON(http.StatusOK, common.MessageResponse{Message: "Settings updated"})
}
