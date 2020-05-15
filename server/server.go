package server

import (
	"github.com/biensupernice/krane/server/handler"
	"github.com/biensupernice/krane/server/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Config : server config
type Config struct {
	Port     string
	LogLevel string // release | debug
}

// Run : start server
func Run(cnf Config) {
	gin.SetMode(cnf.LogLevel)

	client := gin.New()

	// Middleware
	client.Use(gin.Recovery())
	client.Use(gin.Logger())
	client.Use(cors.Default())

	// Routes
	client.POST("/auth", handler.Auth)
	client.GET("/login", handler.Login)

	client.GET("/deployments", middleware.TokenAuthMiddleware(), handler.GetDeployments)
	client.POST("/deployments", middleware.TokenAuthMiddleware(), handler.CreateDeployment)

	client.PUT("/container/:containerID/stop", middleware.TokenAuthMiddleware(), handler.StopContainer)
	client.PUT("/container/:containerID/start", middleware.TokenAuthMiddleware(), handler.StartContainer)

	client.Run(cnf.Port)
}
