package checkers

import (
	"fmt"
	"net/http"
)

//MockHTTPServer run local HTTP server that response 200 for all GET requests
func MockHTTPServer(port string) *http.Server {
	http.HandleFunc(fmt.Sprintf("/%s", port), func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.WriteHeader(http.StatusOK)
		}
	})
	http := &http.Server{Addr: fmt.Sprintf(":%s", port)}
	go http.ListenAndServe()
	return http
}
