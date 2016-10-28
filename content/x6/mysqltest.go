package main

import (
	"database/sql"                     // golang common interface for any database
	"fmt"                              // for Printf function
	_ "github.com/go-sql-driver/mysql" // mysql driver
	"io/ioutil"                        // for ReadFile
	"strings"                          // trim, split operations
	"time"                             // for datetime sql column
)

const password_file = "mysql~.txt"

var db *sql.DB

func readCredentials() (dbuser, dbpassword string) {
	bytes, err := ioutil.ReadFile(password_file)
	if err != nil {
		panic("cannot open " + password_file + " file")
	}

	str := string(bytes)
	str = strings.Trim(str, "\r\n ")
	arr := strings.Split(str, ":")
	if len(arr) != 2 {
		panic("invalid username and password in " + password_file + " file")
	}

	dbuser = arr[0]
	dbpassword = arr[1]
	return
}

func connectServer(dbuser, dbpassword string) {
	var err error

	connection_string := dbuser + ":" + dbpassword + "@tcp(wet.voilokov.com:3306)/mysql?parseTime=true"

	db, err = sql.Open("mysql", connection_string)
	if err != nil {
		panic(err)
	}
}

func createTestDatabase() {
	_, err := db.Exec("create database if not exists test;")
	if err != nil {
		panic(err)
	}

	fmt.Println("test databse created or exists")
}

func selectTestDatabase() {
	_, err := db.Exec("use test;")
	if err != nil {
		panic(err)
	}

	fmt.Println("test databse selected")
}

func createUsersTable() {
	sql := `
		create table if not exists users (
			user_id int(11) not null auto_increment,
			start_date datetime default null,
			username varchar(45) default null,
			PRIMARY KEY (user_id)
		 );`

	_, err := db.Exec(sql)
	if err != nil {
		panic(err)
	}

	fmt.Println("users table created")
}

func insertNewUsers() {
	sql := "insert into users(start_date, username) values(?, ?);"
	stmt, err := db.Prepare(sql)
	if err != nil {
		panic(err)
	}
	defer stmt.Close() // close on exit

	users := []string{"kay", "gerda"}

	for i := 0; i < 2; i++ {
		_, err := stmt.Exec(time.Now(), users[i])
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("2 new users created")
}

type user struct {
	user_id    int64
	start_date time.Time
	username   string
}

func getUsers() {
	sql := "select user_id, start_date, username from users order by start_date desc limit 10;"
	rows, err := db.Query(sql)
	if err != nil {
		panic(err)
	}

	var list []user

	for rows.Next() {
		u := user{}
		err = rows.Scan(&u.user_id, &u.start_date, &u.username)
		if err != nil {
			panic(err)
		}
		list = append(list, u)
	}

	fmt.Printf("print last %d users\n", len(list))

	fmt.Println("num user_id start_date                      username")
	fmt.Println("=== ======= =============================   ==============================")
	for idx, u := range list {
		fmt.Printf("%4d %6d %20s   %-30s\n", idx+1, u.user_id, u.start_date, u.username)
	}
}

func main() {
	dbuser, dbpassword := readCredentials()
	connectServer(dbuser, dbpassword)
	defer db.Close() // close connection on exit
	createTestDatabase()
	selectTestDatabase()
	createUsersTable()
	insertNewUsers()
	getUsers()
}
