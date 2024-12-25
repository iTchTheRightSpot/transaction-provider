package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"
)

func main() {
	con := "postgres://pattern:pattern@localhost:5432/pattern_db?sslmode=disable"
	db, err := ConnectToPostgres(con)
	if err != nil {
		log.Fatal(err)
		return
	}

	mux := http.NewServeMux()
	h := &ReservationHandler{db: db, mux: mux}
	h.Register()

	server := http.Server{
		Addr:              ":8080",
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       5 * time.Second,
		Handler:           mux,
	}

	log.Fatalf("failed to start server %s", server.ListenAndServe())
}

// ConnectToPostgres https://go.dev/doc/tutorial/database-access
func ConnectToPostgres(conn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// https://www.alexedwards.net/blog/configuring-sqldb
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err = db.Ping(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	log.Println("database connection established")
	return db, nil
}
