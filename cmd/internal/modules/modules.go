package modules

import (
	"time"
)

//User is the user module
type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

//Room is the room module
type Room struct {
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

//Reservation module
type Reservation struct {
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
	Room      Room
}

type RoomRestriction struct {
	ID             int
	StartDate      time.Time
	EndDate        time.Time
	RoomID         int
	ReservationsID int
	RestrictionID  int
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Room           Room
	Reservation    Reservation
	Restriction    Restrictions
}

type MailData struct {
	To       string
	From     string
	Subject  string
	Content  string
	Template string
}
