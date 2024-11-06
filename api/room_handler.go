package api

import (
	"github.com/NazarShtiyuk/hotel-reservation/db"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := h.store.RoomStore.GetRooms(c.Context(), bson.M{})
	if err != nil {
		return ErrNotFound("rooms")
	}

	return c.JSON(rooms)
}
