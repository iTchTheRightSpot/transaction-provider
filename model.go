package main

import "time"

type Staff struct {
	StaffId uint64 `json:"staff_id"`
	Email   string `json:"email"`
}

type ReservationEnum string

const (
	CONFIRMED ReservationEnum = "CONFIRMED"
	CANCELLED ReservationEnum = "CANCELLED"
)

type Reservation struct {
	ReservationId uint64          `json:"reservation_id"`
	StaffId       uint64          `json:"staff_id"`
	Name          string          `json:"name"`
	Email         string          `json:"email"`
	Status        ReservationEnum `json:"status"`
	CreatedAt     time.Time       `json:"created_at"`
	ScheduledFor  time.Time       `json:"scheduled_for"`
	ExpireAt      time.Time       `json:"expire_at"`
}
