package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
)

func MigrateDB(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS staff (
			staff_id BIGSERIAL PRIMARY KEY,
			email    VARCHAR(255) NOT NULL
		);

		CREATE TYPE reservationenum AS ENUM ('CONFIRMED', 'CANCELLED');

		CREATE TABLE IF NOT EXISTS reservation
		(
			reservation_id BIGSERIAL        NOT NULL UNIQUE PRIMARY KEY,
			name           VARCHAR(100)     NOT NULL,
			email          VARCHAR(255)     NOT NULL,
			status         reservationenum  NOT NULL,
			created_at     TIMESTAMP        NOT NULL,
			scheduled_for  TIMESTAMP        NOT NULL,
			expire_at      TIMESTAMP        NOT NULL,
			staff_id       BIGINT           NOT NULL,
			CONSTRAINT FK_reservation_to_staff_staff_id
				FOREIGN KEY (staff_id)
					REFERENCES staff(staff_id)
					ON DELETE RESTRICT
					ON UPDATE RESTRICT,
			CONSTRAINT EX_reservation_overlap_constraint
				EXCLUDE USING gist (
					staff_id WITH =,
					tsrange(scheduled_for, expire_at) WITH &&,
					(CASE WHEN status = 'CONFIRMED' THEN TRUE END) WITH =
				)
		);
	`)
	return err
}

type db interface {
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type StaffStore struct {
	db db
}

func NewStaffStore(db db) *StaffStore {
	return &StaffStore{db: db}
}

func (dep *StaffStore) SaveStaff(ctx context.Context, r *Staff) error {
	if r == nil {
		return errors.New("staff object is nil")
	}

	q := `
    	INSERT INTO staff (email)
        VALUES ($1, $2, $3)
        RETURNING staff_id, email
	`

	row := dep.db.QueryRowContext(ctx, q, r.Email)
	if err := row.Scan(&r.StaffId, &r.Email); err != nil {
		log.Print(err.Error())
		return errors.New("exception saving to staff table")
	}

	return nil
}

func (dep *StaffStore) SaveReservation(ctx context.Context, r *Reservation) error {
	if r == nil {
		return errors.New("reservation cannot be nil")
	}

	q := `
        INSERT INTO reservation (staff_id, name, email, status, created_at, scheduled_for, expire_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING reservation_id, staff_id, name, email, status, created_at, scheduled_for, expire_at;
    `

	row := dep.db.QueryRowContext(
		ctx, q, r.StaffId, r.Name, r.Email, r.Status, r.CreatedAt, r.ScheduledFor, r.ExpireAt)

	err := row.Scan(
		&r.ReservationId, &r.StaffId, &r.Name, &r.Email, &r.Status, &r.CreatedAt, &r.ScheduledFor, &r.ExpireAt)

	if err != nil {
		log.Print(err.Error())
		return errors.New("error saving reservation")
	}

	log.Print("reservation saved")
	return nil
}
