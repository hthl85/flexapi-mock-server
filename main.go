package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hthl85/flexapi-mock-server/controller"
	"github.com/hthl85/flexapi-mock-server/router"
	"github.com/hthl85/flexapi-mock-server/service"
	"github.com/hthl85/flexapi-mock-server/storage"
)

// SetUpEngine set up a server engine
func SetUpEngine(c *controller.Controller) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	engine.Use(gin.Logger())
	router.NewRouters(c, engine)
	return engine
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "5000"
	}

	storage := storage.NewStorage("account.db", "accountbucket")
	defer storage.BoltDB.Close()
	service := service.NewService(storage)
	controller := controller.NewController(service)
	server := SetUpEngine(controller)
	fmt.Printf("Start gin server. Listen on port:%s\n", port)
	server.Run(":" + port)
}
