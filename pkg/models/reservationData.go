package models

import "time"

//ReservationData is the reservation data that we want to persist into the database
type ReservationData struct {
	FirstName string
	LastName  string
	Phone     string
	Email     string
}

//Users is the users data that we want to persist into the database
type Users struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

//Rooms is the model for persist the room's data into the database
type Rooms struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

//Restrictions is a model for persisting restrictions data into the database
type Restrictions struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

//Reservations is a model for persisting reservations data into the database
type Reservations struct {
	ID        int
	FirstName string
	LastName  string
	Phone     string
	Email     string
	StartDate time.Time
	EndDate   time.Time
	RoomID    int
	CreatedAt time.Time
	UpdatedAt time.Time
	Room      Rooms
}

//RoomsRestrictions is a model for persisting the rooms restrictions data into the database
type RoomsRestrictions struct {
	ID            int
	StartDate     time.Time
	EndDate       time.Time
	RoomID        int
	ReservationID int
	RestrictionID int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Room          Rooms
	Restriction   Restrictions
	Reservation   Reservations
}
