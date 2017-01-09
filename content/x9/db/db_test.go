package db

import (
	"fmt"
	"testing"
	"time"
)

func TestCreateDatabase(t *testing.T) {
	dbname := time.Now().Format("chat-20060102150405.db")
	ConnectServer(dbname)
	CreateTables()
	err := UpdateVersion(1, "just created")
	if err != nil {
		t.Fatal(err)
	}

	version, err = GetVersion()
	if err != nil {
		t.Fatal(err)
	}

	if version != 1 {
		t.Fatal("version is not 1")
	}
}

var regID int           // registration id
var userID int          // user id
var farsh string        // randomly generated registration farsh
var password string     // randomly generated password
var sessionToken string // randomly generated session token

var username = "user1"

func TestRegistration(t *testing.T) {
	var err error

	fmt.Println("User clicks Register button.")

	email := "user1@ttt.com"

	regID, farsh, err = AddRegistration(username, email)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Server sends email with id: ", regID, ", farsh: ", farsh)
	fmt.Println("User clicks confirmation link.")

	err = FindPendingRegistration(regID, farsh)
	if err != nil {
		t.Fatal(err)
	}

	userID, password, err = CompleteRegistration(regID)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("user created: ", userID)
}

func TestLogin(t *testing.T) {
	var err error

	fmt.Println("User clicks Login button.")

	userID, sessionToken, err = LoginUser(username, password)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Server redirects to chat page.")

	name, err := AuthUser(userID, sessionToken)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("User ", name, " successfully logged in!")
}
