package helpers

import (
	"encoding/json"
	"net/http"
)

func Get_email_password(r *http.Request) (string, string, error) {
	type message struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	msg := message{}
	err := decoder.Decode(&msg)
	if err != nil {
		return "", "", err
	}
	return msg.Email, msg.Password, nil
}
