package jwtgen

import (
	"errors"
	"time"

	"github.com/mikhailshtv/stockLkBack/internal/model"

	"github.com/golang-jwt/jwt/v5"
)

// Структуры для пользователя и токенов.
var (
	secretKey               = []byte("2k0935h84j39k2ks9df8h4fj3dk2s02kj9f8h4g5")
	ErrInvalidSigningMethod = errors.New("invalid signing method")
	ErrInvalidToken         = errors.New("invalid token")
)

// GenerateToken Функция для создания JWT токена.
func GenerateToken(userID int, login string, role model.UserRole) (string, error) {
	// Устанавливаем срок действия токена (1 час)
	expirationTime := time.Now().Add(1 * time.Hour)

	// Создаем JWT токен
	claims := &model.Claims{
		Login:  login,
		Role:   role,
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: expirationTime},
			Issuer:    "go-gin-jwt-example",
		},
	}

	// Создаем JWT токен и шифруем его методом HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен секретным ключом
	return token.SignedString(secretKey)
}

func ParseToken(tokenString string) (*model.Claims, error) {
	// Проверяем и парсим токен
	token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверка алгоритма подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}
		return secretKey, nil
	})

	// Если у нас произошла ошибка или токен не валиден, то вернем ошибку
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*model.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}
