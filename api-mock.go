//
// This contains a small mock that is used to simulate the API-Calls
//
package main

import (
	"time"
)

type ApiUser struct {
	Firstname string
	Lastname  string
	Email     string
}

var dummyUsers []ApiUser = []ApiUser{
	{
		Firstname: "Christal",
		Lastname:  "Goodwill",
		Email:     "cgoodwill0@geocities.jp",
	}, {
		Firstname: "Simonne",
		Lastname:  "Hadaway",
		Email:     "shadaway1@altervista.org",
	}, {
		Firstname: "Yolanthe",
		Lastname:  "Babidge",
		Email:     "ybabidge2@example.com",
	}, {
		Firstname: "Tuck",
		Lastname:  "Frankom",
		Email:     "tfrankom3@elpais.com",
	}, {
		Firstname: "Noach",
		Lastname:  "Besantie",
		Email:     "nbesantie4@jigsy.com",
	}, {
		Firstname: "Ivonne",
		Lastname:  "Morrel",
		Email:     "imorrel5@prweb.com",
	}, {
		Firstname: "Jamill",
		Lastname:  "Grabban",
		Email:     "jgrabban6@ameblo.jp",
	}, {
		Firstname: "Sheri",
		Lastname:  "O'Donohue",
		Email:     "sodonohue7@altervista.org",
	}, {
		Firstname: "Chevy",
		Lastname:  "Gontier",
		Email:     "cgontier8@gizmodo.com",
	}, {
		Firstname: "Jewel",
		Lastname:  "Aizikov",
		Email:     "jaizikov9@youtube.com",
	},
}

type API struct{}

func (a API) getList() ([]ApiUser, error) {
	time.Sleep(API_MOCK_TAKES)
	return dummyUsers, nil
}
