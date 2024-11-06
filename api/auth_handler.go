package api

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/NazarShtiyuk/hotel-reservation/db"
	"github.com/NazarShtiyuk/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	store *db.Store
}

func NewAuthHandler(store *db.Store) *AuthHandler {
	return &AuthHandler{
		store: store,
	}
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var authParams types.AuthParams
	if err := c.BodyParser(&authParams); err != nil {
		return err
	}

	user, err := h.store.UserStore.GetUserByEmail(c.Context(), authParams.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("invalid credentials")
		}
		return err
	}

	if !types.IsValidPassword(user.EncryptedPassword, authParams.Password) {
		return fmt.Errorf("invalid credentials")
	}

	return c.JSON(types.AuthResponse{
		User:  user,
		Token: createTokenFromUser(user),
	})
}

func createTokenFromUser(user *types.User) string {
	now := time.Now()
	claims := jwt.MapClaims{
		"id":      user.ID,
		"email":   user.Email,
		"expires": now.Add(time.Hour * 4).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("failed to sign token with secret:", err)
	}

	return tokenStr
}
