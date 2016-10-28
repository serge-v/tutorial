package main

import (
	"database/sql"                     // golang common interface for any database
	"fmt"                              // for Printf function
	_ "github.com/go-sql-driver/mysql" // mysql driver
	"io/ioutil"                        // for ReadFile
	"strings"                          // trim, split operations
	"time"                             // for datetime sql column
	"os"
)

const password_file = "mysql~.txt"

var db *sql.DB

// readCredentials reads mysql~.txt file. File should contain the line in format
// username:password
// function returns a pair dbuser, dbpassword
func readCredentials() (dbuser, dbpassword string) {
	bytes, err := ioutil.ReadFile(password_file) // read all bytes from file
	if err != nil {
		panic("cannot open " + password_file + " file") // print message and exit if error
	}

	str := string(bytes)             // convert bytes to the string
	str = strings.Trim(str, "\r\n ") // trim space or line endings from string
	arr := strings.Split(str, ":")   // split string into slice
	if len(arr) != 2 { // if not 2 elemens in slice then fail
		panic("invalid username and password in " + password_file + " file")
	}

	dbuser = arr[0]
	dbpassword = arr[1]
	return
}

// connectServer makes a connection to Mysql server. It uses global db variable.
func connectServer(dbuser, dbpassword string) {
	var err error

	// create connection string to connect to default mysql database
	// parseTime=true enables converting to golang time.Time time
	connection_string := dbuser + ":" + dbpassword + "@tcp(wet.voilokov.com:3306)/mysql?parseTime=true"

	db, err = sql.Open("mysql", connection_string)
	if err != nil {
		panic(err)
	}
	
	// global db variable now holds the connection
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
	stmt, err := db.Prepare(sql)    // prepare statement
	if err != nil {
		panic(err)
	}
	defer stmt.Close() // close on exit

	users := []string{"kay", "gerda"} // test user names

	for i := 0; i < 2; i++ {
		_, err := stmt.Exec(time.Now(), users[i]) // execute insert with parameters 
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

func getUsers() (list []user) {
	sql := "select user_id, start_date, username from users order by start_date desc limit 10;"
	rows, err := db.Query(sql)
	if err != nil {
		panic(err)
	}

	for rows.Next() {            // while next row available
		u := user{}              // create empty user structure
		err = rows.Scan(&u.user_id, &u.start_date, &u.username) // read into structure members
		if err != nil {
			panic(err)
		}
		list = append(list, u)
	}

	return list
}

func printUsers(list []user) {
	fmt.Printf("print last %d users\n", len(list))
	fmt.Println("num user_id start_date                      username")
	fmt.Println("=== ======= =============================   ==============================")
	for idx, u := range list {
		fmt.Printf("%4d %6d %20s   %-30s\n", idx+1, u.user_id, u.start_date, u.username)
	}
}

func deleteExtraUsers(list []user) {
	if len(list) <= 10 {
		return
	}

	sql := "delete from users where user_id = ?;"
	stmt, err := db.Prepare(sql)    // prepare statement
	if err != nil {
		panic(err)
	}
	defer stmt.Close() // close on exit

	for i := 10; i < len(list); i++ {
		_, err := stmt.Exec(list[i].user_id) // execute delete with parameter
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("%d users deleted\n", len(list)-10)
}

func saveUsers(list []user) {
	f, err := os.Create("users.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close() // close file on function exit

	for idx, u := range list {
		fmt.Fprintf(f, "%d: %+v\n", idx+1, u)
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
	list := getUsers()
	printUsers(list)
	deleteExtraUsers(list)
}
