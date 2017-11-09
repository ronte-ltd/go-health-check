package example

/*
package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/ronte-ltd/go-health-check/checkers"
)

const (
	addr  = "localhost:11911"
	route = "/health"
)

func main() {
	fmt.Println("Start.")
	innerServer := checkers.MockHTTPServer("22922")
	//mc, _ := memcache.New("127.0.0.1:11211")
	//defer mc.Close()

	//DB, _ := sql.Open("postgres", "postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable")
	//defer DB.Close()

	checker := checkers.NewHealthChecker("MyChecker")
	checker.RegistryURL("Yandex", "http://localhost:22922/22922")
	//checker.RegistryDB("Postgres", DB)
	//checker.RegistryMemcacheed("memcached", mc)
	checker.RegistryFunc("InnerMyChecker", myChecker)
	go func() {
		fmt.Println("Error:", checker.Handle("localhost:11911", "/health").Error())
	}()

	ticker := time.NewTicker(time.Second * 1)
	go func() {
		for t := range ticker.C {
			fmt.Println("time ", t)
			resp, err := http.Get("http://" + addr + route)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
			}
			bytes, _ := ioutil.ReadAll(resp.Body)
			fmt.Printf("Health: %+v\n", string(bytes))
			resp.Body.Close()
		}
	}()

	time.Sleep(time.Second * 100)
	fmt.Println("Stop.")
	innerServer.Shutdown(context.Background())
}

var nextInt = intSeq()

func myChecker() checkers.Health {
	next := nextInt()
	if next < 30 {
		return checkers.Health{
			Name:   "MyChecker",
			Status: checkers.UP,
			Msg:    "All ok",
		}
	}
	return checkers.Health{
		Name:   "MyChecker",
		Status: checkers.DOWN,
		Msg:    fmt.Sprintf("My checker DOWN, because 29 less than %d", next),
	}
}

func intSeq() func() int {
	i := 0
	return func() int {
		i++
		return i
	}
}
*/
