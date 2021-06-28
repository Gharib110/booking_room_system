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

func (r *PostgresDBRepo) InsertRestrictionRoom(restrict models.RoomsRestrictions) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*6)
	defer cancel()

	stmt := `INSERT INTO `

	_, err := r.DB.ExecContext(ctx, stmt)
	if err != nil {
		log.Println("We have some error for inserting restriction data into database : " + err.Error())
		return err
	}

	return nil
}
