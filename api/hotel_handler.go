package api

import (
	"errors"

	"github.com/NazarShtiyuk/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

func (h *HotelHandler) HandleGetHotelByID(c *fiber.Ctx) error {
	id := c.Params("id")
	hotel, err := h.store.HotelStore.GetHotelByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrNotFound("hotel")
		}
	}

	return c.JSON(hotel)
}

func (h *HotelHandler) HandleGetRoomsFromHotelByID(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID()
	}
	filter := bson.M{"hotelID": oid}
	rooms, err := h.store.RoomStore.GetRooms(c.Context(), filter)
	if err != nil {
		return ErrNotFound("rooms")
	}

	return c.JSON(rooms)
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	hotels, err := h.store.HotelStore.GetHotels(c.Context(), nil)
	if err != nil {
		return ErrNotFound("hotels")
	}

	return c.JSON(hotels)
}
