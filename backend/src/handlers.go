package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func healthCheck(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"status": "healthy"})
}
