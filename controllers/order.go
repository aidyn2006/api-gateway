package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

type OrderItem struct {
	ProductID uint `json:"product_id"`
	Quantity  uint `json:"quantity"`
}

type Order struct {
	UserID uint        `json:"user_id"`
	Status string      `json:"status,omitempty"`
	Total  float64     `json:"total"`
	Items  []OrderItem `json:"items"`
}

func CreateOrder(c *gin.Context) {
	var incomingOrder Order
	if err := c.ShouldBindJSON(&incomingOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	if len(incomingOrder.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Order must have at least one item",
		})
		return
	}

	orderData, err := json.Marshal(incomingOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to marshal order data: " + err.Error()})
		return
	}

	orderServiceURL := "http://localhost:8081/api/v1/orders"
	resp, err := http.Post(orderServiceURL, "application/json", bytes.NewBuffer(orderData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to send request to order service: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read response from order service: " + err.Error()})
		return
	}

	if resp.StatusCode != http.StatusCreated {
		c.JSON(resp.StatusCode, gin.H{
			"error": fmt.Sprintf("Failed to create order, status code: %d, response: %s", resp.StatusCode, string(body)),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Order created successfully",
		"data":    json.RawMessage(body),
	})
}
