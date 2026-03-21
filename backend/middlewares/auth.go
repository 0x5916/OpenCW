package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"
	"opencw/common"
	"strconv"
	"strings"

	"opencw/configs"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			slog.Warn("Auth required: missing header", "ip", c.ClientIP(), "path", c.Request.RequestURI)
			c.JSON(http.StatusUnauthorized, common.ErrorResponse{Error: "Authorization header is required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		if tokenString == authHeader {
			slog.Warn("Auth required: invalid header format", "ip", c.ClientIP(), "path", c.Request.RequestURI)
			c.JSON(http.StatusUnauthorized, common.ErrorResponse{Error: "Authorization header must be in the format 'Bearer <token>'"})
			c.Abort()
			return
		}

		var claims jwt.RegisteredClaims
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return configs.App.JWTSecret, nil
		})

		if err != nil || !token.Valid {
			slog.Warn("Auth required: invalid token", "ip", c.ClientIP(), "path", c.Request.RequestURI, "err", err)
			c.JSON(http.StatusUnauthorized, common.ErrorResponse{Error: "Invalid token"})
			c.Abort()
			return
		}

		userID, err := strconv.ParseUint(claims.Subject, 10, 32)
		if err != nil {
			slog.Warn("Auth required: invalid token subject", "ip", c.ClientIP(), "subject", claims.Subject, "err", err)
			c.JSON(http.StatusUnauthorized, common.ErrorResponse{Error: "Invalid token"})
			c.Abort()
			return
		}

		c.Set("userID", uint(userID))
		c.Next()
	}
}
