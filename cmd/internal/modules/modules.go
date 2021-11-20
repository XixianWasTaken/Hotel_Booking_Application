package modules

import "time"

//Reservation holds the reservation data
type Reservation struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
}

//Users is the user module
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

//Rooms is the room module
type Rooms struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Restrictions struct {
	ID               int
	RestrictionsName string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

//Reservations module
type Reservations struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	RoomID    int
	CreatedAt time.Time
	UpdatedAt time.Time
	Room      Rooms
}

type RoomRestrictions struct {
	ID             int
	StartDate      time.Time
	EndDate        time.Time
	RoomID         int
	ReservationsID int
	RestrictionID  int
	CreatedAt      time.Time
	Room           Rooms
	Reservation    Reservations
	Restriction    Restrictions
}
