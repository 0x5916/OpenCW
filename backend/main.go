package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"opencw/common"
	"opencw/utils"
	"os"
	"os/signal"
	"syscall"
	"time"

	"opencw/configs"
	"opencw/databases"
	"opencw/handlers/v1"
	"opencw/middlewares"
	"opencw/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if configs.GetGinMode() != "release" {
		if err := godotenv.Load(".env"); err != nil {
			slog.Error("Failed to load environment variables", "err", err)
			os.Exit(1)
		}
	}
	configs.Load()
	databases.Connect()

	r := gin.Default()

	if err := utils.RegisterCustomValidators(); err != nil {
		slog.Error("Failed to register custom validators", "err", err)
		os.Exit(1)
	}

	if !configs.App.IsRelease() {
		slog.Warn("App is not in production mode. Set GIN_MODE=release for production")
	}

	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" && isCORSAllowed(origin) {
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
		}
		c.Next()
	})

	v1 := r.Group("/v1")
	{
		// Health check endpoint
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":    "healthy",
				"timestamp": time.Now().Unix(),
			})
		})

		auth := v1.Group("/auth")
		authHandler := handlers.AuthHandler{DB: databases.DB}
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
			auth.POST("/refresh", authHandler.Refresh)
		}

		protected := v1.Group("/")
		protected.Use(middlewares.AuthRequired())
		protected.Use(middlewares.LoadUser(databases.DB))
		{
			authProtected := protected.Group("/auth")
			{
				authProtected.POST("/send-verification-email", authHandler.SendVerificationEmail)
				authProtected.POST("/verify-email", authHandler.VerifyEmail)
			}

			settings := protected.Group("/settings")
			settingsHandler := handlers.SettingsHandler{DB: databases.DB}
			{
				settings.GET("/all", settingsHandler.GetAllSettings)
				settings.GET("/cw", settingsHandler.GetCWSettings)
				settings.GET("/page", settingsHandler.GetPageSettings)
				settings.POST("/cw", settingsHandler.UpdateCWSettings)
				settings.POST("/page", settingsHandler.UpdatePageSettings)
			}

			user := protected.Group("/user")
			userHandler := handlers.UserHandler{DB: databases.DB}
			{
				user.GET("/me", userHandler.GetUserInfo)
				user.PUT("/callsign", userHandler.UpdateCallSign)
				user.PUT("/email", userHandler.UpdateEmail)
				user.PUT("/password", userHandler.UpdatePassword)
			}

			cwProgress := protected.Group("/cw")
			progressHandler := handlers.ProgressHandler{DB: databases.DB}
			{
				cwProgress.GET("/progress", progressHandler.GetAllProgress)
				cwProgress.PUT("/progress", progressHandler.AddProgress)
			}

			protected.GET("/hello", func(c *gin.Context) {
				user := c.MustGet("user").(models.User)
				c.JSON(http.StatusOK, common.MessageResponse{Message: "Hello, authenticated user {" + user.Username + "}!"})
			})
		}
	}

	// --- Graceful shutdown setup ---

	srv := &http.Server{
		Addr:    ":" + configs.App.Port,
		Handler: r,
	}

	// Context that cancels on SIGINT/SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Background clean-up goroutine — exits when ctx is cancelled
	go func() {
		ticker := time.NewTicker(3 * time.Hour)
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

	// Start server in a goroutine
	go func() {
		slog.Info("Server starting", "port", configs.App.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server error", "err", err)
			os.Exit(1)
		}
	}()

	// Block until signal received
	<-ctx.Done()
	slog.Info("Shutdown signal received")

	// Give in-flight requests 10 seconds to complete
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server forced to shutdown", "err", err)
		os.Exit(1)
	}

	// Close database connection
	sqlDB, err := databases.DB.DB()
	if err != nil {
		return
	}
	err = sqlDB.Close()
	if err != nil {
		return
	}

	slog.Info("Server exited gracefully")
}

// isCORSAllowed reports whether origin should receive CORS headers.
// Priority order:
//  1. CORS_ORIGINS env var (explicit allowlist, comma-separated)
//  2. Non-release mode → allow everything
//  3. Production → only https://opencw.net
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
	return origin == "https://opencw.net"
}
