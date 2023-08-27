package main

import (
	"log"
	"net/http"
)

// type rootHandle struct{}

// func (rootHandle) ServeHTTP(http.ResponseWriter, *http.Request) {}

func main() {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./"))
	mux.Handle("/", fs)

	// add middleware to mux
	corsMux := middlewareCors(mux)

	serv := &http.Server{
		Handler: corsMux,
		Addr:    "0.0.0.0:8080",
	}

	err := serv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}

}

// middlewareCors required to allow browser interaction for bootdev course.
func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
