package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Build DSN from environment
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5430")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "med_platform")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer pool.Close()

	fmt.Println("✅ Connected to database\n")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📋 EQUIPMENT MANUFACTURERS:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	// Check if manufacturer_id column exists
	checkColumnQuery := `
		SELECT column_name, data_type, is_nullable
		FROM information_schema.columns
		WHERE table_name = 'equipment_registry'
		AND column_name IN ('manufacturer_id', 'manufacturer_name')
		ORDER BY column_name
	`

	fmt.Println("Checking columns in equipment_registry table:")
	rows, err := pool.Query(ctx, checkColumnQuery)
	if err != nil {
		log.Fatalf("Failed to check columns: %v", err)
	}
	defer rows.Close()

	hasManufacturerID := false
	hasManufacturerName := false

	for rows.Next() {
		var colName, dataType, nullable string
		err := rows.Scan(&colName, &dataType, &nullable)
		if err != nil {
			continue
		}
		fmt.Printf("  ✅ Column: %-20s Type: %-15s Nullable: %s\n", colName, dataType, nullable)
		if colName == "manufacturer_id" {
			hasManufacturerID = true
		}
		if colName == "manufacturer_name" {
			hasManufacturerName = true
		}
	}

	fmt.Println()

	if !hasManufacturerID {
		fmt.Println("⚠️  manufacturer_id column does NOT exist!")
		fmt.Println("    The backend filtering won't work without this column.\n")
	}

	if !hasManufacturerName {
		fmt.Println("⚠️  manufacturer_name column does NOT exist!\n")
	}

	// Check equipment data
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📊 SAMPLE EQUIPMENT DATA:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	dataQuery := `
		SELECT 
			id,
			equipment_name,
			COALESCE(manufacturer_id::text, 'NULL') as manufacturer_id,
			COALESCE(manufacturer_name, 'NULL') as manufacturer_name
		FROM equipment_registry
		LIMIT 10
	`

	dataRows, err := pool.Query(ctx, dataQuery)
	if err != nil {
		log.Fatalf("Failed to query equipment: %v", err)
	}
	defer dataRows.Close()

	count := 0
	nullManufacturerID := 0
	nullManufacturerName := 0

	for dataRows.Next() {
		var id, name, mfgID, mfgName string
		err := dataRows.Scan(&id, &name, &mfgID, &mfgName)
		if err != nil {
			continue
		}
		count++
		fmt.Printf("%d. %s\n", count, name)
		fmt.Printf("   Manufacturer ID:   %s\n", mfgID)
		fmt.Printf("   Manufacturer Name: %s\n\n", mfgName)

		if mfgID == "NULL" {
			nullManufacturerID++
		}
		if mfgName == "NULL" {
			nullManufacturerName++
		}
	}

	// Get total counts
	var totalCount, totalNullID, totalNullName int
	countQuery := `
		SELECT 
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE manufacturer_id IS NULL) as null_id,
			COUNT(*) FILTER (WHERE manufacturer_name IS NULL OR manufacturer_name = '') as null_name
		FROM equipment_registry
	`
	err = pool.QueryRow(ctx, countQuery).Scan(&totalCount, &totalNullID, &totalNullName)
	if err != nil {
		log.Printf("Failed to get counts: %v", err)
	} else {
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println("📈 SUMMARY:")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("Total Equipment: %d\n", totalCount)
		fmt.Printf("With NULL manufacturer_id: %d (%.1f%%)\n", totalNullID, float64(totalNullID)/float64(totalCount)*100)
		fmt.Printf("With NULL manufacturer_name: %d (%.1f%%)\n\n", totalNullName, float64(totalNullName)/float64(totalCount)*100)

		if totalNullID > 0 {
			fmt.Println("⚠️  ISSUE FOUND:")
			fmt.Println("   Equipment has NULL manufacturer_id values!")
			fmt.Println("   The backend filtering by manufacturer_id won't work.")
			fmt.Println("   You need to populate manufacturer_id from manufacturer_name.\n")
		}
	}

	// Get manufacturer organization IDs
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🏢 MANUFACTURER ORGANIZATIONS:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	orgQuery := `
		SELECT id, name FROM organizations WHERE org_type = 'manufacturer' ORDER BY name LIMIT 10
	`
	orgRows, err := pool.Query(ctx, orgQuery)
	if err != nil {
		log.Printf("Failed to get manufacturers: %v", err)
	} else {
		defer orgRows.Close()
		for orgRows.Next() {
			var id, name string
			err := orgRows.Scan(&id, &name)
			if err != nil {
				continue
			}
			fmt.Printf("  • %s (ID: %s)\n", name, id)
		}
	}

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
