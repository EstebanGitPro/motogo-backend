package api

import (
	//mid "github.com/EstebanGitPro/motogo-backend/internal/middleware"
	//schema "github.com/EstebanGitPro/motogo-backend/internal/platform/schema"
	"log"
	"path/filepath"

	"github.com/EstebanGitPro/motogo-backend/tools/utils"
	"github.com/gin-gonic/gin"
)

func routing(app *gin.Engine, dependencies *Dependencies) {
	handler := New(dependencies.PersonService)

	moduleRoot, err := utils.FindModuleRoot()
	if err != nil {
		log.Fatalf("Error finding module root: %v", err)
	}

	templatePath := filepath.Join(moduleRoot, "cmd", "api", "template", "*")
	log.Printf("Template path: %s", templatePath)

	app.LoadHTMLGlob(templatePath)

	app.POST("/v1/motogo/users", handler.Save())
	app.GET("/v1/motogo/auth/verify-email/:token", handler.VerifyEmail())
	app.GET("/v1/motogo/email/status", handler.CheckEmailStatus())
	app.POST("/v1/motogo/auth/login", handler.Login())
	app.GET("/v1/motogo/users/:id", handler.GetByID())  
	app.PATCH("/v1/motogo/users/:id", handler.Update()) 

}

func Boostrap(app *gin.Engine) {
	dependencies := initDependencies()
	if dependencies == nil {
		panic("dependencies not initialized")
	}
	dependencies.config.PrintConfig()
	routing(app, dependencies)
}
