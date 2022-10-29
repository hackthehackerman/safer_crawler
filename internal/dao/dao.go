package dao

import (
	"database/sql"
	"log"

	"carrierleads.com/internal/dao/carrierleads"
	_ "github.com/go-sql-driver/mysql"
)

type Dao struct {
	DB      *sql.DB
	Queries *carrierleads.Queries
}

func Instance(s string) (d Dao) {
	db, err := sql.Open("mysql", s)
	if err != nil {
		log.Fatal("failed to open database", err)
	}

	return Dao{
		DB:      db,
		Queries: carrierleads.New(db),
	}
}
