package main

import (
	"database/sql"
	"just-a-payment-gateway/backend/internal/auth/tokens"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/lib/pq"
)

func main() {
	err := godotenv.Load("../default.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	cfg := pq.Config{
		Host:           "localhost",
		Port:           5432,
		User:           "servicer",
		Password:       "SuperSecret1234",
		Database:       "payments",
		SSLMode:        pq.SSLModeDisable,
		ConnectTimeout: 5 * time.Second,
	}
	c, err := pq.NewConnectorConfig(cfg)
	if err != nil {
		log.Fatal(err)
	}
	db := sql.OpenDB(c)
	defer db.Close()

	tokenHandler := tokens.NewHandler(db)
	router := gin.Default()

	router.GET("/health", healthCheck)
	router.POST("/v1/tokens", tokenHandler.PostTokenizeCard)
	router.Run()
}
