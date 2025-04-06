package main

import (
	"api-gateway/controllers"
	"api-gateway/database"
	"api-gateway/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	// Подключение к базе данных и автоматическая миграция
	database.Connect()

	// Инициализация маршрутов
	r := gin.Default()

	// Роуты регистрации и логина
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	// Защищённый маршрут для профиля
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	protected.GET("/profile", controllers.Profile)

	// Роут для создания заказа
	r.POST("/create-order", controllers.CreateOrder)

	// Запуск сервера на порту 8085
	r.Run(":8085") // Здесь указываем порт 8085
}
