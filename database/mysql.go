package database

import (
	"database/sql"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	"github.com/kelseyhightower/envconfig"
	"log"
)

const defaultAppPrefix = "APP"

// Config describes MySQL config structure
type Config struct {
	User                 string `envconfig:"APP_USER" default:"root"`
	Passwd               string `envconfig:"APP_PASS" default:""`
	Net                  string `envconfig:"APP_NET" default:"tcp"`
	Addr                 string `envconfig:"APP_ADDR" required:"true"`
	DBName               string `envconfig:"APP_DB_NAME" default:"test"`
	AllowNativePasswords bool   `envconfig:"APP_ALLOW_NATIVE" default:"true"`
}

// FromEnv reads env variables
func FromEnv() (*Config, error) {
	config := &Config{}

	err := envconfig.Process(defaultAppPrefix, config)
	if err != nil {
		return config, fmt.Errorf("failed to load config: %w", err)
	}

	return config, nil
}

// Validate validates env variables requirements
func (c *Config) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}

// InitDb prepares MySQL session
func InitDb(conf *Config) *sql.DB {
	cfg := mysql.Config{
		User:                 conf.User,
		Passwd:               conf.Passwd,
		Net:                  conf.Net,
		Addr:                 conf.Addr,
		DBName:               conf.DBName,
		AllowNativePasswords: conf.AllowNativePasswords,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatalln(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		fmt.Println(pingErr)
	}
	return db
}
