package utils

import (
	"fmt"
	"lesson_server/constants"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetID(session string) (int, error) {
	if strings.TrimSpace(session) == "" {
		return 0, fmt.Errorf(constants.ErrUserNotFound)
	}

	sessionSplit := strings.Split(session, ":")
	if len(sessionSplit) < 2 {
		return 0, fmt.Errorf(constants.ErrUserNotFound)
	}

	id, err := strconv.Atoi(sessionSplit[0])
	if err != nil {
		return 0, fmt.Errorf(constants.ErrUserNotFound)
	}

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
