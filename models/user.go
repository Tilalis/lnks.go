package models

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
)

// User struct represents user
type User struct {
	ID       int
	Username string
	Hash     string
}

var userStmts struct {
	getByUsername *sql.Stmt
	insert        *sql.Stmt
	delete        *sql.Stmt
}

// NewUser creates new user
func NewUser(username, password string) (*User, error) {
	if username == "" || password == "" {
		return nil, ErrEmptyUserCredentials
	}

	hash := sha256.Sum256([]byte(password))
	hashString := fmt.Sprintf("%x", hash)

	user := &User{
		Username: username,
		Hash:     hashString,
	}

	return user, nil
}

// GetUser gets user from database
func GetUser(username string) (*User, error) {
	if connection == nil {
		return nil, ErrNoConnection
	}

	user := &User{}
	err = userStmts.getByUsername.QueryRow(username).Scan(&user.ID, &user.Username, &user.Hash)

	return user, err
}

// AuthenticateUser gets user from database and checks password
func AuthenticateUser(username, password string) (*User, error) {
	if connection == nil {
		return nil, ErrNoConnection
	}

	hash := sha256.Sum256([]byte(password))
	hashString := fmt.Sprintf("%x", hash)

	user := &User{}

	err = userStmts.getByUsername.QueryRow(username).Scan(&user.ID, &user.Username, &user.Hash)

	if err != nil {
		return nil, err
	}

	if user.Hash != hashString {
		return nil, ErrWrongUserPassword
	}

	return user, err
}

// Save saves user to database
func (user *User) Save() error {
	if connection == nil {
		return ErrNoConnection
	}

	result, err := userStmts.insert.Exec(user.Username, user.Hash)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	if err == nil {
		user.ID = int(id)
	} else {
		err = userStmts.getByUsername.QueryRow(user.Username).Scan(&user.ID, &user.Username, &user.Hash)

		if err != nil {
			return err
		}
	}

	return err
}

// Delete deletes user from database
func (user *User) Delete() error {
	if connection == nil {
		return ErrNoConnection
	}

	_, err = userStmts.delete.Exec(user.Username)
	return err
}

func prepareUser(connection *sql.DB) error {
	userStmts.getByUsername, err = connection.Prepare("SELECT id, username, hash FROM user WHERE username = ?")

	if err != nil {
		return err
	}

	userStmts.insert, err = connection.Prepare("INSERT INTO user (username, hash) VALUES (?, ?)")

	if err != nil {
		return err
	}

	userStmts.delete, err = connection.Prepare("DELETE FROM user WHERE username = ?")

	return err
}

func closeUser() {
	userStmts.getByUsername.Close()
	userStmts.insert.Close()
	userStmts.delete.Close()
}
