package main

import (
	"context"
	"log/slog"
	"net/http"
	"opencw/common"
	"opencw/handlers/v1"
	"opencw/middlewares"
	"opencw/utils"
	"os"
	"os/signal"
	"syscall"
	"time"

	"opencw/configs"
	"opencw/databases"
	"opencw/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func RouterV1Setup(engine *gin.Engine) {
	v1 := engine.Group("/v1")

	// Health check endpoint
	v1.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
		})
	})

	authHandler := handlers.AuthHandler{DB: databases.DB}
	auth := v1.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.POST("/logout", authHandler.Logout)
	auth.POST("/refresh", authHandler.Refresh)

	protected := v1.Group("/")
	protected.Use(middlewares.AuthRequired())
	protected.Use(middlewares.LoadUser(databases.DB))

	authProtected := protected.Group("/auth")
	authProtected.POST("/send-verification-email", authHandler.SendVerificationEmail)
	authProtected.POST("/verify-email", authHandler.VerifyEmail)

	settingsHandler := handlers.SettingsHandler{DB: databases.DB}
	settings := protected.Group("/settings")
	settings.GET("/all", settingsHandler.GetAllSettings)
	settings.GET("/cw", settingsHandler.GetCWSettings)
	settings.GET("/page", settingsHandler.GetPageSettings)
	settings.POST("/cw", settingsHandler.UpdateCWSettings)
	settings.POST("/page", settingsHandler.UpdatePageSettings)

	userHandler := handlers.UserHandler{DB: databases.DB}
	user := protected.Group("/user")
	user.GET("/me", userHandler.GetUserInfo)
	user.PUT("/callsign", userHandler.UpdateCallSign)
	user.PUT("/email", userHandler.UpdateEmail)
	user.PUT("/password", userHandler.UpdatePassword)

	progressHandler := handlers.ProgressHandler{DB: databases.DB}
	cwProgress := protected.Group("/cw")
	cwProgress.GET("/progress", progressHandler.GetAllProgress)
	cwProgress.PUT("/progress", progressHandler.AddProgress)

	forumHandler := handlers.ForumHandler{DB: databases.DB}
	forum := v1.Group("/forum")
	forum.GET("/categories", forumHandler.GetCategories)
	forum.GET("/categories/:categoryID/threads", forumHandler.GetThreadsByCategory)
	forum.GET("/threads/:threadID/posts", forumHandler.GetPostsByThread)

	forumProtected := protected.Group("/forum")
	forumProtected.POST("/threads", forumHandler.CreateThread)
	forumProtected.POST("/threads/:threadID/posts", forumHandler.CreatePost)

	protected.GET("/hello", func(c *gin.Context) {
		user := c.MustGet("user").(models.User)
		c.JSON(http.StatusOK, common.MessageResponse{Message: "Hello, authenticated user {" + user.Username + "}!"})
	})
}

func GracefulShutdown(srv *http.Server) {
	// Context that cancels on SIGINT/SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	startRefreshTokenCleanup(ctx)
	runServer(srv)

	// Block until signal received
	<-ctx.Done()
	slog.Info("Shutdown signal received")

	// Give in-flight requests a short window to complete.
	shutdownServer(context.Background(), srv)
	closeDatabase()

	slog.Info("Server exited gracefully")
}

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

	CORSSetup(r)
	RouterV1Setup(r)

	srv := &http.Server{
		Addr:              ":" + configs.App.Port,
		Handler:           r,
		ReadTimeout:       configs.App.ReadTimeout,
		ReadHeaderTimeout: configs.App.ReadHeaderTimeout,
		WriteTimeout:      configs.App.WriteTimeout,
		IdleTimeout:       configs.App.IdleTimeout,
	}

	GracefulShutdown(srv)
}
