package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"

	log "github.com/sirupsen/logrus"
)

func ConnectToPostgres(
	host string, port int, user, password, dbname, sslmode string,
) (*sql.DB, error) {
	pgInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	log.Infof("Try to connect to `%s` database on %s:%d", dbname, host, port)
	db, err := sql.Open("postgres", pgInfo)
	if err != nil {
		return nil, err
	}
	log.Info("Ping to database")
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Infof("Successfully connected to `%s` database on %s:%d", dbname, host, port)
	return db, nil
}
