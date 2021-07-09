package repository

import (
	"github.com/DapperBlondie/booking_system/pkg/models"
	"time"
)

// DatabaseRepo Holds all common functionalities that are common between all DBs Drivers
type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservations) (int, error)
	InsertRestrictionRoom(restrict models.RoomsRestrictions) error
	SearchAvailabilityByDateByRoomID(start, end time.Time, roomID int) (int, error)
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Rooms, error)
	GetRoomByID(id int) (models.Rooms, error)
	GetUserByID(id int) (models.Users, error)
	UpdateUser(user models.Users) error
	Authenticate(email, testPass string) (int, string, error)
}
