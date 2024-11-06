package api

import (
	"fmt"

	"github.com/NazarShtiyuk/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
)

func AdminAuth(store *db.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Context().Value("userID").(string)
		if !ok {
			return fmt.Errorf("not authorized")
		}

		user, _ := store.UserStore.GetUserByID(c.Context(), userID)
		if !user.Admin {
			return fmt.Errorf("not admin")
		}

		return c.Next()
	}

}
