package example

/*
package main

import (
	"github.com/rainycape/memcache"
	"github.com/ronte-ltd/go-health-check/checkers"
	"database/sql"
	_ "github.com/lib/pq"
	"time"
	"fmt"
)

func main() {
	fmt.Println("Start.")
	mc, _ := memcache.New("127.0.0.1:11211")
	mcChecker := checkers.NewMemcachedChecker("memcache", mc)
	defer mc.Close()

	http := checkers.NewHttpChecker("Yandex", "https://ya.ru")

	DB, _ := sql.Open("postgres", "postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable")
	defer DB.Close()
	dbChecker := checkers.NewDBChecker("Postgres", DB)

	funcChecker := checkers.NewFuncChecker("MyChecker", MyChecker)
	composite := checkers.SimpleChecker("CompositeChecker")
	composite.AddChecker(&mcChecker)
	composite.AddChecker(&http)
	composite.AddChecker(&funcChecker)
	composite.AddChecker(&dbChecker)

	handler := checkers.NewHandler("localhost:11911")
	go handler.AddRoute("/health", &composite)



	ticker := time.NewTicker(time.Second * 1)
	go func() {
		for t := range ticker.C {
			fmt.Println("time ", t)
			health, _ := composite.Check()
			fmt.Printf("Health: %+v\n", health)
		}
	}()

	time.Sleep(time.Second * 100)
	fmt.Println("Stop.")
}

var nextInt = intSeq()

func MyChecker() checkers.Health {
	next := nextInt()
	if  next < 30 {
		return checkers.Health{
			Name:   "MyChecker",
			Status: checkers.UP,
			Msg:    "All ok",
		}
	}else {
		return checkers.Health{
			Name:   "MyChecker",
			Status: checkers.DOWN,
			Msg:    fmt.Sprintf("My checker DOWN, because 29 less than %d", next),
		}
	}
}

func intSeq() func() int {
	i := 0
	return func() int {
		i += 1
		return i
	}
}
*/
