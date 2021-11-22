package dbrepo

import (
	"errors"
	"learningGo/cmd/internal/modules"
	"time"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

//Insert reservation
func (m *testDBRepo) InsertReservation(res modules.Reservation) (int, error) {
	// if the room id is 2, then fail; otherwise, pass
	if res.RoomID == 2 {
		return 0, errors.New("some error)")
	}
	return 1, nil
}

//Insert room restrictions to database
func (m *testDBRepo) InsertRoomRestriction(res modules.RoomRestriction) error {
	if res.RoomID == 1000 {
		return errors.New("some error")
	}
	return nil
}

//SearchAvailabilityByDates return true if availability
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomId int) (bool, error) {
	return false, nil
}

//SearchAvailabilityForAllRooms return a slice of available rooms
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]modules.Room, error) {
	var rooms []modules.Room
	return rooms, nil
}

//GetRoomByID get room by id
func (m *testDBRepo) GetRoomByID(id int) (modules.Room, error) {
	var room modules.Room
	if id > 2 {
		return room, errors.New("Error")
	}
	return room, nil
}
