package main

import (
	"net/http"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Gunakan CORS paling dasar
	r.Use(cors.Default())

	// Buat satu endpoint super sederhana
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Jalankan di port 8080
	r.Run(":8080")
}