package data

import (
	"gopkg.in/mgo.v2"
)

type Storage struct {
	Server   string
	Database string
}

var db *mgo.Database

func NewStorage(server string, database string) *Storage {
	return &Storage{server, database}
}

func (s *Storage) Connect() error {
	session, err := mgo.Dial(s.Server)
	if err != nil {
		return err
	}
	db = session.DB(s.Database)

	return nil
}
