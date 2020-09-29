package db

import (
	"Proxy/src/models"
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
	"os"
)

type Database struct {
	dbConn *pgx.Conn
}

func CreateNewDatabaseConnection() (*Database, error) {
	username, _ := os.LookupEnv("DB_USER")
	dbName, _ := os.LookupEnv("DB_NAME")
	dbConf := pgx.ConnConfig{
		User:                 username,
		Database:             dbName,
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
		"INSERT INTO requests VALUES ($1, $2)",
		dbReq.Host,
		dbReq.Request,
	)

	if err != nil {
		logrus.Warn("Can't save request to database")
		logrus.Error(err)
	}
}

func (db *Database) GetRequestList() ([]models.DatabaseReq, error) {
	rows, err := db.dbConn.Query(
		"SELECT * FROM requests",
	)

	if err != nil {
		return nil, err
	}

	requests := make([]models.DatabaseReq, 0)
	for rows.Next() {
		req := models.DatabaseReq{}
		err := rows.Scan(&req.Host, &req.Request, &req.Id)
		if err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}

	rows.Close()
	return requests, nil
}

func (db *Database) GetReqById(id int) models.DatabaseReq {
	req := models.DatabaseReq{}
	row := db.dbConn.QueryRow(
		"SELECT * FROM requests WHERE id=$1", id)
	err := row.Scan(&req.Host, &req.Request, &req.Id)
	if err != nil {
		logrus.Warn("Can't get data from database")
		logrus.Error(err)
	}

	return req
}

func (db *Database) Close() {
	err := db.dbConn.Close()
	if err != nil {
		logrus.Warn(err)
	}
}
