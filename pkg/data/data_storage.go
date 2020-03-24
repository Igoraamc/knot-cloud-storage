package data

import (
	"time"

	. "github.com/CESARBR/knot-cloud-storage/pkg/entities"
	"gopkg.in/mgo.v2/bson"
)

const collection = "things"

type DataStore struct{}

type IDataStore interface {
	Get()
	Create()
	Delete()
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

func (ds *DataStore) DeleteAll() error {
	_, err := db.C(collection).RemoveAll(nil)
	return err
}
func (ds *DataStore) DeleteByDevice(deviceId string) error {
	_, err := db.C(collection).RemoveAll(bson.M{"from": deviceId})
	return err
}
