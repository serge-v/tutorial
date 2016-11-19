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

func CreateTables() {
	users := `create table if not exists users (
	id integer,
        start_date datetime default null,
        username varchar(45) default null,
        PRIMARY KEY (id)
     );`

	registrations := `create table if not exists registrations(test text);`

	sqls := []string{users, registrations}

	for idx, sql := range sqls {
		_, err := db.Exec(sql)
		if err != nil {
			panic(err)
		}
		fmt.Println("table created:", idx)
	}
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
