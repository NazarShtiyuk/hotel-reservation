package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func getAuthUserID(c *fiber.Ctx) (string, error) {
	userID, ok := c.Context().Value("userID").(string)
	if !ok {
		return "", fmt.Errorf("unauthorized")
	}

	return userID, nil
}
