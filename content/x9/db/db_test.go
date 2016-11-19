package db

import (
	"testing"
	"time"
	"fmt"
)

func TestCreateDatabase(t *testing.T) {
	dbname := time.Now().Format("chat-20060102150405.db")
	ConnectServer(dbname)
	CreateTables()
}

var reg_id int           // registration id
var user_id int          // user id
var farsh string         // randomly generated registration farsh
var password string      // randomly generated password
var session_token string // randomly generated session token

var username = "user1"

func TestRegistration(t *testing.T) {
	var err error

	fmt.Println("User clicks Register button.")

	email := "user1@ttt.com"

	reg_id, farsh, err = AddRegistration(username, email)
	if err != nil {
		t.Fatal(err)
	}
	
	fmt.Println("Server sends email with id: ", reg_id, ", farsh: ", farsh)
	fmt.Println("User clicks confirmation link.")
	
	err = FindPendingRegistration(reg_id, farsh)
	if err != nil {
		t.Fatal(err)
	}
	
	user_id, password, err = CompleteRegistration(reg_id)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("user created: ", user_id)
}

func TestLogin(t *testing.T) {
	var err error

	fmt.Println("User clicks Login button.")

	user_id, session_token, err = LoginUser(username, password)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Server redirects to chat page.")
	
	name, err := AuthUser(user_id, session_token)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("User ", name, " successfully logged in!")
}
