package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5430"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "med_platform"))

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer pool.Close()

	fmt.Println("\n🔍 Checking Siemens equipment...")
	
	query := `
		SELECT id, manufacturer_id, manufacturer_name, equipment_name 
		FROM equipment_registry 
		WHERE manufacturer_name LIKE '%Siemens%'
		LIMIT 5
	`
	
	rows, err := pool.Query(ctx, query)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()

	fmt.Println("\nSiemens Equipment:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	
	count := 0
	for rows.Next() {
		var id, manufacturerID, manufacturerName, equipmentName string
		err := rows.Scan(&id, &manufacturerID, &manufacturerName, &equipmentName)
		if err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}
		count++
		fmt.Printf("%d. %s\n", count, equipmentName)
		fmt.Printf("   ID: %s\n", id)
		fmt.Printf("   Manufacturer ID: %s\n", manufacturerID)
		fmt.Printf("   Manufacturer Name: %s\n\n", manufacturerName)
	}
	
	if count == 0 {
		fmt.Println("❌ No Siemens equipment found!")
		fmt.Println("\nLet's check what manufacturer_ids exist:")
		
		query2 := `
			SELECT DISTINCT manufacturer_id, manufacturer_name 
			FROM equipment_registry 
			WHERE manufacturer_id IS NOT NULL
			LIMIT 10
		`
		rows2, _ := pool.Query(ctx, query2)
		defer rows2.Close()
		
		fmt.Println("\nManufacturers with IDs:")
		for rows2.Next() {
			var mID, mName string
			rows2.Scan(&mID, &mName)
			fmt.Printf("  - %s: %s\n", mName, mID)
		}
	} else {
		fmt.Printf("\n✅ Found %d Siemens equipment\n", count)
		fmt.Println("\n📌 Expected manufacturer_id: 11afdeec-5dee-44d4-aa5b-952703536f10")
	}
	
	fmt.Println()
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
