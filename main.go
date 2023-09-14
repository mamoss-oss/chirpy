package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	conf := apiConfig{fileserverHits: 0}
	// create new mux (chi router)
	r := chi.NewRouter()

	// second router mux for api
	r_api := chi.NewRouter()

	// admin namespace router
	r_admin := chi.NewRouter()

	// mount api and admin router mux behind main router r
	r.Mount("/api/", r_api)
	r.Mount("/admin/", r_admin)
	// add middleware to mux
	corsRouter := middlewareCors(r)

	// create fileserve http handler at root /
	fs := http.FileServer(http.Dir("./"))

	// Middle ware wrapping done extensive. Can be made much more compact.
	// StripPrefix is a built in middleware. It takes the handler, does something, returns the handler.
	fs_wrapped_one := http.StripPrefix("/app", fs)
	// takes the handler, does something, returns the handler
	fs_wrapped_two := conf.middlewareMetricsInc(fs_wrapped_one)

	// You can either use Handle or HandleFunc. The result is basically the same.
	// It all depends if you got a http.handler or a func that matches the
	// implementation of a handlefunc

	//register handler with different mux routers
	r.Handle("/app", fs_wrapped_two)
	r.Handle("/app/*", fs_wrapped_two)
	r.HandleFunc("/reset", conf.resetMetrics)

	r_api.Get("/healthz", healthz)
	r_api.Get("/chirps", API_GetChirps)
	r_api.Post("/chirps", API_PostChirp)
	r_api.Get("/chirps/{id}", API_Get_Single_Chirp)
	r_api.Post("/users", API_Create_User)
	r_api.Post("/login", API_User_Login)

	r_admin.Get("/metrics", conf.displayMetrics)

	serv := &http.Server{
		Handler: corsRouter,
		Addr:    "0.0.0.0:8080",
	}
	log.Print("Starting server on :8080")
	err := serv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}

}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code >= 500 {
		fmt.Printf("Server side error: %d", code)
	}
	type err_resp struct {
		Err_msg string `json:"error"`
	}
	respondWithJSON(w, code, err_resp{
		Err_msg: msg,
	})
}
