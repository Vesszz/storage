package secutils

import (
	"golang.org/x/crypto/bcrypt"
	"math/rand"
)

const saltSymbols = "0123456789qwertyuiopasdfghjklzxcvbnm"

type HashedPassword struct {
	Value string
	Salt  string
}

func HashPassword(password string, saltLength int) (HashedPassword, error) {
	// Generate a random salt of the specified length
	salt := generateSalt(saltLength)
	// Hash the password with the generated salt
	hashedPassword, err := bcrypt.GenerateFromPassword(append([]byte(password), salt...), bcrypt.DefaultCost)
	if err != nil {
		return HashedPassword{}, err
	}

	// Return the hashed password and the salt
	return HashedPassword{
		Value: string(hashedPassword),
		Salt:  string(salt),
	}, nil
}

func generateSalt(length int) []byte {
	salt := make([]byte, length)
	for j := 0; j < length; j++ {
		randomInt := rand.Intn(len(saltSymbols))
		salt[j] = saltSymbols[randomInt]
	}
	return salt
}
