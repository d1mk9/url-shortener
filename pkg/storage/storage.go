package storage

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	reform "gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"
)

type DBs struct {
	SQL    *sql.DB
	Reform *reform.DB
}

func MustInitPostgres(dsn string) *DBs {
	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("postgres open: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("postgres ping: %v", err)
	}

	logger := log.New(os.Stdout, "SQL: ", log.LstdFlags)
	reformDB := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(logger.Printf))

	return &DBs{SQL: sqlDB, Reform: reformDB}
}
