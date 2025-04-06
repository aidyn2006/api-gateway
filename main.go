package main

import (
	"api-gateway/controllers"
	handlers "api-gateway/handler"
	"api-gateway/middleware"
)

func main() {
	productServiceURL := "http://product-service:8082/api/v1"
	categoryServiceURL := "http://category-service:8083/api/v1"

	productHandler := &handlers.ProductHandler{ProductServiceURL: productServiceURL}
	categoryHandler := &handlers.CategoryHandler{CategoryServiceURL: categoryServiceURL}

	router := handlers.NewRouter(productHandler, categoryHandler)

	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)

	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	protected.GET("/profile", controllers.Profile)

	router.POST("/create-order", controllers.CreateOrder)

	router.Run(":8085")
}
