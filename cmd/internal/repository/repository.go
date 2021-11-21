package repository

import "learningGo/cmd/internal/modules"

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res modules.Reservation) error
}
