package main

import (
	"database/sql"
	"fmt"
	"just-a-payment-gateway/backend/internal/auth/tokens"
	"log"
	"net"
	"net/http"
	"time"

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
	http.HandleFunc("/health", healthCheck)
	http.HandleFunc("/v1/tokens", tokenHandler.PostTokenizeCard)
	// router := gin.Default()

	// router.GET("/health", healthCheck)
	// router.POST("/v1/tokens", tokenHandler.PostTokenizeCard)
	// router.Run()
	PORT := "8080"
	listener, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Server started listening on port: %s\nBase url is: http://localhost:%s\n", PORT, PORT)
	http.Serve(listener, nil)

	if err := http.ListenAndServe(":"+PORT, nil); err != nil {
		log.Fatal(err)
	}
	// log.Fatal(http.ListenAndServe(":"+PORT, nil))
}
