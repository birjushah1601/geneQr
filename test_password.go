package main

import (
	"database/sql"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	_ "github.com/lib/pq"
)

func main() {
	// Connect to database
	connStr := "host=localhost port=5430 user=postgres password=postgres dbname=med_platform sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Get user's password hash
	var passwordHash string
	err = db.QueryRow("SELECT password_hash FROM users WHERE email = $1", "admin@geneqr.com").Scan(&passwordHash)
	if err != nil {
		log.Fatal("Failed to get user:", err)
	}

	fmt.Println("Testing password verification...")
	fmt.Println("Email: admin@geneqr.com")
	fmt.Println("Password hash length:", len(passwordHash))
	
	// Test passwords
	passwords := []string{"Admin@123456", "admin123", "password"}
	
	for _, pwd := range passwords {
		err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(pwd))
		if err == nil {
			fmt.Printf("✅ Password '%s' matches!\n", pwd)
		} else {
			fmt.Printf("❌ Password '%s' does not match\n", pwd)
		}
	}
}
