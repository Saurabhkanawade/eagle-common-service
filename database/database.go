package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	defaultMaxIdleConnections = 5
	defaultMaxOpenConnections = 50
)

// DbConfig - allows for the configuration of a database connection pool for PostgreSQL
// * Host - database host - required
// * Post - database port - required
// * User - database user name - required
// * Pass - database user password - required
// * DbName - Name of the database to connect to - required
// * MaxIdleConn - the maximum number of idle connection in the pool - optional, if missing will be set to a default
// * MaxOpenConn - the maximum number of open connections - optional, if missing will be set to a default
// * MaxLifetimeMins - the maximum number of mins a connection can be reused before being closed - optional, if missing connections will not
type DbConfig struct {
	Host            string
	Port            int
	User            string
	Pass            string
	DbName          string
	MaxIdleConn     int
	MaxOpenConn     int
	MaxLifetimeMins int
}

var db *sql.DB

type DbConnection struct {
	Conn *sql.DB
}

func (con DbConnection) PingDB() error {
	conn := con.Conn.Ping()

	if conn != nil {
		logrus.Errorf("Error while pinging the database: %v", con)
		return conn
	}

	logrus.Infof("Database ping successful")
	return nil
}

func InitDatabase(config DbConfig) (DbConnection, error) {
	logrus.Info("connecting to the db .....")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.User, config.Pass, config.DbName)

	connection := DbConnection{}

	if db != nil {
		connection.Conn = db
	} else {
		pDb, err := sql.Open("postgres", psqlInfo)

		if err != nil {
			return DbConnection{}, errors.Wrap(err, "error while connecting to the database")
		}

		connection.Conn = pDb
	}

	if config.MaxIdleConn == 0 {
		logrus.Info("InitDatabase:: setting SetMaxIdleConns to %d", defaultMaxIdleConnections)
		connection.Conn.SetMaxIdleConns(defaultMaxIdleConnections)
	} else {
		logrus.Info("InitDatabase:: setting SetMaxIdleConns to %d", config.MaxIdleConn)
		connection.Conn.SetMaxIdleConns(config.MaxIdleConn)
	}

	if config.MaxOpenConn == 0 {
		logrus.Info("InitDatabase:: setting SetMaxOpenConns to %d", defaultMaxOpenConnections)
		connection.Conn.SetMaxOpenConns(defaultMaxOpenConnections)
	} else {
		logrus.Info("InitDatabase:: setting SetMaxOpenConns to %d", config.MaxOpenConn)
		connection.Conn.SetMaxOpenConns(config.MaxOpenConn)
	}

	if config.MaxLifetimeMins != 0 {
		logrus.Info("InitDatabase:: setting ConnMaxLifetime to %d mins", config.MaxLifetimeMins)
		connection.Conn.SetConnMaxLifetime(time.Minute * time.Duration(config.MaxLifetimeMins))
	}
	logrus.Info("connected to the db successfully.....")
	return connection, nil
}
