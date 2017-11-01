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

func (c *DBChecker) Check() (Health, error) {
	var result string
	if c.DB == nil {
		err := errors.New("empty connection")
		c.HealthChecker.Msg = err.Error()
		return c.HealthChecker.Health, err
	}
	err := c.DB.QueryRow(c.QuerySQL).Scan(&result)
	if err != nil {
		c.HealthChecker.Msg = err.Error()
		return c.HealthChecker.Health, err
	}
	c.HealthChecker.Up()
	return c.HealthChecker.Health, nil
}
