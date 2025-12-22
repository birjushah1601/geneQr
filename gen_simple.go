package main
import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)
func main() {
	password := "password"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	fmt.Printf("Password: password\n")
	fmt.Printf("Hash: %s\n", string(hash))
}
