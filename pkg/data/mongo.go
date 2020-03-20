package data

import (
	"log"

	"gopkg.in/mgo.v2"
)

type DataDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

func (m *DataDAO) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}
