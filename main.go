package main

import (
	"context"
	"flag"

	"log"

	"github.com/NazarShtiyuk/hotel-reservation/api"
	"github.com/NazarShtiyuk/hotel-reservation/db"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: api.ErrorHandler,
}

func main() {
	listenAddr := flag.String("listenAddr", ":3000", "The listen address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	var (
		app          = fiber.New(config)
		userStore    = db.NewMongoUserStore(client)
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		bookingStore = db.NewMongoBookingStore(client)
		store        = db.Store{UserStore: userStore, HotelStore: hotelStore, RoomStore: roomStore, BookingStore: bookingStore}
		apiv1        = app.Group("/api/v1", api.JWTAuthentication)
		adminv1      = apiv1.Group("/admin", api.AdminAuth(&store))
	)

	userHandler := api.NewUserHandler(&store)
	hotelHandler := api.NewHotelHandler(&store)
	authHandler := api.NewAuthHandler(&store)
	roomHandler := api.NewRoomHandler(&store)
	bookingHandler := api.NewBookingHandler(&store)

	//auth
	app.Post("/api/auth", authHandler.HandleAuthenticate)

	//users
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUserByID)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)

	//hotels
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotelByID)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRoomsFromHotelByID)

	//rooms
	apiv1.Get("/room", roomHandler.HandleGetRooms)

	//booking
	apiv1.Post("/room/:id/book", bookingHandler.HandleBookRoom)
	apiv1.Get("/booking/:id", bookingHandler.HandleGetBookingByID)
	apiv1.Get("/booking/:id/cancel", bookingHandler.HandleCancelBooking)
	//admin
	adminv1.Get("/booking", bookingHandler.HandleGetBookings)

	app.Listen(*listenAddr)
}
