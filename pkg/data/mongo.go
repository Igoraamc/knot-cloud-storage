package data

import (
	"gopkg.in/mgo.v2"
)

var db *mgo.Database

type MongoDB struct {
	Server   string
	Database string
}

func NewMongoDB(server string, database string) *MongoDB {
	return &MongoDB{server, database}
}

func (s *MongoDB) Connect() (*mgo.Database, error) {
	session, err := mgo.Dial(s.Server)
	if err != nil {
		return nil, err
	}
	database := session.DB(s.Database)

	return database, nil
}
