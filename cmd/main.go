package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var ip = flag.String("addr", "127.0.0.1:8001", "Serving host and port")
var user_server = flag.String("user_service", "127.0.0.1:", "User microservice endpoint")
var apiRoot = flag.String("api_root", "/v1", "api root path")

func loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Inside register handler")
	return
}

func main() {
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc(*apiRoot+"/user/login", loginHandler).Methods("GET")
	srv := &http.Server{
		Handler: r,
		Addr:    *ip,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Printf("Running server on %s/user/login.\n", *ip)
	log.Fatal(srv.ListenAndServe())
}
