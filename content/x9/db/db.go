package db

import (
	"crypto/rand"     // for random generation
	"database/sql"    // generic go database package
	"encoding/base64" // for random generation
	"errors"          // for creating return errors
	"fmt"             // for formatted printing

	_ "github.com/mattn/go-sqlite3" // driver for sqlite
)

var db *sql.DB

// ConnectServer creates a connection.
func ConnectServer(dbname string) {
	var err error
	db, err = sql.Open("sqlite3", dbname) // different open parameters
	if err != nil {
		panic(err)
	}
}

const sqlVersion = `
create table if not exists version(
	revnum int, comment text, primary key(revnum)
);`

const sqlUsers = `
create table if not exists users (
	id integer,
	start_date datetime default null,
	username varchar(45) default null,
	PRIMARY KEY (id)
);`

const sqlRegistrations = `
create table if not exists registrations(
	test text
);`

var createStatements = []string{
	sqlVersion,
	sqlUsers,
	sqlRegistrations,
}

// CreateTables creates all tables if not exist.
func CreateTables() {
	for idx, sql := range createStatements {
		_, err := db.Exec(sql)
		if err != nil {
			panic(err)
		}
		fmt.Println("table created:", idx)
	}
}

// UpdateVersion adds version record to the version table.
func UpdateVersion(revisionNumber int, comment string) error {
	sql := "insert into version(revnum, comment) values (?, ?);"
	stmt, err := db.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(revisionNumber, comment)
	if err != nil {
		return err
	}
	return nil
}

// GetVersion returns max version from version table.
func GetVersion() (int, error) {
	sql := "select max(revnum) from version;"
	stmt, err := db.Prepare(sql)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return 0, err
	}

	if !rows.Next() {
		return 0, errors.New("no results")
	}

	var version int
	err = rows.Scan(&version)
	if err != nil {
		return 0, err
	}

	return version, nil
}

func generateRandomString(length int) string {
	b := make([]byte, length)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// AddRegistration adds new user registration.
func AddRegistration(name, email string) (id int, farsh string, err error) {
	id = 0
	farsh = ""
	err = errors.New("cannot create registration")
	return
}

// FindPendingRegistration gets pending user registration.
func FindPendingRegistration(id int, farsh string) error {
	return errors.New("registration not found")
}

// CompleteRegistration removes registration info and creates user record.
func CompleteRegistration(id int) (userID int, password string, err error) {
	userID = 0
	err = errors.New("cannot complete registration")
	return
}

// LoginUser finds the user and generates session token.
func LoginUser(user, password string) (userID int, sessionToken string, err error) {
	userID = 0
	sessionToken = ""
	err = errors.New("cannot login user. Wrong password")
	return
}

// AuthUser checks user token.
func AuthUser(userID int, sessionToken string) (name string, err error) {
	name = ""
	err = errors.New("invalid session token")
	return
}
