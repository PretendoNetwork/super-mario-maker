package database

import (
	"database/sql"
	"os"

	"github.com/PretendoNetwork/super-mario-maker/globals"
	_ "github.com/lib/pq"
)

var Postgres *sql.DB

func ConnectPostgres() {
	var err error

	Postgres, err = sql.Open("postgres", os.Getenv("PN_SMM_POSTGRES_URI"))
	if err != nil {
		globals.Logger.Critical(err.Error())
	}

	globals.Logger.Success("Connected to Postgres!")

	initPostgres()
}
