package types

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID     primitive.ObjectID `bson:"userID" json:"userID"`
	RoomID     primitive.ObjectID `bson:"roomID" json:"roomID"`
	NumPersons int                `bson:"numPersons" json:"numPersons"`
	From       time.Time          `bson:"from" json:"from"`
	To         time.Time          `bson:"to" json:"to"`
	Canceled   bool               `bson:"canceled" json:"canceled"`
}

type BookRoomParams struct {
	From       time.Time `json:"from"`
	To         time.Time `json:"to"`
	NumPersons int       `json:"numPersons"`
}

func (p BookRoomParams) Validate() error {
	now := time.Now()
	if now.After(p.From) || now.After(p.To) {
		return fmt.Errorf("cannot book a room in the past")
	}
	if p.From.Unix() > p.To.Unix() {
		fmt.Println("to bigger than from")
		return fmt.Errorf("invalid data")
	}

	return nil
}
