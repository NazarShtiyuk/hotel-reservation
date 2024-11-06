package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Room struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	SeaSide bool               `bson:"seaside" json:"seaside"`
	Size    string             `bson:"size" json:"size"`
	Price   float64            `bson:"price" json:"price"`
	HotelID primitive.ObjectID `bson:"hotelID" json:"hotelID"`
}
