package main

import (
	"context"
	"time"

	"log"

	"github.com/NazarShtiyuk/hotel-reservation/db"
	"github.com/NazarShtiyuk/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client       *mongo.Client
	userStore    db.UserStore
	roomStore    db.RoomStore
	hotelStore   db.HotelStore
	bookingStore db.BookingStore
	ctx          = context.Background()
)

func seedUser(admin bool, fname, lname, email, password string) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: fname,
		LastName:  lname,
		Email:     email,
		Password:  password,
	})
	if err != nil {
		log.Fatal(err)
	}

	user.Admin = admin

	insertedUser, err := userStore.CreateUser(ctx, user)
	if err != nil {
		log.Fatal(err)
	}

	return insertedUser
}

func seedHotel(name, location string, rating int) *types.Hotel {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	insertedHotel, err := hotelStore.CreateHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	return insertedHotel
}

func seedRoom(size string, seaside bool, price float64, hotelID primitive.ObjectID) *types.Room {
	room := types.Room{
		Size:    size,
		SeaSide: seaside,
		Price:   price,
		HotelID: hotelID,
	}

	insertedRoom, err := roomStore.CreateRoom(ctx, &room)
	if err != nil {
		log.Fatal(err)
	}

	return insertedRoom
}

func seedBooking(userID, roomID primitive.ObjectID, from, to time.Time) *types.Booking {
	booking := types.Booking{
		UserID: userID,
		RoomID: roomID,
		From:   from,
		To:     to,
	}

	insertedBooking, err := bookingStore.CreateBooking(ctx, &booking)
	if err != nil {
		log.Fatal(err)
	}

	return insertedBooking
}

func main() {
	admin := seedUser(true, "admin", "admin", "admin@gmail.com", "admin")
	seedUser(false, "Nazar", "Shtiyuk", "nazarshtiyuk@gmail.com", "qwerty12345")
	hotel := seedHotel("Bellucia", "France", 5)
	seedHotel("The cozy hotel", "The Nederlands", 4)
	seedHotel("Dont die in your sleep", "London", 4)
	seedRoom("small", false, 49.9, hotel.ID)
	seedRoom("normal", false, 99.9, hotel.ID)
	room := seedRoom("kingsize", true, 149.9, hotel.ID)
	seedBooking(admin.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 2))
}

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	userStore = db.NewMongoUserStore(client)
	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	bookingStore = db.NewMongoBookingStore(client)
}
