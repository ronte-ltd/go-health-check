package checkers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	handlerChecker Checker
)

// Handler struct provide httpClient and listen addr for HTTP request via current health
type Handler struct {
	Client http.Client
	Server http.Server
}

// NewHandler return new handler with address
func NewHandler(checker Checker, addr, route string) Handler {
	handlerChecker = checker
	http.HandleFunc(route, defaultHandler)
	return Handler{
		Client: http.Client{},
		Server: http.Server{Addr: addr},
	}
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	health, err := handlerChecker.Check()
	checkError(err)

	bytes, err := json.Marshal(health)
	checkError(err)
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func checkError(err error) {
	if err != nil {
		msg := fmt.Sprintf("Was error: %s", err.Error())
		fmt.Println(msg)
	}
}
