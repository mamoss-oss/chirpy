package main

import (
	"fmt"
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) resetMetrics(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	w.WriteHeader(200)
}

func (cfg *apiConfig) displayMetrics(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits)))
}

func main() {
	conf := apiConfig{fileserverHits: 0}
	// create new mux
	mux := http.NewServeMux()
	// create fileserve http handler at root /
	fs := http.FileServer(http.Dir("./"))

	// Middle ware wrapping done extensive. Can be made much more compact.
	// StripPrefix is a built in middleware. It takes the handler, does something, returns the handler.
	fs_wrapped_one := http.StripPrefix("/app/", fs)
	// takes the handler, does something, returns the handler
	fs_wrapped_two := conf.middlewareMetricsInc(fs_wrapped_one)

	// You can either use Handle or HandleFunc. The result is basically the same.
	// It all depends if you got a http.handler or a func that matches the
	// implementation of a handlefunc

	//register handler with mux
	mux.Handle("/app/", fs_wrapped_two)

	// register healthz function with the mux
	mux.HandleFunc("/healthz", healthz)

	mux.HandleFunc("/metrics", conf.displayMetrics)

	mux.HandleFunc("/reset", conf.resetMetrics)

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
