package main

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func healthCheckGin(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"status": "healthy"})
}

func healthCheck(w http.ResponseWriter, t *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := struct {
		Status string `json:"status"`
	}{
		Status: "healthy",
	}
	json.NewEncoder(w).Encode(response)
}
