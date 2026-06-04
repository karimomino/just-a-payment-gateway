package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

func healthCheck(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"status": "healthy"})
}

func main() {
	router := gin.Default()

	router.GET("/health", healthCheck)
	router.Run()
}
