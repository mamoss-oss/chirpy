package db

import "sort"

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.LoadDB()
	if err != nil {
		return []Chirp{}, err
	}
	chirps := []Chirp{}

	if dbStructure.Chirps == nil {
		return chirps, err
	}

	for _, val := range dbStructure.Chirps {
		chirps = append(chirps, val)
	}
	return chirps, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.LoadDB()
	if err != nil {
		return Chirp{}, err
	}
	ch := Chirp{}
	next, err := db.nextID_chirp()
	if err != nil {
		UNUSED(next)
		return ch, err
	}
	ch.ID = next
	ch.Body = body
	dbStructure.Chirps[ch.ID] = ch
	db.writeDB(dbStructure)
	return ch, nil
}

func (db *DB) nextID_chirp() (int, error) {
	chirps, err := db.GetChirps()
	if err != nil {
		UNUSED(chirps)
		return -1, err
	}
	if len(chirps) == 0 {
		return 1, nil
	}

	sort.Slice(chirps, func(i, j int) bool { return chirps[i].ID < chirps[j].ID })
	last_id := chirps[len(chirps)-1].ID
	return last_id + 1, nil

}
