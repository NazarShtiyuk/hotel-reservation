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

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (h *BookingHandler) HandleBookRoom(c *fiber.Ctx) error {
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

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	booking, err := h.store.BookingStore.GetBookingByID(c.Context(), oid)
	if err != nil {
		return err
	}

	userId, err := getAuthUserID(c)
	if err != nil {
		return err
	}

	if booking.UserID.Hex() != userId {
		return c.Status(http.StatusUnauthorized).JSON(types.GenericResponse{
			Type: "error",
			Msg:  "not authorized",
		})
	}

	if err := h.store.BookingStore.UpdateBooking(c.Context(), booking.ID.Hex(), bson.M{"$set": bson.M{"canceled": true}}); err != nil {
		return err
	}

	return c.JSON(types.GenericResponse{
		Type: "msg",
		Msg:  "updated",
	})
}

func (h *BookingHandler) HandleGetBookingByID(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	booking, err := h.store.BookingStore.GetBookingByID(c.Context(), oid)
	if err != nil {
		return err
	}

	return c.JSON(booking)
}

func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.BookingStore.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return err
	}

	return c.JSON(bookings)
}

func (h *BookingHandler) isRoomAvailableForBooking(ctx context.Context, roomID primitive.ObjectID, params types.BookRoomParams) (bool, error) {
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
