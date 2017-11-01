package checkers

import (
	"database/sql"
	"fmt"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

func IgnoreNewDBChecker(t *testing.T) {
	DB, err := sql.Open("postgres", "postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable")
	defer DB.Close()
	if err != nil {
		t.Fatalf("Some error when open connection: %s", err.Error())
	}
	var dbChecker = NewDBChecker("DataBase", DB)
	health, err := dbChecker.Check()
	if err != nil {
		t.Fatalf("Unhealthy: %s", err.Error())
	}
	t.Log(fmt.Sprintf("Health: %+v", health))
}

func TestStatusUp(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Cannot create mock: %s", err.Error())
	}
	mock.ExpectQuery("SELECT 1").WillReturnRows(
		sqlmock.NewRows([]string{"?"}).
			AddRow("1"))

	var dbChecker = NewDBChecker("DataBase", db)
	health, err := dbChecker.Check()
	if err != nil {
		t.Fatalf("Unhealthy: %s", err.Error())
	}
	t.Log(fmt.Sprintf("Health: %+v", health))
}
