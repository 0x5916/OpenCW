package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"opencw/utils"

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
		CWSettings   utils.CWSettingsResponse   `json:"cw_settings"`
		PageSettings utils.PageSettingsResponse `json:"page_settings"`
	}{}

	if err := h.DB.Model(&models.CWSettings{}).Take(&response.CWSettings, "user_id = ?", user.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			defaultCWSettings := models.GetDefaultCWSettings()
			response.CWSettings = utils.CWSettingsResponse{
				CharWPM:    defaultCWSettings.CharWPM,
				EffWPM:     defaultCWSettings.EffWPM,
				Freq:       defaultCWSettings.Freq,
				StartDelay: defaultCWSettings.StartDelay,
			}
		} else {
			slog.Error("Failed to get settings", "user_id", user.ID, "err", err)
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: "Failed to get settings"})
			return
		}
	}

	if err := h.DB.Model(&models.PageSettings{}).Take(&response.PageSettings, "user_id = ?", user.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			defaultPageSettings := models.GetDefaultPageSettings()
			response.PageSettings = utils.PageSettingsResponse{
				Theme:     defaultPageSettings.Theme,
				Lang:      defaultPageSettings.Lang,
				CurLesson: defaultPageSettings.CurLesson,
			}
		} else {
			slog.Error("Failed to get settings", "user_id", user.ID, "err", err)
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: "Failed to get settings"})
			return
		}
	}

	c.JSON(http.StatusOK, response)
}

func (h SettingsHandler) GetCWSettings(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	var settings utils.CWSettingsResponse
	if err := h.DB.Model(&models.CWSettings{}).Take(&settings, "user_id = ?", user.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, common.FromCwSettingsModel(models.GetDefaultCWSettings()))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: "Failed to get settings"})
		slog.Error("Failed to get settings", "user_id", user.ID, "err", err)
		return
	}

	c.JSON(http.StatusOK, settings)
}

func (h SettingsHandler) GetPageSettings(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	var settings utils.PageSettingsResponse
	if err := h.DB.Model(&models.PageSettings{}).Take(&settings, "user_id = ?", user.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, common.FromPageSettingsModel(models.GetDefaultPageSettings()))
			return
		}
		slog.Error("Failed to get settings", "user_id", user.ID, "err", err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: "Failed to get settings"})
		return
	}

	c.JSON(http.StatusOK, settings)
}

func (h SettingsHandler) UpdateCWSettings(c *gin.Context) {
	var input common.CWSettingsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: "Invalid request body"})
		return
	}

	user := c.MustGet("user").(*models.User)
	settings := models.CWSettings{UserID: user.ID}

	if err := h.DB.Where(&settings).Assign(&input).FirstOrCreate(&settings).Error; err != nil {
		slog.Error("Failed to update settings", "user_id", user.ID, "err", err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: "Failed to update settings"})
		return
	}

	c.JSON(http.StatusOK, utils.MessageResponse{Message: "Settings updated"})
}

func (h SettingsHandler) UpdatePageSettings(c *gin.Context) {
	var input common.PageSettingsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: "Invalid request body"})
		return
	}

	user := c.MustGet("user").(*models.User)
	settings := models.PageSettings{UserID: user.ID}

	if err := h.DB.Where(&settings).Assign(&input).FirstOrCreate(&settings).Error; err != nil {
		slog.Error("Failed to update settings", "user_id", user.ID, "err", err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: "Failed to update settings"})
		return
	}

	c.JSON(http.StatusOK, utils.MessageResponse{Message: "Settings updated"})
}
