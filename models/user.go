package models

import (
	"database/sql"
	"log"
)

// User model
type User struct {
	ID        int64
	Email     string
	Password  string
	Address   string
	Telephone string
	Name      string
	Picture   string
}

// ProfileData  model for view page
type ProfileData struct {
	Error string
	User  User
}

// CreateUser insert user data to db
func CreateUser(db *sql.DB, user User) error {
	ins, err := db.Prepare("INSERT INTO USERS(email, password, address, telephone, full_name) values(?,?,?,?,?)")

	if err != nil {
		log.Panic(err.Error())
		return err
	}

	_, err = ins.Exec(user.Email, user.Password, user.Address, user.Telephone, user.Name)
	log.Println("insert user email: " + user.Email)

	return err
}

// UpdateUser update user tables
func UpdateUser(db *sql.DB, user User) error {
	ins, err := db.Prepare("UPDATE users set address=?, telephone=?, full_name=? where email=?")
	if err != nil {
		log.Fatal(err)
		return err
	}

	_, err = ins.Exec(user.Address, user.Telephone, user.Name, user.Email)
	log.Println("user with email " + user.Email + " data updated")

	return err
}

// FindUser find user by email
func FindUser(db *sql.DB, email string) (User, error) {
	var eml, password, address, telephone, fullname, imageURL sql.NullString
	user := User{}
	err := db.QueryRow("SELECT email, password, address, telephone, full_name, image_url FROM users WHERE email=? limit 1", email).
		Scan(&eml, &password, &address, &telephone, &fullname, &imageURL)

	user.Email = eml.String
	user.Password = password.String
	user.Address = address.String
	user.Telephone = telephone.String
	user.Name = fullname.String
	user.Picture = imageURL.String

	return user, err
}
