package api

import (
	"log"
	"log/slog"
	"path/filepath"

	middleware "github.com/EstebanGitPro/motogo-backend/internal/middleware"
	"github.com/EstebanGitPro/motogo-backend/internal/platform/schema"
	"github.com/EstebanGitPro/motogo-backend/tools/utils"
	"github.com/gin-gonic/gin"
)

func routing(app *gin.Engine, dependencies *Dependencies) {
	slog.Info("Setting up routes")

	handler := New(dependencies.PersonService)

	moduleRoot, err := utils.FindModuleRoot()
	if err != nil {
		slog.Error("Error finding module root", slog.String("error", err.Error()))
		return
	}

	templatePath := filepath.Join(moduleRoot, "cmd", "api", "template", "*")
	slog.Debug("Template configuration", slog.String("template_path", templatePath))

	app.LoadHTMLGlob(templatePath)

	validators, err := schema.NewValidator(&schema.DefaultFileReader{})
	if err != nil {
		slog.Error("Error creating validator", slog.String("error", err.Error()))
		return
	}
	validator := middleware.NewMiddlewareValidator(validators)

	public := app.Group("/v1/motogo")
	{
		public.POST("/users", validator.WithValidateRegister(), handler.Save())
		public.POST("/auth/login", validator.WithValidateLogin(), handler.Login())
		public.GET("/auth/verify-email/:token", handler.VerifyEmail())
		public.GET("/email/status", handler.CheckEmailStatus())

		public.POST("/auth/password-recovery/send", handler.SendPasswordRecoveryEmail())
		public.POST("/auth/password-recovery/reset", handler.RecoveryPassword())
	}

	protected := app.Group("/v1/motogo")
	protected.Use(middleware.JWTAuthMiddleware(dependencies.Config.JWT))
	{
		protected.GET("/users/:id", handler.GetByID())
		protected.PATCH("/users/:id", validator.WithValidateUpdate(), handler.Update())
	}

	slog.Info("API routes configured successfully",
		slog.Int("public_routes", 6),
		slog.Int("protected_routes", 2))
}

func Bootstrap(app *gin.Engine) *Dependencies {

	dependencies, err := initDependencies()
	if err != nil {
		log.Fatal("Error initializing dependencies")
		return nil
	}

	routing(app, dependencies)

	return dependencies
}
