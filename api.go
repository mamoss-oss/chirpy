package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/mamoss-oss/chirpy/internal/db"
	"github.com/mamoss-oss/chirpy/internal/helpers"
)

func API_PostChirp(w http.ResponseWriter, r *http.Request) {
	db, err := db.NewDB("./database.json")
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	message, err := validate_chirp(r)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	chirp, err := db.CreateChirp(message)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	respondWithJSON(w, 201, chirp)
}

func API_GetChirps(w http.ResponseWriter, r *http.Request) {
	db, err := db.NewDB("./database.json")
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	chirps, err := db.GetChirps()
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	respondWithJSON(w, 200, chirps)
}

func API_Get_Single_Chirp(w http.ResponseWriter, r *http.Request) {
	db, err := db.NewDB("./database.json")
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	dbStructure, err := db.LoadDB()
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	chirp_id_str := chi.URLParam(r, "id")
	chirp_id_int, err := strconv.Atoi(chirp_id_str)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	data, ok := dbStructure.Chirps[chirp_id_int]
	if !ok {
		respondWithError(w, 404, "not found")
		return
	}
	respondWithJSON(w, 200, data)

}

func API_Create_User(w http.ResponseWriter, r *http.Request) {
	db, err := db.NewDB("./database.json")

	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	email, password, err := helpers.Get_email_password(r)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	user, err := db.CreateUser(email, password)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	user_no_pass := struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}{
		ID:    user.ID,
		Email: user.Email,
	}
	respondWithJSON(w, 201, user_no_pass)
}

func API_User_Login(w http.ResponseWriter, r *http.Request) {
	database, err := db.NewDB("./database.json")
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	database_loaded, err := database.LoadDB()
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	email, password, err := helpers.Get_email_password(r)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	user, err := db.FindByMail(&database_loaded, email)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	password_ok := db.ComparePasswords(user.Password, password)
	if password_ok {
		user_no_pass := struct {
			ID    int    `json:"id"`
			Email string `json:"email"`
		}{
			ID:    user.ID,
			Email: user.Email,
		}
		respondWithJSON(w, 200, user_no_pass)
		return
	} else {
		respondWithError(w, 401, "Unauthorized")
		return
	}

}
