package repository

import "github.com/DapperBlondie/booking_system/pkg/models"

// DatabaseRepo Holds all common functionalities that are common between all DBs Drivers
type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservations) (int, error)
	InsertRestrictionRoom(restrict models.RoomsRestrictions) error
}
