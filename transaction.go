package main

import (
	"database/sql"
	"errors"
	"log"
)

type Adapters struct {
	StaffStore *StaffStore
	TxProvider *transactionProvider
}

func NewAdapters(db db, tx *transactionProvider) *Adapters {
	return &Adapters{StaffStore: NewStaffStore(db), TxProvider: tx}
}

type transactionProvider struct {
	db *sql.DB
}

func (p *transactionProvider) RunInTransaction(txFunc func(*Adapters) error) error {
	log.Printf("transaction beginning isolation level options")

	err := p.runInTx(p.db, func(tx *sql.Tx) error { return txFunc(NewAdapters(tx, nil)) })
	if err != nil {
		log.Printf("transaction not committed %s", err.Error())
		return err
	}

	log.Print("transaction commited successfully")
	return nil
}

func (p *transactionProvider) runInTx(db *sql.DB, fn func(*sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		log.Printf("error starting transaction %s", err.Error())
		return err
	}

	if err = fn(tx); err == nil {
		log.Print("committing transaction")
		return tx.Commit()
	}

	if rollbackErr := tx.Rollback(); rollbackErr != nil {
		log.Printf("failed to rollback transaction original %s, rollback error  %s", err.Error(), rollbackErr.Error())
		return errors.Join(err, rollbackErr)
	}

	log.Printf("transaction rolled back due to error %s", err.Error())
	return err
}
