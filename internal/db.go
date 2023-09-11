package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	_, err := os.ReadFile(path)
	var db = DB{
		path: path,
		mux:  &sync.RWMutex{},
	}
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = db.ensureDB()
			if err != nil {
				return &DB{}, err
			} else {
				return &db, err
			}
		} else {
			return &DB{}, fmt.Errorf("file error %s", err)
		}
	} else {
		return &db, err
	}
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	_, err := os.Create(db.path)
	if err != nil {
		return err
	}
	err = db.writeDB(DBStructure{})
	if err != nil {
		return err
	}
	return nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {

	content, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}
	db.mux.Lock()
	defer db.mux.Unlock()
	err = os.WriteFile(db.path, content, 0644)
	if err != nil {
		return err
	}
	return nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	data, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}

	var dbStructure DBStructure

	err = json.Unmarshal(data, &dbStructure)

	if err != nil {
		return DBStructure{}, err
	}
	return dbStructure, nil

}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return []Chirp{}, err
	}
	chirps := []Chirp{}

	if dbStructure.Chirps == nil {
		return chirps, err
	}

	// todo after writing the add chirp function

}
