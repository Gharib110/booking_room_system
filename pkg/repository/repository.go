package repository

// DatabaseRepo Holds all common functionalities that are common between all DBs Drivers
type DatabaseRepo interface {
	AllUsers() bool
}
