package postgres

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type DbConfig struct {
	Host               string
	Port               uint16
	User               string
	Password           string
	DbName             string
	SSLMode            string
	MaxOpenConnections uint16
	MaxIdleConnections uint16
	QueryTimeout       time.Duration
}

type Database interface {
	GetInstance() *gorm.DB
	GetConfig() DbConfig
	Connect()
	Disconnect()
}

type database struct {
	*sql.DB
	instance *gorm.DB
	context  context.Context
	config   DbConfig
}

func NewDatabase(ctx context.Context, config DbConfig) Database {
	db := database{
		context: ctx,
		config:  config,
	}
	return &db
}

func (db *database) GetInstance() *gorm.DB {
	return db.instance
}

func (db *database) GetConfig() DbConfig {
	return db.config
}

func (db *database) Connect() {
	uri := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		db.config.Host,
		db.config.Port,
		db.config.User,
		db.config.Password,
		db.config.DbName,
		db.config.SSLMode,
	)

	postgresDb, err := sql.Open("postgres", uri)
	if err != nil {
		log.Fatal("connection to postgres failed: ", err)
	}
	postgresDb.SetMaxIdleConns(int(db.config.MaxIdleConnections))
	postgresDb.SetMaxOpenConns(int(db.config.MaxOpenConnections))

	err = postgresDb.Ping()
	if err != nil {
		log.Fatalln("pinging to postgres failed: ", err)
	}

	gormDB, err := gorm.Open(
		postgres.New(postgres.Config{Conn: postgresDb}),
		&gorm.Config{},
	)
	if err != nil {
		log.Fatalln("gorm connection to postgres failed: ", err)
	}

	log.Info("successfully connected to postgres database")
	db.instance = gormDB
}

func (db *database) Disconnect() {
	log.Info("disconnecting postgres...")
	sqlDB, err := db.instance.DB()
	if err != nil {
		log.Panic(err)
	}
	if err := sqlDB.Close(); err != nil {
		log.Panic(err)
	}
	log.Info("successfully disconnected postgres")
}
