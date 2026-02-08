package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("🔄 Applying Timeline Overrides Migration")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using environment variables")
	}

	// Build database URL from env vars
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Build from individual vars
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		
		if host == "" || port == "" || user == "" || dbname == "" {
			log.Fatal("❌ Database configuration not found. Set DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME")
		}
		
		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)
		fmt.Printf("📡 Built connection string from env vars: postgres://%s:***@%s:%s/%s\n", user, host, port, dbname)
	}

	fmt.Println("📡 Connecting to database...")
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("❌ Failed to connect:", err)
	}
	defer pool.Close()

	// Test connection
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal("❌ Failed to ping database:", err)
	}
	fmt.Println("✅ Connected successfully\n")

	// Apply migrations
	migrations := []struct {
		name string
		sql  string
	}{
		{
			name: "Add timeline_overrides column",
			sql:  "ALTER TABLE service_tickets ADD COLUMN IF NOT EXISTS timeline_overrides JSONB",
		},
		{
			name: "Add parts_override column",
			sql:  "ALTER TABLE service_tickets ADD COLUMN IF NOT EXISTS parts_override JSONB",
		},
		{
			name: "Create GIN index on timeline_overrides",
			sql:  "CREATE INDEX IF NOT EXISTS idx_service_tickets_timeline_overrides ON service_tickets USING GIN (timeline_overrides)",
		},
		{
			name: "Add column comments",
			sql:  "COMMENT ON COLUMN service_tickets.timeline_overrides IS 'Admin-adjusted milestone data (JSON array of PublicMilestone)'",
		},
		{
			name: "Add column comments",
			sql:  "COMMENT ON COLUMN service_tickets.parts_override IS 'Admin-adjusted parts status and ETA (JSON object)'",
		},
	}

	fmt.Println("📝 Applying migrations...\n")
	successCount := 0
	
	for i, migration := range migrations {
		fmt.Printf("   [%d/%d] %s...", i+1, len(migrations), migration.name)
		
		_, err := pool.Exec(context.Background(), migration.sql)
		if err != nil {
			fmt.Printf(" ⚠️  Warning: %v\n", err)
		} else {
			fmt.Println(" ✅")
			successCount++
		}
	}

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("✅ Migration complete! (%d/%d successful)\n", successCount, len(migrations))
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	// Verify columns exist
	fmt.Println("🔍 Verifying columns...")
	var timelineCol, partsCol bool
	err = pool.QueryRow(context.Background(), `
		SELECT 
			EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='service_tickets' AND column_name='timeline_overrides'),
			EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='service_tickets' AND column_name='parts_override')
	`).Scan(&timelineCol, &partsCol)
	
	if err == nil {
		if timelineCol {
			fmt.Println("   ✅ timeline_overrides column exists")
		}
		if partsCol {
			fmt.Println("   ✅ parts_override column exists")
		}
		
		if timelineCol && partsCol {
			fmt.Println("\n🎉 Migration successful! Timeline adjustments will now persist.")
		}
	}
}
