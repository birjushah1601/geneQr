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

	siemensOrgID := "11afdeec-5dee-44d4-aa5b-952703536f10"
	
	var count int
	err = pool.QueryRow(ctx, `SELECT COUNT(*) FROM equipment_registry WHERE manufacturer_id = $1`, siemensOrgID).Scan(&count)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	fmt.Printf("\n✅ Equipment with manufacturer_id = %s: %d\n\n", siemensOrgID, count)
	
	// Show total
	var total int
	pool.QueryRow(ctx, `SELECT COUNT(*) FROM equipment_registry`).Scan(&total)
	fmt.Printf("Total equipment in database: %d\n\n", total)
	
	if count == 73 || count == total {
		fmt.Println("⚠️  WARNING: All equipment belongs to Siemens!")
		fmt.Println("   This means the manufacturer sees ALL equipment, which is expected.")
		fmt.Println("   To test filtering, you need equipment from other manufacturers.\n")
	} else {
		fmt.Printf("✅ Siemens has %d out of %d equipment (%.1f%%)\n", count, total, float64(count)/float64(total)*100)
		fmt.Println("   The filtering should show only these Siemens equipment.\n")
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
