package data

import (
	"time"

	. "github.com/CESARBR/knot-cloud-storage/pkg/entities"
	"gopkg.in/mgo.v2/bson"
)

const collection = "things"

type DataStore struct{}

type IDataStore interface {
	Get(order string, skip int, take int, startDate time.Time, finishDate time.Time) ([]Data, error)
	Create(data Data) error
	Delete(param bson.M) error
}

func NewDataStore() DataStore {
	return DataStore{}
}

func (ds *DataStore) Get(order string, skip int, take int, startDate time.Time, finishDate time.Time) ([]Data, error) {
	var data []Data

	err := db.C(collection).Find(bson.M{
		"timestamp": bson.M{
			"$gt": startDate,
			"$lt": finishDate,
		},
	}).Select(bson.M{
		"timestamp": 1,
		"payload":   1,
		"from":      1,
	}).Skip(skip).Sort(order).Limit(take).All(&data)

	if data == nil {
		data = []Data{}
	}

	return data, err
}

func (ds *DataStore) Save(data Data) error {
	err := db.C(collection).Insert(&data)
	return err
}

func (ds *DataStore) Delete(param bson.M) error {
	_, err := db.C(collection).RemoveAll(param)
	return err
}
