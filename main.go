package main

import (
	"log"
	"net/http"
)

// type rootHandle struct{}

// func (rootHandle) ServeHTTP(http.ResponseWriter, *http.Request) {}

func main() {
	// create new mux
	mux := http.NewServeMux()
	// create fileserve http handler at root /
	fs := http.FileServer(http.Dir("./"))
	//register handler with mux
	mux.Handle("/app/", http.StripPrefix("/app/", fs))

	// register healthz function with the mux
	mux.HandleFunc("/healthz", healthz)

	// add middleware to mux
	corsMux := middlewareCors(mux)

	serv := &http.Server{
		Handler: corsMux,
		Addr:    "0.0.0.0:8080",
	}
	log.Print("Starting server on :8080")
	err := serv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}

}
