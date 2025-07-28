package loader

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/buger/jsonparser"
)

type EC2PriceLoader struct {
	db *sql.DB
}

func NewEC2PriceLoader(db *sql.DB) *EC2PriceLoader {
	return &EC2PriceLoader{db: db}
}

func (loader *EC2PriceLoader) LoadProducts(jsonPath string) error {
	file, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	// offset := 0
	each := jsonparser.ObjectEach(file, func(key []byte, value []byte, dataType jsonparser.ValueType, offset2 int) error {
		sku := string(key)
		attributes, _, _, err := jsonparser.Get(value, "attributes")
		if err != nil {
			return fmt.Errorf("failed to get attributes: %v", err)
		}

		var attrs map[string]string
		if err := json.Unmarshal(attributes, &attrs); err != nil {
			return fmt.Errorf("failed to unmarshal attributes: %v", err)
		}

		instanceType := attrs["instanceType"]
		region := attrs["location"]
		os := attrs["operatingSystem"]
		tenancy := attrs["tenancy"]
		rawJSON := string(attributes)

		_, err = loader.db.Exec(`
			INSERT INTO aws_ec2_pricing (sku, instance_type, region, operating_system, tenancy, price_per_hour, raw_attributes)
			VALUES ($1, $2, $3, $4, $5, 0.0, $6)
			ON CONFLICT (sku) DO NOTHING
		`, sku, instanceType, region, os, tenancy, rawJSON)
		if err != nil {
			return fmt.Errorf("failed to insert SKU: %v", err)
		}
		log.Println("Inserted SKU:", sku)
		// Uncomment the next line if you want to print each SKU insertion
		// fmt.Printf("Inserted SKU: %s, Instance Type: %s, Region: %s, OS: %s, Tenancy: %s\n", sku, instanceType, region, os, tenancy)
		return nil
	}, "products")

	return each
}

func (loader *EC2PriceLoader) LoadPrices(jsonPath string) error {
	file, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	each := jsonparser.ObjectEach(file, func(skuKey []byte, skuValue []byte, vt jsonparser.ValueType, offset int) error {
		sku := string(skuKey)
		return jsonparser.ObjectEach(skuValue, func(pdKey []byte, pdValue []byte, vt jsonparser.ValueType, offset int) error {
			priceStr, err := jsonparser.GetString(pdValue, "pricePerUnit", "USD")
			if err != nil {
				return fmt.Errorf("failed to get pricePerUnit: %v", err)
			}
			var price float64
			fmt.Sscanf(priceStr, "%f", &price)

			_, err = loader.db.Exec(`
				UPDATE aws_ec2_pricing SET price_per_hour = $1 WHERE sku = $2
			`, price, sku)
			if err != nil {
				return fmt.Errorf("failed to update price: %v", err)
			}
			log.Println("Updated price for SKU:", sku, "to", price)
			// Uncomment the next line if you want to print each price update
			// fmt.Printf("Updated price for SKU: %s to %f\n", sku, price)
			return nil
		}, "priceDimensions")
	}, "terms", "OnDemand")

	return each
}
