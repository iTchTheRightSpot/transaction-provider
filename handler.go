package main

import (
	"database/sql"
	"net/http"
)

type ReservationHandler struct {
	mux *http.ServeMux
	db  *sql.DB
}

func (dep *ReservationHandler) Register() {
	dep.mux.HandleFunc("GET /reservation", dep.dates)
	dep.mux.HandleFunc("POST /reservation", dep.create)
}

func (dep *ReservationHandler) dates(w http.ResponseWriter, r *http.Request) {}

func (dep *ReservationHandler) create(w http.ResponseWriter, r *http.Request) {}
