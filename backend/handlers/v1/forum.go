package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ForumHandler struct {
	DB *gorm.DB
}

func (h ForumHandler) CreateThread(c *gin.Context) {

}
