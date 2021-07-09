package dbrepo

import (
	"context"
	"errors"
	"github.com/DapperBlondie/booking_system/pkg/models"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

func (r *PostgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservation Insert a reservation in the database
func (r *PostgresDBRepo) InsertReservation(res models.Reservations) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	stmt := `INSERT INTO reservations 
    		 (first_name,last_name,email,phone,start_date,end_date,room_id,created_at,updated_at) 
			 VALUES ($1, $2, $3, $4, $5, $6,$7,$8,$9) RETURNING id`

	var newID int
	err := r.DB.QueryRowContext(ctx, stmt,
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
		return 0, err
	}

	return newID, nil
}

// InsertRestrictionRoom Insert a restriction for a reserved room into the database
func (r *PostgresDBRepo) InsertRestrictionRoom(restrict models.RoomsRestrictions) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*6)
	defer cancel()

	stmt := `INSERT INTO room_restrictions 
    (start_Date,end_date,room_id,reservation_id,restriction_id,created_at,updated_at) VALUES 
    ($1,$2,$3,$4,$5,$6,$7)`

	_, err := r.DB.ExecContext(ctx, stmt,
		restrict.StartDate,
		restrict.EndDate,
		restrict.RoomID,
		restrict.ReservationID,
		restrict.RestrictionID,
		time.Now(),
		time.Now())

	if err != nil {
		log.Println("We have some error for inserting restriction data into database : " + err.Error())
		return err
	}

	return nil
}

// SearchAvailabilityByDateByRoomID search across the availability of the times we have there
func (r *PostgresDBRepo) SearchAvailabilityByDateByRoomID(start, end time.Time, roomID int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*6)
	defer cancel()

	query := `SELECT count(id) FROM room_restrictions WHERE room_id=$1 
                                          AND $2 < end_date AND $3 > start_date`
	var numRows int

	err := r.DB.QueryRowContext(ctx, query, roomID, start, end).Scan(&numRows)
	if err != nil {
		log.Println("[SearchAvailabilityByDate] Error in getting the number of rows : " + err.Error() + "\n")
		return -1, err
	}

	return numRows, nil
}

//SearchAvailabilityForAllRooms Using for scanning the result of the searching availability of all rooms
func (r *PostgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Rooms, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*6)
	defer cancel()

	var rooms []models.Rooms

	query := `SELECT rooms.id, rooms.room_name
				FROM rooms WHERE rooms.id NOT IN 
				SELECT room_id FROM room_restrictions WHERE $1 < room_restrictions.end_date 
				AND  $2 > room_restrictions.start_date`

	rows, err := r.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		log.Println("[SearchAvailabilityForAllRooms] We have an error for getting the information : " + err.Error())
		return nil, err
	}

	for rows.Next() {
		var room models.Rooms
		err := rows.Scan(&room.ID, room.RoomName)
		if err != nil {
			log.Println("[SearchAvailabilityForAllRooms] We have an error for scanning the result : " + err.Error())
			return rooms, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

//GetRoomByID Using to get the specific room information by ID
func (r *PostgresDBRepo) GetRoomByID(id int) (models.Rooms, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*6)
	defer cancel()

	var room models.Rooms
	query := "SELECT id,room_name,created_at,updated_at FROM rooms WHERE id=$1"

	row := r.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt,
	)
	if err != nil {
		log.Println("An error occurred during scanning the row in GetRoomByID : " + err.Error() + "\n")
		return room, err
	}

	return room, nil
}

//GetUserByID Using for getting the User by its ID
func (r *PostgresDBRepo) GetUserByID(id int) (models.Users, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*6)
	defer cancel()

	query := `SELECT id,first_name,last_name,email,password,created_at,updated_at,access_level 
		FROM users WHERE id=$1`
	row := r.DB.QueryRowContext(ctx, query, id)
	var user models.Users

	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.AccessLevel,
	)

	if err != nil {
		log.Println("Error in scanning data from the queryRow : " + err.Error() + "\n")
		return user, err
	}

	return user, nil
}

//UpdateUser Using for updating a user
func (r *PostgresDBRepo) UpdateUser(user models.Users) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*6)
	defer cancel()

	query := `UPDATE users set first_name=$1,last_name=$2,email=$3,password=$4,created_at=$5,updated_at=$6,access_level=$7 
			WHERE id=$8`

	_, err := r.DB.ExecContext(ctx, query,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		time.Now(),
		&user.AccessLevel,
		&user.ID,
	)

	if err != nil {
		log.Println("Error during updating user data occurred : " + err.Error() + "\n")
		return err
	}

	return nil
}

//Authenticate Using for compare test and hashed password for qualifying the correct password
func (r *PostgresDBRepo) Authenticate(email, testPass string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*6)
	defer cancel()

	var id int
	var hashedPass string
	query := `SELECT id,password FROM users WHERE email=$1`

	row := r.DB.QueryRowContext(ctx, query, email)

	err := row.Scan(&id, &hashedPass)
	if err != nil {
		log.Println("Error during scanning the data for authenticate : " + err.Error() + "\n")
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(testPass))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		log.Println("Mismatch found between hashPass and testPass : " + err.Error() + "\n")
		return 0, "", errors.New("incorrect password passed")
	} else if err != nil {
		log.Println("Error occurred in CompareHashAndPass : " + err.Error() + "\n")
		return 0, "", err
	}

	return id, hashedPass, nil
}
