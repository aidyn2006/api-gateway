package database

import (
	"api-gateway/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	// Настройки подключения к PostgreSQL
	dsn := "host=localhost user=postgres password=Na260206 dbname=go port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to the database", err)
	}

	// Автоматическая миграция таблиц
	err = database.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("failed to auto migrate", err)
	}

	// Устанавливаем глобальную переменную DB для доступа к базе данных
	DB = database
}
