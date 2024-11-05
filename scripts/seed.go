package main

import (
	"context"

	"log"

	"github.com/NazarShtiyuk/hotel-reservation/db"
	"github.com/NazarShtiyuk/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	userStore  db.UserStore
	roomStore  db.RoomStore
	hotelStore db.HotelStore
	ctx        = context.Background()
)

func seedUser(fname, lname, email string) {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: fname,
		LastName:  lname,
		Email:     email,
		Password:  "qwerty12345",
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = userStore.CreateUser(ctx, user)
	if err != nil {
		log.Fatal(err)
	}
}

func seedHotel(name, location string, rating int) error {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	rooms := []types.Room{
		{
			Size:      "small",
			BasePrice: 99.9,
		},
		{
			Size:      "normal",
			BasePrice: 149.9,
		},
		{
			Size:      "kingsize",
			BasePrice: 199.9,
		},
	}
	insertedHotel, err := hotelStore.CreateHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		_, err := roomStore.CreateRoom(ctx, &room)

		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

func main() {
	seedUser("Nazar", "Shtiyuk", "nazarshtiyuk@gmail.com")
	seedUser("Andrew", "Shway", "andrewshway@gmail.com")
	seedHotel("Bellucia", "France", 5)
	seedHotel("The cozy hotel", "The Nederlands", 4)
	seedHotel("Dont die in your sleep", "London", 4)
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
}
