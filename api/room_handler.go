package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/NazarShtiyuk/hotel-reservation/db"
	"github.com/NazarShtiyuk/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		return err
	}

	return c.JSON(rooms)
}

func (h *RoomHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.BookingStore.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return err
	}

	return c.JSON(bookings)
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params types.BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := params.Validate(); err != nil {
		return err
	}

	roomID := c.Params("id")
	rid, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return err
	}

	userID := c.Context().Value("userID").(string)
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	ok, err := h.isRoomAvailableForBooking(c.Context(), rid, params)
	if err != nil {
		return err
	}
	if !ok {
		return c.Status(http.StatusBadRequest).JSON(types.GenericResponse{
			Type: "error",
			Msg:  fmt.Sprintf("room %s already booked", roomID),
		})
	}

	booking := types.Booking{
		RoomID:     rid,
		UserID:     uid,
		From:       params.From,
		To:         params.To,
		NumPersons: params.NumPersons,
	}

	inserted, err := h.store.BookingStore.CreateBooking(c.Context(), &booking)
	if err != nil {
		return err
	}

	return c.JSON(inserted)
}

func (h *RoomHandler) isRoomAvailableForBooking(ctx context.Context, roomID primitive.ObjectID, params types.BookRoomParams) (bool, error) {
	filter := bson.M{
		"roomID": roomID,
		"to": bson.M{
			"$gte": params.From,
		},
	}

	bookings, err := h.store.BookingStore.GetBookings(ctx, filter)
	if err != nil {
		return false, err
	}

	ok := len(bookings) == 0
	return ok, nil
}
