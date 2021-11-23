package dbrepo

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"learningGo/cmd/internal/modules"
	"time"
)

func (m *PostgreDBRepo) AllUsers() bool {
	return true
}

//Insert reservation
func (m *PostgreDBRepo) InsertReservation(res modules.Reservation) (int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int

	stmt := "insert into reservations (first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id"

	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return -1, err
	}

	return newID, nil
}

//Insert room restrictions to database
func (m *PostgreDBRepo) InsertRoomRestriction(r modules.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := "insert into room_restrictions (start_date, end_date, room_id, reservation_id, created_at, updated_at, restriction_id) values ($1, $2, $3, $4, $5, $6, $7)"

	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.RestrictionID,
		time.Now(),
		time.Now(),
		r.RestrictionID,
	)
	if err != nil {
		return err
	}
	return nil
}

//SearchAvailabilityByDates return true if availability
func (m *PostgreDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomId int) (bool, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var numRows int

	query := "select count(id) from room_restrictions where room_id = $1 and $2 < end_date and $3 > start_date;"
	row := m.DB.QueryRowContext(ctx, query, roomId, start, end)

	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}
	return false, nil
}

//SearchAvailabilityForAllRooms return a slice of available rooms
func (m *PostgreDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]modules.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var rooms []modules.Room

	query := "select r.id, r.room_name from rooms r where r.id not in (select rr.room_id from room_restrictions rr where $1 < rr.end_date and $2 > rr.start_date);"

	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var room modules.Room
		err := rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil
}

//GetRoomByID get room by id
func (m *PostgreDBRepo) GetRoomByID(id int) (modules.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room modules.Room

	query := "select id , room_name, created_at, updated_at from rooms where id = $1"

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&room.ID, &room.RoomName, &room.CreatedAt, &room.UpdatedAt)

	if err != nil {
		return room, err
	}

	return room, nil
}

//GetUserByID return user by id
func (m *PostgreDBRepo) GetUserByID(id int) (modules.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "select id , first_name, last_name, email, password, access_level, created_at, updated_at from users where id = $1"

	row := m.DB.QueryRowContext(ctx, query, id)

	var u modules.User
	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.AccessLevel,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return u, err
	}
	return u, nil
}

//UpdateUser update user info
func (m *PostgreDBRepo) UpdateUser(u modules.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "update users set first_name = $1, last_name = $2, email = $3, access_level = $4, updated_at = $5"

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.AccessLevel,
		time.Now(),
	)
	if err != nil {
		return err
	}
	return nil
}

//Authenticate the user
func (m *PostgreDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, "select id ,password from users where email = $1", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}

//AllReservations returns a slice of all reservations
func (m *PostgreDBRepo) AllReservations() ([]modules.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []modules.Reservation

	query := "select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date,r.end_date, r.room_id, r.created_at, r.updated_at, rm.id, rm.room_name from reservations r left join rooms rm on (r.room_id = rm.id) order by r.start_date asc"

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()
	for rows.Next() {
		var i modules.Reservation
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Room.ID,
			&i.Room.RoomName,
		)
		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}
	if err = rows.Err(); err != nil {
		return reservations, err
	}
	return reservations, nil

}

func (m *PostgreDBRepo) AllNewReservations() ([]modules.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []modules.Reservation

	query := "select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date,r.end_date, r.room_id, r.created_at, r.updated_at,r.processed, rm.id, rm.room_name from reservations r left join rooms rm on (r.room_id = rm.id) where processed = 0 order by r.start_date asc"

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()
	for rows.Next() {
		var i modules.Reservation
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Processed,
			&i.Room.ID,
			&i.Room.RoomName,
		)
		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}
	if err = rows.Err(); err != nil {
		return reservations, err
	}
	return reservations, nil

}

func (m *PostgreDBRepo) GetReservationByID(id int) (modules.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var res modules.Reservation

	query := "select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date,r.end_date, r.room_id, r.created_at, r.updated_at,r.processed, rm.id, rm.room_name from reservations r left join rooms rm on (r.room_id = rm.id) where r.id = $1 "

	rows := m.DB.QueryRowContext(ctx, query, id)
	err := rows.Scan(
		&res.ID,
		&res.FirstName,
		&res.LastName,
		&res.Email,
		&res.Phone,
		&res.StartDate,
		&res.EndDate,
		&res.RoomID,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.Processed,
		&res.Room.ID,
		&res.Room.RoomName,
	)
	if err != nil {
		return res, err
	}

	return res, nil

}
