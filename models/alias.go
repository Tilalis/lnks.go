package models

import (
	"database/sql"
	"net/url"
	"regexp"
)

// Alias for Urls
type Alias struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	userID int
}

var aliasStmts struct {
	get       *sql.Stmt
	insert    *sql.Stmt
	delete    *sql.Stmt
	getByUser *sql.Stmt
}

var err error

// NewAlias creates alias
func NewAlias(name, url string, user *User) (*Alias, error) {
	alias := &Alias{
		Name: name,
		URL:  url,
	}

	if user != nil {
		alias.userID = user.ID
	}

	err = alias.Validate()

	if err != nil {
		return nil, err
	}

	return alias, nil
}

// GetAlias by its name
func GetAlias(name string) (*Alias, error) {
	if connection == nil {
		return nil, ErrNoConnection
	}

	alias := &Alias{}
	err = aliasStmts.get.QueryRow(name).Scan(&alias.Name, &alias.URL, &alias.userID)

	if err != nil {
		return nil, err
	}

	return alias, err
}

// GetAliases by user
func GetAliases(user *User) ([]Alias, error) {
	if connection == nil {
		return nil, ErrNoConnection
	}

	var rows *sql.Rows

	if user != nil {
		rows, err = aliasStmts.getByUser.Query(user.ID)
	} else {
		rows, err = aliasStmts.getByUser.Query(0)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var alias Alias
	var aliases = make([]Alias, 0)

	for rows.Next() {
		rows.Scan(&alias.Name, &alias.URL, &alias.userID)
		aliases = append(aliases, Alias{
			Name:   alias.Name,
			URL:    alias.URL,
			userID: alias.userID,
		})
	}

	return aliases, nil
}

// SetUser sets userid
func (alias *Alias) SetUser(user *User) {
	if user == nil {
		alias.userID = 0
	} else {
		alias.userID = user.ID
	}
}

// Validate validates name and url fields
func (alias *Alias) Validate() error {
	match, _ := regexp.MatchString("[a-z0-9]+", alias.Name)

	if !match {
		return ErrWrongAlias
	}

	u, err := url.Parse(alias.URL)
	isURL := err == nil && u.Host != ""

	if !isURL {
		return ErrWrongURL
	}

	return nil
}

// Save saves alias if it does not exist yet
func (alias *Alias) Save() error {
	if connection == nil {
		return ErrNoConnection
	}

	_, err = aliasStmts.insert.Exec(alias.Name, alias.URL, alias.userID)

	return err
}

// Delete alias
func (alias *Alias) Delete() error {
	if connection == nil {
		return ErrNoConnection
	}

	_, err = aliasStmts.delete.Exec(alias.Name)

	return err
}

// PrepareAlias prepares statements for Alias
func prepareAlias(connection *sql.DB) error {
	aliasStmts.get, err = connection.Prepare("SELECT name, url, userid FROM `alias` WHERE name = ?")

	if err != nil {
		return err
	}

	aliasStmts.insert, err = connection.Prepare("INSERT INTO `alias` (name, url, userid) VALUES (?, ?, ?)")

	if err != nil {
		return err
	}

	aliasStmts.getByUser, err = connection.Prepare("SELECT name, url, userid FROM `alias` WHERE userid = ?")

	if err != nil {
		return err
	}

	aliasStmts.delete, err = connection.Prepare("DELETE FROM `alias` WHERE name = ?")

	return err
}

// CloseAlias closes statements for Alias
func closeAlias() {
	aliasStmts.get.Close()
	aliasStmts.delete.Close()
	aliasStmts.insert.Close()
	aliasStmts.getByUser.Close()
}
