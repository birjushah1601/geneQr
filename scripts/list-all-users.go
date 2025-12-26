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
	fmt.Println("📋 ALL USERS WITH THEIR ORGANIZATIONS:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	query := `
		SELECT 
			u.email,
			u.full_name,
			u.status,
			CASE WHEN u.password_hash IS NOT NULL THEN 'YES' ELSE 'NO' END as has_password,
			COALESCE(o.name, 'No Organization') as org_name,
			COALESCE(o.org_type, 'N/A') as org_type,
			COALESCE(uo.role, 'N/A') as role,
			COALESCE(array_to_string(uo.permissions, ', '), 'N/A') as permissions
		FROM users u
		LEFT JOIN user_organizations uo ON u.id = uo.user_id AND uo.status = 'active'
		LEFT JOIN organizations o ON uo.organization_id = o.id
		ORDER BY o.org_type, u.email
	`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		log.Fatalf("Failed to query: %v", err)
	}
	defer rows.Close()

	count := 0
	currentOrgType := ""
	for rows.Next() {
		var email, fullName, status, hasPassword, orgName, orgType, role, permissions string
		err := rows.Scan(&email, &fullName, &status, &hasPassword, &orgName, &orgType, &role, &permissions)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		// Print org type header
		if orgType != currentOrgType {
			if count > 0 {
				fmt.Println()
			}
			fmt.Printf("═══════════════════════════════════════════════════════════════════════════════\n")
			fmt.Printf("📁 %s ORGANIZATIONS\n", orgType)
			fmt.Printf("═══════════════════════════════════════════════════════════════════════════════\n")
			currentOrgType = orgType
		}

		count++
		fmt.Printf("\n%d. %s\n", count, orgName)
		fmt.Printf("   📧 Email:       %s\n", email)
		fmt.Printf("   👤 Name:        %s\n", fullName)
		fmt.Printf("   🔐 Password:    %s\n", hasPassword)
		fmt.Printf("   📊 Status:      %s\n", status)
		fmt.Printf("   🎭 Role:        %s\n", role)
		fmt.Printf("   ⚙️  Permissions: %s\n", permissions)
	}

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("📊 Total Users: %d\n", count)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	// Now show organizations without users
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("⚠️  ORGANIZATIONS WITHOUT USERS:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	orgQuery := `
		SELECT o.name, o.org_type, o.status
		FROM organizations o
		LEFT JOIN user_organizations uo ON o.id = uo.organization_id
		WHERE uo.id IS NULL
		ORDER BY o.org_type, o.name
	`

	orgRows, err := pool.Query(ctx, orgQuery)
	if err != nil {
		log.Printf("Failed to query organizations: %v", err)
	} else {
		defer orgRows.Close()
		orgCount := 0
		for orgRows.Next() {
			var name, orgType, status string
			err := orgRows.Scan(&name, &orgType, &status)
			if err != nil {
				continue
			}
			orgCount++
			fmt.Printf("%d. %s (%s) - Status: %s\n", orgCount, name, orgType, status)
		}
		if orgCount == 0 {
			fmt.Println("✅ All organizations have users assigned!\n")
		} else {
			fmt.Printf("\n⚠️  %d organizations need users created\n\n", orgCount)
		}
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("💡 DEFAULT TEST LOGINS:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Println("🔐 System Admin:")
	fmt.Println("   Email:    admin@geneqr.com")
	fmt.Println("   Password: Admin@123\n")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
