//
// This contains a small mock that is used to simulate the Database-Queries
//
package main

import (
	"errors"
	"math/rand"
	"time"
)

type DBUser struct {
	Firstname string
	Lastname  string
	Email     string
}

type DB struct{}

var UserNotFoundError error = errors.New("User not found")

func (db DB) getUserIdByEmail(email string) (int, error) {
	time.Sleep(DB_MOCK_TAKES_FOR_GET_USERID)
	if rand.Intn(3) == 0 {
		return -1, UserNotFoundError
	}
	return rand.Int(), nil
}

func (db DB) updateUser(id int, firstname string, lastname string, email string) error {
	time.Sleep(DB_MOCK_TAKES_FOR_UPDATE_USER)
	return nil
}

func (db DB) createUser(firstname string, lastname string, email string) error {
	time.Sleep(DB_MOCK_TAKES_FOR_CREATE_USER)
	return nil
}
