package db

import (
	"Proxy/models"
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
)

type Database struct {
	dbConn *pgx.Conn
}

func CreateNewDatabaseConnection() (*Database, error) {
	dbConf := pgx.ConnConfig{
		User:                 "farcoad",
		Database:             "proxy",
		PreferSimpleProtocol: false,
	}

	dbConn, err := pgx.Connect(dbConf)
	if err != nil {
		return nil, err
	}
	logrus.Debug("Successfully connected to db")

	return &Database{
		dbConn: dbConn,
	}, nil
}

func (db *Database) InsertRequest(dbReq models.DatabaseReq) {
	_, err := db.dbConn.Exec(
		"INSERT INTO requests VALUES ($1, $2, $3)",
		dbReq.Host,
		dbReq.IsHttps,
		dbReq.Request,
	)

	if err != nil {
		logrus.Warn("Can't save request to database")
	}
}

func (db *Database) Close() {
	err := db.dbConn.Close()
	if err != nil {
		logrus.Warn(err)
	}
}
