package salt_utils

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
)

const saltSize = 16

func GenerateSalt() ([]byte, error) {
	var salt = make([]byte, saltSize)
	_, err := rand.Read(salt[:])
	return salt, err
}

func HashPassword(password string, salt []byte) string {
	var sha512Hasher = sha512.New()
	password_salt := append([]byte(password), salt...)
	sha512Hasher.Write(password_salt)
	var hashed = sha512Hasher.Sum(nil)
	return hex.EncodeToString(hashed)
}

func DoPasswordsMatch(password string, salt []byte, hashed string) bool {
	other_hashed := HashPassword(password, salt)
	return other_hashed == hashed
}
