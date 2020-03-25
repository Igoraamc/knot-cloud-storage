package data

import (
	"gopkg.in/mgo.v2"
)

var db *mgo.Database

type Storage struct {
	Server   string
	Database string
}

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
