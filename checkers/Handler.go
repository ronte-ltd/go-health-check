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
	newInnerHandler(route)
	return Handler{
		Client: http.Client{},
		Server: http.Server{Addr: addr},
	}
}

func newInnerHandler(route string) http.Handler {
	handler := new(innerHandler)
	http.HandleFunc(route, handler.ServeHTTP)
	return handler
}

type innerHandler func(http.ResponseWriter, *http.Request)

func (i *innerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		fmt.Println("Error:", err.Error())
	}
}
