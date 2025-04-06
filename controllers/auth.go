package controllers

import (
	"api-gateway/database"
	"api-gateway/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("secret")

type Claims struct {
	UserID   uint
	Username string
	jwt.RegisteredClaims
}

func Register(c *gin.Context) {
	var body struct {
		Username string
		Password string
		Name     string
		Email    string
	}
	c.Bind(&body)

	hash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)

	user := models.User{
		Username: body.Username,
		Name:     body.Name,
		Email:    body.Email,
		Password: string(hash),
	}
	result := database.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User registration failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func Login(c *gin.Context) {
	var body struct {
		Username string
		Password string
	}
	c.Bind(&body)

	var user models.User
	database.DB.Where("username = ?", body.Username).First(&user)

	if user.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	expirationTime := time.Now().Add(time.Hour * 24)
	claims := &Claims{
		UserID:   user.UserID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtKey)

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func Profile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	result := database.DB.Where("user_id = ?", userIDUint).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username": user.Username,
		"name":     user.Name,
		"email":    user.Email,
	})
}
