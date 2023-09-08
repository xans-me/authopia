package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // postgre library
	log "github.com/sirupsen/logrus"
	"github.com/xans-me/authopia/core/configuration"
)

// NewPostgres constructor
func NewPostgres(conf *configuration.AppConfig) *sql.DB {
	connInfo := fmt.Sprintf(`%s:%s@tcp(%s:%s)/%s`, conf.SQLDatabase.User, conf.SQLDatabase.Password,
		conf.SQLDatabase.Host, conf.SQLDatabase.Port, conf.SQLDatabase.Name)
	switch conf.SQLDatabase.Connection {
	case "postgres":
		connInfo = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", conf.SQLDatabase.Host, conf.SQLDatabase.Port, conf.SQLDatabase.User, conf.SQLDatabase.Name, conf.SQLDatabase.Password)
		break
	}

	db, err := sql.Open(conf.SQLDatabase.Connection, connInfo)
	if err != nil {
		log.Panic(err)
	}

	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}
	db.SetMaxOpenConns(conf.SQLDatabase.MaximumOpenConnection)
	db.SetMaxIdleConns(conf.SQLDatabase.MaximumIdleConnection)

	log.Info("postgres : Ready")
	return db
}
