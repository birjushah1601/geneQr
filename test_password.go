package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "admin123"
	hash := "$2a$12$EixZaYVK1fsbw1ZfbX3OXePaWxn96p36WQoeG6Lruj3vjPGga31lW"
	
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == nil {
		fmt.Println("✅ Password matches hash!")
	} else {
		fmt.Printf("❌ Password doesn't match: %v\n", err)
	}
}
