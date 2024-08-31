package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type RoomsPostgres struct {
	db *sqlx.DB
}

func NewRoomsPostgres(db *sqlx.DB) *RoomsPostgres {
	return &RoomsPostgres{
		db: db,
	}
}

func (r *RoomsPostgres) JoinRoom(userId int) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	var roomId int
	query1 := fmt.Sprintf("SELECT id FROM %s WHERE users_quantity < $1 AND type = $2", roomsTable)
	row := tx.QueryRow(query1, 12, "public")
	err = row.Scan(&roomId)
	if err != nil {
		roomId = 0
	}
	if roomId == 0 {
		query2 := fmt.Sprintf("INSERT INTO %s (users_quantity, type) values ($1, $2) RETURNING id", roomsTable)
		row = tx.QueryRow(query2, 1, "public")
		err = row.Scan(&roomId)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

	}
	query3 := fmt.Sprintf("UPDATE %s SET users_quantity = users_quantity + 1 WHERE id = $1", roomsTable)
	_, err = tx.Exec(query3, roomId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	var name string
	query4 := fmt.Sprintf("SELECT name FROM %s WHERE id=$1", usersTable)
	row = tx.QueryRow(query4, userId)
	err = row.Scan(&name)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	query5 := fmt.Sprintf("INSERT INTO %s (user_id, room_id, user_name) values ($1, $2, $3)", playersTable)
	_, err = tx.Exec(query5, userId, roomId, name)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return roomId, tx.Commit()
}

func (r *RoomsPostgres) LeaveRoom(userId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	var roomId int
	query := fmt.Sprintf("SELECT room_id FROM %s WHERE user_id = $1", playersTable)
	row := tx.QueryRow(query, userId)
	err = row.Scan(&roomId)
	if err != nil {
		tx.Rollback()
		return err
	}
	query = fmt.Sprintf("DELETE FROM %s WHERE user_id = $1", playersTable)
	_, err = tx.Exec(query, userId)
	if err != nil {
		tx.Rollback()
		return err
	}
	query = fmt.Sprintf("UPDATE %s SET users_quantity = users_quantity - 1 WHERE id = $1", roomsTable)
	_, err = tx.Exec(query, roomId)
	if err != nil {
		tx.Rollback()
		return err
	}
	query = fmt.Sprintf("DELETE FROM %s WHERE users_quantity = $1", roomsTable)
	_, err = tx.Exec(query, 0)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
