package checkers

import (
	"database/sql"
	"errors"
)

// DBChecker sturct provide health check for SQL DB by `sql.DB`
type DBChecker struct {
	HealthChecker HealthChecker
	QuerySQL      string
	DB            *sql.DB
}

// NewDBChecker return new instance DBChecker with name and DB resource args
func NewDBChecker(name string, DB *sql.DB) *DBChecker {
	return &DBChecker{
		HealthChecker: NewHealthChecker(name),
		QuerySQL:      "SELECT 1",
		DB:            DB,
	}
}

// Check health DB service
func (dbc *DBChecker) Check() (Health, error) {
	var result string
	if dbc.DB == nil {
		err := errors.New("empty connection")
		dbc.HealthChecker.Msg = err.Error()
		dbc.HealthChecker.PushHealth()
		return dbc.HealthChecker.Health, err
	}
	err := dbc.DB.QueryRow(dbc.QuerySQL).Scan(&result)
	if err != nil {
		dbc.HealthChecker.Msg = err.Error()
		dbc.HealthChecker.PushHealth()
		return dbc.HealthChecker.Health, err
	}
	dbc.HealthChecker.Up()
	dbc.HealthChecker.PushHealth()
	return dbc.HealthChecker.Health, nil
}

// Name return name this checker
func (dbc *DBChecker) Name() string {
	return dbc.HealthChecker.Health.Name
}
