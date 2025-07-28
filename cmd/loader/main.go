package main

import (
	"fmt"
	"log"
	"os"
	"pricing-loader-service/internal/db"
	"pricing-loader-service/internal/loader"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or failed to load .env file")
	}
	// Load DB credentials from environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	if host == "" || port == "" || user == "" || password == "" || dbname == "" || sslmode == "" {
		log.Fatal("Database environment variables are not set properly")
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)
	dbConn, err := db.Connect(connStr)
	if err != nil {
		log.Fatalf("Failed to connect DB: %v", err)
	}
	defer dbConn.Close()

	log.Println("Connected to the database successfully.")

	// Initialize EC2PriceLoader
	ec2 := loader.NewEC2PriceLoader(dbConn)
	jsonPath := "data/index-ec2.json"

	if err := ec2.LoadProducts(jsonPath); err != nil {
		log.Fatalf("Failed to load products: %v", err)
	}

	log.Println("EC2 products loaded successfully.")

	if err := ec2.LoadPrices(jsonPath); err != nil {
		log.Fatalf("Failed to load prices: %v", err)
	}

	log.Println("EC2 pricing data loaded successfully.")
}
