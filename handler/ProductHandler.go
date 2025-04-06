package handlers

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type ProductHandler struct {
	ProductServiceURL string
}

type CategoryHandler struct {
	CategoryServiceURL string
}

func NewRouter(productHandler *ProductHandler, categoryHandler *CategoryHandler) *gin.Engine {
	router := gin.Default()
	api := router.Group("/api/v1")
	{
		products := api.Group("/products")
		{
			products.GET("", productHandler.ProxyToService(http.MethodGet, "/products"))
			products.POST("", productHandler.ProxyToService(http.MethodPost, "/products"))
			products.PUT("/:id", productHandler.ProxyWithID(http.MethodPut, "/products"))
			products.DELETE("/:id", productHandler.ProxyWithID(http.MethodDelete, "/products"))
			products.GET("/:id", productHandler.ProxyWithID(http.MethodGet, "/products"))
		}
		categories := api.Group("/categories")
		{
			categories.GET("", categoryHandler.ProxyToService(http.MethodGet, "/categories"))
			categories.POST("", categoryHandler.ProxyToService(http.MethodPost, "/categories"))
			categories.GET("/:id", categoryHandler.ProxyWithID(http.MethodGet, "/categories"))
			categories.PATCH("/:id", categoryHandler.ProxyWithID(http.MethodPatch, "/categories"))
			categories.DELETE("/:id", categoryHandler.ProxyWithID(http.MethodDelete, "/categories"))
		}
	}
	return router
}

func (h *ProductHandler) ProxyToService(method, path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		fullURL := h.ProductServiceURL + path

		var body io.Reader
		if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
			body = c.Request.Body
		}

		req, err := http.NewRequest(method, fullURL, body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request: " + err.Error()})
			return
		}

		copyHeaders(c.Request.Header, req.Header)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to proxy request to service: " + err.Error()})
			return
		}
		defer resp.Body.Close()

		proxyResponse(resp, c)
	}
}

func (h *ProductHandler) ProxyWithID(method, path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		fullURL := fmt.Sprintf("%s%s/%s", h.ProductServiceURL, path, id)

		var body io.Reader
		if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
			body = c.Request.Body
		}

		req, err := http.NewRequest(method, fullURL, body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request: " + err.Error()})
			return
		}

		copyHeaders(c.Request.Header, req.Header)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to proxy request to service: " + err.Error()})
			return
		}
		defer resp.Body.Close()

		proxyResponse(resp, c)
	}
}

func (h *CategoryHandler) ProxyToService(method, path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		fullURL := h.CategoryServiceURL + path

		var body io.Reader
		if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
			body = c.Request.Body
		}

		req, err := http.NewRequest(method, fullURL, body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request: " + err.Error()})
			return
		}

		copyHeaders(c.Request.Header, req.Header)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to proxy request to service: " + err.Error()})
			return
		}
		defer resp.Body.Close()

		proxyResponse(resp, c)
	}
}

func (h *CategoryHandler) ProxyWithID(method, path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		fullURL := fmt.Sprintf("%s%s/%s", h.CategoryServiceURL, path, id)

		var body io.Reader
		if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
			body = c.Request.Body
		}

		req, err := http.NewRequest(method, fullURL, body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request: " + err.Error()})
			return
		}

		copyHeaders(c.Request.Header, req.Header)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to proxy request to service: " + err.Error()})
			return
		}
		defer resp.Body.Close()

		proxyResponse(resp, c)
	}
}

func copyHeaders(source, target http.Header) {
	for key, values := range source {
		for _, value := range values {
			target.Add(key, value)
		}
	}
}

func proxyResponse(resp *http.Response, c *gin.Context) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response: " + err.Error()})
		return
	}
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), buf.Bytes())
}
