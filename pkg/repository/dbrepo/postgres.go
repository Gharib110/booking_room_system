package dbrepo

import (
	"context"
	"github.com/DapperBlondie/booking_system/pkg/models"
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

// SearchAvailabilityByDate search across the availability of the times we have there
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
