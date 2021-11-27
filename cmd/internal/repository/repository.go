package repository

import (
	"learningGo/cmd/internal/modules"
	"time"
)

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res modules.Reservation) (int, error)
	InsertRoomRestriction(r modules.RoomRestriction) error
	SearchAvailabilityByDatesByRoomID(start, end time.Time, roomId int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time) ([]modules.Room, error)
	GetRoomByID(id int) (modules.Room, error)

	UpdateUser(u modules.User) error
	GetUserByID(id int) (modules.User, error)
	Authenticate(email, testPassword string) (int, string, error)

	AllReservations() ([]modules.Reservation, error)
	AllNewReservations() ([]modules.Reservation, error)

	GetReservationByID(id int) (modules.Reservation, error)

	UpdateReservation(u modules.Reservation) error
	DeleteReservation(id int) error
	UpdateProcessed(id, processed int) error

	AllRooms() ([]modules.Room, error)

	GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]modules.RoomRestriction, error)

	InsertBlockForRoom(id int, startDate time.Time) error
	DeleteBlockById(id int) error
}
