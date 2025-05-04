package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"time"
)

type DbConfig struct {
	Type               string
	Host               string
	Port               uint16
	User               string
	Password           string
	DbName             string
	SSLMode            string
	MaxOpenConnections uint16
	MaxIdleConnections uint16
	QueryTimeout       time.Duration
	MigrationsPath     string
}

type Database interface {
	GetInstance() *sql.DB
	Ping() error
	Connect() error
	Disconnect() error
}

type database struct {
	instance *sql.DB
	config   DbConfig
}

func NewDatabase(config DbConfig) Database {
	return &database{config: config}
}

func (db *database) GetInstance() *sql.DB {
	return db.instance
}

func (db *database) Ping() error {
	return db.instance.Ping()
}

func (db *database) Connect() error {
	const op = "[internal/postgres/database Connect]"

	uri := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		db.config.Host,
		db.config.Port,
		db.config.User,
		db.config.Password,
		db.config.DbName,
		db.config.SSLMode,
	)

	dbObj, err := sql.Open("postgres", uri)
	if err != nil {
		return err
	}

	dbObj.SetMaxIdleConns(int(db.config.MaxIdleConnections))
	dbObj.SetMaxOpenConns(int(db.config.MaxOpenConnections))

	err = dbObj.Ping()
	if err != nil {
		return err
	}

	// "postgres://user:password@localhost:5432/db_name?sslmode=disable"
	dbUriForMigrations := fmt.Sprintf(
		"%s://%s:%s@%s:%d/%s?sslmode=%s",
		db.config.Type,
		db.config.User,
		db.config.Password,
		db.config.Host,
		db.config.Port,
		db.config.DbName,
		db.config.SSLMode,
	)
	err = db.MigrateUp(db.config.MigrationsPath, dbUriForMigrations)
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Warnf("%s Error running migrations: %v\n", op, err)
		return err
	}

	log.Infoln("Successfully connected to database!")
	db.instance = dbObj
	return nil
}

func (db *database) Disconnect() error {
	log.Infoln("Disconnecting database...")
	if err := db.instance.Close(); err != nil {
		return err
	}
	log.Infoln("Successfully disconnected database!")
	db.instance = nil
	return nil
}

func (db *database) MigrateUp(migrationPath string, dbUri string) error {
	m, err := migrate.New(migrationPath, dbUri)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		return err
	}
	return nil
}
