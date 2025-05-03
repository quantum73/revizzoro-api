package config

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

type ServerModeEnum string

const (
	DEBUG   ServerModeEnum = "debug"
	TEST    ServerModeEnum = "test"
	RELEASE ServerModeEnum = "release"
)

type Env struct {
	// server
	ServerMode            ServerModeEnum `mapstructure:"SERVER_MODE"`
	ServerHost            string         `mapstructure:"SERVER_HOST"`
	ServerPort            uint16         `mapstructure:"SERVER_PORT"`
	ServeWriteTimeout     uint16         `mapstructure:"SERVER_WRITE_TIMEOUT"`
	ServeReadTimeout      uint16         `mapstructure:"SERVER_READ_TIMEOUT"`
	ServerIdleTimeout     uint16         `mapstructure:"SERVER_IDLE_TIMEOUT"`
	ServerGracefulTimeout uint16         `mapstructure:"SERVER_GRACEFUL_TIMEOUT"`
	// database
	DBHType              string `mapstructure:"DB_TYPE"`
	DBHost               string `mapstructure:"DB_HOST"`
	DBPort               uint16 `mapstructure:"DB_PORT"`
	DBUser               string `mapstructure:"DB_USER"`
	DBPassword           string `mapstructure:"DB_PASS"`
	DBName               string `mapstructure:"DB_NAME"`
	DBSSLMode            string `mapstructure:"DB_SSLMODE"`
	DbMaxOpenConnections uint16 `mapstructure:"DB_MAX_OPEN_CONNECTIONS"`
	DbMaxIdleConnections uint16 `mapstructure:"DB_MAX_IDLE_CONNECTIONS"`
	DBQueryTimeout       uint16 `mapstructure:"DB_QUERY_TIMEOUT"`
	DbMigrationsPath     string `mapstructure:"DB_MIGRATIONS_PATH"`
}

func NewEnv(filename string, override bool) *Env {
	env := Env{}
	viper.SetConfigFile(filename)

	if override {
		viper.AutomaticEnv()
	}

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading environment file: %s", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatalf("Error loading environment file: %s", err)
	}

	return &env
}
