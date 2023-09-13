package db

func (db *DB) CreateUser(email string, password string) (User, error) {
	dbStructure, err := db.LoadDB()
	if err != nil {
		return User{}, err
	}
	user := User{}

	id := len(dbStructure.Users) + 1
	user.ID = id
	user.Email = email
	user.Password = password
	dbStructure.Users[user.ID] = user
	db.writeDB(dbStructure)
	return user, nil
}
