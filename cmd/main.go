package main

import (
	"github.com/EstebanGitPro/motogo-backend/cmd/api"
	"github.com/EstebanGitPro/motogo-backend/config"
	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()
	api.Boostrap(app)
	if err := app.Run(config.MustLoadConfig().GetServerAddress()); err != nil {
		panic(err)
	}
}
