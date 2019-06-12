package models

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
)

// User struct represents user
type User struct {
	ID       int
	Username string
	Hash     string
}

var userStmts struct {
	get    *sql.Stmt
	insert *sql.Stmt
	delete *sql.Stmt
}

// NewUser creates new user
func NewUser(username, password string) *User {
	hash := sha256.Sum256([]byte(password))
	hashString := fmt.Sprintf("%x", hash)

	user := &User{
		Username: username,
		Hash:     hashString,
	}

	return user
}

// GetUser gets user from database
func GetUser(username, password string) (*User, error) {
	if connection == nil {
		return nil, ErrNoConnection
	}

	hash := sha256.Sum256([]byte(password))
	hashString := fmt.Sprintf("%x", hash)

	user := &User{}

	err = userStmts.get.QueryRow(username).Scan(&user.ID, &user.Username, &user.Hash)

	if user.Hash != hashString {
		return nil, errors.New("")
	}

	return user, err
}

// Save saves user to database
func (user *User) Save() error {
	if connection == nil {
		return ErrNoConnection
	}

	_, err = userStmts.insert.Exec(user.Username, user.Hash)
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
	userStmts.get, err = connection.Prepare("SELECT id, username, hash FROM user WHERE username = ?")

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
	userStmts.get.Close()
	userStmts.insert.Close()
	userStmts.delete.Close()
}
