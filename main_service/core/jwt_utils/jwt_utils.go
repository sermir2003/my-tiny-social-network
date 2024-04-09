package jwt_utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var JWT_SECRET []byte = nil

func StartUpJWT() {
	JWT_SECRET = []byte(os.Getenv("JWT_SECRET"))
}

func CreateJWT(id uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":       id.String(),
			"exp_time": time.Now().Add(15 * time.Minute).Format(time.RFC3339),
		},
	)

	tokenString, err := token.SignedString(JWT_SECRET)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetStrIdFromJWT(tokenString string) (id string, err error) {
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return JWT_SECRET, nil
		},
	)
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		exp_time, is_string := claims["exp_time"].(string)
		if !is_string {
			return "", errors.New("invalid token")
		}

		token_time, err := time.Parse(time.RFC3339, exp_time)
		if err != nil {
			return "", err
		}
		if token_time.Compare(time.Now()) == -1 {
			return "", errors.New("expired token")
		}

		id_str, is_string := claims["id"].(string)
		if !is_string {
			return "", errors.New("invalid token")
		}

		return id_str, nil
	}
	return "", errors.New("invalid token")
}
