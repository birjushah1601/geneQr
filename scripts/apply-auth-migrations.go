package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	// Database connection string from .env
	connStr := "host=localhost port=5430 user=postgres password=postgres dbname=med_platform sslmode=disable"

	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("✅ Connected to database successfully")

	// Read migration files
	migrations := []string{
		"database/migrations/020_authentication_system.sql",
		"database/migrations/021_enhanced_tickets.sql",
	}

	// Apply each migration
	for _, migrationFile := range migrations {
		fmt.Printf("\n📄 Applying migration: %s\n", migrationFile)

		// Read migration file
		content, err := os.ReadFile(migrationFile)
		if err != nil {
			log.Printf("❌ Failed to read migration file: %v", err)
			continue
		}

		// Execute migration
		_, err = db.Exec(string(content))
		if err != nil {
			log.Printf("❌ Failed to apply migration: %v", err)
			log.Printf("Migration content length: %d bytes", len(content))
			continue
		}

		fmt.Printf("✅ Migration applied successfully\n")
	}

	// Verify tables were created
	fmt.Println("\n🔍 Verifying created tables...")
	
	tables := []string{
		"users",
		"otp_codes",
		"refresh_tokens",
		"auth_audit_log",
		"user_organizations",
		"roles",
		"notification_preferences",
		"ticket_notifications",
		"whatsapp_conversations",
		"whatsapp_messages",
		"recaptcha_scores",
	}

	for _, table := range tables {
		var exists bool
		query := `SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = $1
		)`
		err := db.QueryRow(query, table).Scan(&exists)
		if err != nil {
			log.Printf("❌ Error checking table %s: %v", table, err)
			continue
		}
		
		if exists {
			// Count rows
			var count int
			db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&count)
			fmt.Printf("  ✅ Table '%s' exists (rows: %d)\n", table, count)
		} else {
			fmt.Printf("  ❌ Table '%s' does not exist\n", table)
		}
	}

	fmt.Println("\n✅ Migration process complete!")
}
