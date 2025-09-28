package utils

import (
	"fmt"
	"lesson_server/constants"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = "SECRaqwem32mpomfwmfwlkm3034324klnlakndlnpq21230304923820938"

func GenerateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"id":  userID,
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("ошибка подписи токена: %w", err)
	}

	return tokenString, nil
}

func VerifyJWT(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return 0, fmt.Errorf("ошибка парсинга токена: %w", err)
	}

	if !token.Valid {
		return 0, fmt.Errorf("невалидный токен")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("неверный формат claims")
	}

	idStr, ok := claims["id"].(string)
	if !ok {
		return 0, fmt.Errorf("ID не найден или неверного формата")
	}

	id, _ := strconv.Atoi(idStr)

	return id, nil
}

func DeleteCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     constants.SessionKey,
		Value:    "",
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
	})

	c.Cookie(&fiber.Cookie{
		Name:     constants.StatusKey,
		Value:    "",
		MaxAge:   -1,
		HTTPOnly: false,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
	})
}

func AddCookie(c *fiber.Ctx, sessionID, status string) {
	c.Cookie(&fiber.Cookie{
		Name:     constants.SessionKey,
		Value:    sessionID,
		Expires:  time.Now().Add(48 * time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
	})

	c.Cookie(&fiber.Cookie{
		Name:     constants.StatusKey,
		Value:    status,
		Expires:  time.Now().Add(48 * time.Hour),
		HTTPOnly: false,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
	})
}

func LogginAPI(path, method string, status int, ip string, data any, message string) {
	fileName := "logger.txt"
	loggerStr := fmt.Sprintf("[%s]: %s - %s - %d - %s - %v - %s\n", time.Now().Format("2006-01-02 15:04:05"), path, method, status, ip, data, message)

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("[ERROR LOGGER]:", err.Error())
		return
	}
	defer file.Close()

	_, err = file.WriteString(loggerStr)
	if err != nil {
		log.Println("[ERROR LOGGER]:", err.Error())
		return
	}
}
