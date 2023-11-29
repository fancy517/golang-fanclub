package main

import (
	"fanclub/app"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	ginMode := gin.DebugMode
	if os.Getenv("mode") == "prod" {
		ginMode = gin.ReleaseMode
	}
	gin.SetMode(ginMode)
	app, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}
	defer app.CloseDB()
	app.Serve()
}
