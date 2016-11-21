package db

import (
	"database/sql"                  // generic go database package
	_ "github.com/mattn/go-sqlite3" // driver for sqlite
	"errors"                        // for creating return errors
	"fmt"                           // for formatted printing
	"crypto/rand"                   // for random generation
	"encoding/base64"               // for random generation
)

var db *sql.DB

func ConnectServer(dbname string) {
	var err error
	db, err = sql.Open("sqlite3", dbname) // different open parameters
	if err != nil {
		panic(err)
	}
}

const sql_version = `
create table if not exists version(
	revnum int, comment text, primary key(revnum)
);`

const sql_users = `
create table if not exists users (
	id integer,
	start_date datetime default null,
	username varchar(45) default null,
	PRIMARY KEY (id)
);`

const sql_registrations = `
create table if not exists registrations(
	test text
);`

var create_statements = []string{
	sql_version,
	sql_users,
	sql_registrations,
}

func CreateTables() {
	for idx, sql := range create_statements {
		_, err := db.Exec(sql)
		if err != nil {
			panic(err)
		}
		fmt.Println("table created:", idx)
	}
}

func UpdateVersion(revision_number int, comment string) error {
	sql := "insert into version(revnum, comment) values (?, ?);"
	stmt, err := db.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(revision_number, comment)
	if err != nil {
		return err
	}
	return nil
}

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

func AddRegistration(name, email string) (id int, farsh string, err error) {
	id = 0
	farsh = ""
	err = errors.New("cannot create registration")
	return
}

func FindPendingRegistration(id int, farsh string) error {
	return errors.New("registration not found")
}

func CompleteRegistration(id int) (user_id int, password string, err error) {
	user_id = 0
	err = errors.New("cannot complete registration")
	return
}

func LoginUser(user, password string) (user_id int, session_token string, err error) {
	user_id = 0
	session_token = ""
	err = errors.New("cannot login user. Wrong password.")
	return
}

func AuthUser(user_id int, session_token string) (name string, err error) {
	name = ""
	err = errors.New("invalid session token")
	return
}
