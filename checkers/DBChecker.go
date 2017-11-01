package checkers

import (
	"database/sql"
	"errors"
)

type DBChecker struct {
	HealthChecker HealthChecker
	QuerySQL      string
	DB            *sql.DB
}

func NewDBChecker(name string, DB *sql.DB) DBChecker {
	return DBChecker{
		HealthChecker: NewHealthChecker(name),
		QuerySQL:      "SELECT 1",
		DB:            DB,
	}
}

func (dbc *DBChecker) Check() (Health, error) {
	var result string
	if dbc.DB == nil {
		err := errors.New("empty connection")
		dbc.HealthChecker.Msg = err.Error()
		return dbc.HealthChecker.Health, err
	}
	err := dbc.DB.QueryRow(dbc.QuerySQL).Scan(&result)
	if err != nil {
		dbc.HealthChecker.Msg = err.Error()
		return dbc.HealthChecker.Health, err
	}
	dbc.HealthChecker.Up()
	return dbc.HealthChecker.Health, nil
}

func (dbc *DBChecker) Name() string {
	return dbc.HealthChecker.Name
}
