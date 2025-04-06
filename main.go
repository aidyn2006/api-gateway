package main

import (
	"api-gateway/controllers"
	"api-gateway/database"
	"api-gateway/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()

	r := gin.Default()

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	protected.GET("/profile", controllers.Profile)

	r.POST("/create-order", controllers.CreateOrder)

	r.Run(":8085")

}
