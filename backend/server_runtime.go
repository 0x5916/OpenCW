package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"opencw/configs"
	"opencw/databases"
	"opencw/models"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	refreshTokenCleanupInterval = 3 * time.Hour
	gracefulShutdownTimeout     = 10 * time.Second
	defaultProductionOrigin     = "https://opencw.net"
)

func CORSSetup(engine *gin.Engine) {
	engine.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "" || !isCORSAllowed(origin) {
			c.Next()
			return
		}

		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Header("Vary", "Origin")

		if c.Request.Method == http.MethodOptions {
			c.Header("Access-Control-Max-Age", "43200")
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})
}

func startRefreshTokenCleanup(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(refreshTokenCleanupInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				err := databases.DB.Unscoped().Where(
					"expires_at < ? OR revoked = true",
					time.Now(),
				).Delete(&models.RefreshToken{}).Error
				if err != nil {
					slog.Error("Failed to cleanup refresh tokens", "err", err)
				}
			case <-ctx.Done():
				slog.Info("Stopping refresh token cleanup goroutine")
				return
			}
		}
	}()
}

func runServer(srv *http.Server) {
	go func() {
		slog.Info("Server starting", "port", configs.App.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server error", "err", err)
			os.Exit(1)
		}
	}()
}

func shutdownServer(ctx context.Context, srv *http.Server) {
	shutdownCtx, cancel := context.WithTimeout(ctx, gracefulShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server forced to shutdown", "err", err)
		os.Exit(1)
	}
}

func closeDatabase() {
	sqlDB, err := databases.DB.DB()
	if err != nil {
		slog.Error("Failed to get SQL DB instance", "err", err)
		return
	}

	if err := sqlDB.Close(); err != nil {
		slog.Error("Failed to close database connection", "err", err)
	}
}

// isCORSAllowed reports whether origin should receive CORS headers.
// Priority order:
//  1. Non-release mode -> allow everything
//  2. CORS_ORIGINS env var (explicit allowlist, comma-separated)
//  3. Production -> only https://opencw.net
func isCORSAllowed(origin string) bool {
	if !configs.App.IsRelease() {
		return true
	}
	if len(configs.App.CORSOrigins) > 0 {
		for _, o := range configs.App.CORSOrigins {
			if o == origin {
				return true
			}
		}
		return false
	}
	return origin == defaultProductionOrigin
}
