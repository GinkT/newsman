package db

import (
	"fmt"
	"github.com/ginkt/newsman/config"
	"github.com/mailru/dbr"
)

func CreatePostgresClient(cfg *config.Config) (dbConn *dbr.Connection, err error) {
	dbConn, err = dbr.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.PostgresUsername, cfg.PostgresPassword, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDatabase),
		&dbr.NullEventReceiver{})
	if err != nil {
		return
	}
	return
}