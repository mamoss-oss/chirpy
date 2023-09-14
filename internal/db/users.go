package db

import (
	"errors"
	"slices"

	"golang.org/x/crypto/bcrypt"
)

func (db *DB) CreateUser(email string, password string) (User, error) {
	dbStructure, err := db.LoadDB()
	if err != nil {
		return User{}, err
	}

	emails := GenerateEmailSlice(&dbStructure)
	exists := checkEmailExists(&emails, email)
	if exists {
		return User{}, errors.New("email already registered")
	}
	user := User{}

	id := len(dbStructure.Users) + 1
	user.ID = id
	user.Email = email
	hashed_pw, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return User{}, err
	}
	user.Password = string(hashed_pw)

	dbStructure.Users[user.ID] = user
	db.writeDB(dbStructure)
	return user, nil
}

func GenerateEmailSlice(users *DBStructure) []string {
	var emails []string
	for _, user := range users.Users {
		emails = append(emails, user.Email)
	}
	return emails
}

func checkEmailExists(emails *[]string, email string) bool {
	return slices.Contains(*emails, email)
}

func FindByMail(users *DBStructure, email string) (User, error) {
	for key, value := range users.Users {
		if value.Email == email {
			return users.Users[key], nil
		}
	}
	return User{}, errors.New("user email not found")
}

func ComparePasswords(hashed string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}
