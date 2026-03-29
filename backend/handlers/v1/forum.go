package handlers

import (
	// "opencw/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ForumHandler struct {
	DB *gorm.DB
}

func (h ForumHandler) CreateThread(c *gin.Context) {

}

func (h ForumHandler) CreatePost(c *gin.Context) {
	// user := utils.MustGetUser(c)
}
