package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("🔄 Applying Authentication Migrations (Safe Mode)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	// Database connection string
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

	fmt.Println("✅ Connected to database successfully\n")

	// Apply migration 020 first (authentication core)
	fmt.Println("📄 Applying migration: 020_authentication_system.sql")
	if err := applyMigration(db, "database/migrations/020_authentication_system.sql"); err != nil {
		log.Printf("⚠️  Migration 020 had issues: %v\n", err)
		log.Println("Continuing to check what was created...\n")
	} else {
		fmt.Println("✅ Migration 020 applied successfully\n")
	}

	// Check service_tickets table exists before applying 021
	var ticketsExists bool
	err = db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'service_tickets')").Scan(&ticketsExists)
	if err != nil {
		log.Fatal("Failed to check service_tickets existence:", err)
	}

	if !ticketsExists {
		fmt.Println("⚠️  Warning: service_tickets table doesn't exist")
		fmt.Println("   Migration 021 requires service_tickets table from existing system")
		fmt.Println("   Skipping migration 021 for now\n")
	} else {
		fmt.Println("📄 Applying migration: 021_enhanced_tickets.sql")
		if err := applyMigration(db, "database/migrations/021_enhanced_tickets.sql"); err != nil {
			log.Printf("⚠️  Migration 021 had issues: %v\n", err)
		} else {
			fmt.Println("✅ Migration 021 applied successfully\n")
		}
	}

	// Verify tables were created
	fmt.Println("🔍 Verifying authentication tables...")
	
	coreTables := []string{
		"users",
		"otp_codes",
		"refresh_tokens",
		"auth_audit_log",
		"user_organizations",
		"roles",
		"notification_preferences",
	}

	createdCount := 0
	for _, table := range coreTables {
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
			var count int
			db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&count)
			fmt.Printf("  ✅ Table '%s' exists (rows: %d)\n", table, count)
			createdCount++
		} else {
			fmt.Printf("  ❌ Table '%s' does not exist\n", table)
		}
	}

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	if createdCount == len(coreTables) {
		fmt.Println("✅ All core authentication tables created successfully!")
	} else if createdCount > 0 {
		fmt.Printf("⚠️  %d/%d core tables created\n", createdCount, len(coreTables))
		fmt.Println("   Some tables may already exist or had errors")
	} else {
		fmt.Println("❌ No tables were created - check errors above")
	}
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
}

func applyMigration(db *sql.DB, filename string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Split into individual statements
	statements := strings.Split(string(content), ";")
	
	successCount := 0
	errorCount := 0
	
	for i, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" || strings.HasPrefix(stmt, "--") {
			continue
		}

		_, err := db.Exec(stmt)
		if err != nil {
			// Log but don't stop - table might already exist
			if strings.Contains(err.Error(), "already exists") {
				// Ignore "already exists" errors
				continue
			}
			errorCount++
			if errorCount <= 3 {
				log.Printf("  ⚠️  Statement %d error: %v", i+1, err)
			}
		} else {
			successCount++
		}
	}

	if errorCount > 0 {
		return fmt.Errorf("%d statements had errors (some may be expected)", errorCount)
	}

	return nil
}
