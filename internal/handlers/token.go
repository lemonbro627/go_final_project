package handlers

import (
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/golang-jwt/jwt/v5"
)

// token создаёт подписанный токен
func createToken(userPassedPassword string, encKey string) (string, error) {
	// Secret key to sign and verify the token lately
	secret := []byte(encKey)

	signedPassword := HashPassword([]byte(userPassedPassword), secret)

	// создаём payload
	claims := jwt.MapClaims{
		"password": signedPassword, //захэшированный пароль вместе с секретным словом
	}

	// создаём jwt токен и указываем алгоритм хеширования и payload
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// получаем подписанный токен
	signedToken, err := jwtToken.SignedString(secret)
	if err != nil {
		log.Printf("failed to sign jwt: %s\n", err)
		return "", err
	}
	
	fmt.Println("Result token: " + string(signedToken[:]))

	return signedToken, nil
}

// функция для создания подписи
// HashPassword - это hash
func HashPassword(password []byte, secretKey []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(append(password, secretKey...)))
}
