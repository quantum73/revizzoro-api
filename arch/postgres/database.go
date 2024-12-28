package postgres

import (
	"context"
	"fmt"
	_ "github.com/lib/pq"
	dishes "github.com/quantum73/revizzoro-api/api/dishes/model"
	restaurants "github.com/quantum73/revizzoro-api/api/restaurants/model"
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
	instance *gorm.DB
	context  context.Context
	config   DbConfig
}

func NewDatabase(ctx context.Context, config DbConfig) Database {
	db := database{context: ctx, config: config}
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

	postgresDb, err := gorm.Open(postgres.Open(uri), &gorm.Config{})
	if err != nil {
		log.Fatal("connection to postgres failed: ", err)
	}

	sqlDB, err := postgresDb.DB()
	if err != nil {
		log.Fatal("connection to postgres failed: ", err)
	}
	sqlDB.SetMaxIdleConns(int(db.config.MaxIdleConnections))
	sqlDB.SetMaxOpenConns(int(db.config.MaxOpenConnections))

	err = sqlDB.Ping()
	if err != nil {
		log.Fatalln("pinging to postgres failed: ", err)
	}

	log.Info("successfully connected to postgres database")
	db.instance = postgresDb

	err = db.instance.AutoMigrate(&restaurants.Restaurant{}, &dishes.Dish{})
	if err != nil {
		log.Fatalln("migrations has failed: ", err)
	}
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
