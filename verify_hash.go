package main
import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)
func main() {
	password := "admin123"
	hash := "$2a$12$iGiLH0yA9AZv75byI.F2B.tJ.E1IxrMYd.0XV8No59WyfA4EswXo."
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == nil {
		fmt.Println("✅ Password 'admin123' MATCHES the hash!")
	} else {
		fmt.Printf("❌ Password doesn't match: %v\n", err)
	}
}
