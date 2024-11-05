package db

const DBNAME = "hotel-reservation"
const DBURI = "mongodb://localhost:27017"

type Store struct {
	UserStore  UserStore
	HotelStore HotelStore
	RoomStore  RoomStore
}
