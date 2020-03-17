package interactor

import (
	"log"
	"strconv"
	"time"

	. "github.com/Igoraamc/knot-cloud-storage/pkg/entities"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DataDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	COLLECTION = "things"
)

func (m *DataDAO) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

func (m *DataDAO) GetAll(order int, skip int, take int, startDate time.Time, finishDate time.Time) ([]Data, error) {
	var data []Data

	selectOrder := "timestamp"
	if order == -1 {
		selectOrder = "-timestamp"
	}

	err := db.C(COLLECTION).Find(bson.M{
		"timestamp": bson.M{
			"$gt": startDate,
			"$lt": finishDate,
		},
	}).Select(bson.M{
		"timestamp": 1,
		"payload":   1,
		"from":      1,
	}).Skip(skip).Sort(selectOrder).Limit(take).All(&data)

	if data == nil {
		data = []Data{}
	}

	return data, err
}

func (m *DataDAO) GetByID(id string, order int, skip int, take int, startDate time.Time, finishDate time.Time) ([]Data, error) {
	var data []Data

	selectOrder := "timestamp"
	if order == -1 {
		selectOrder = "-timestamp"
	}

	s, _ := strconv.ParseInt(id, 10, 64)

	err := db.C(COLLECTION).Find(bson.M{
		"timestamp": bson.M{
			"$gt": startDate,
			"$lt": finishDate,
		},
	}).Select(bson.M{
		"timestamp": 1,
		"payload":   1,
		"from":      1,
	}).Skip(skip).Sort(selectOrder).Limit(take).All(&data)

	if data == nil {
		data = []Data{}
	}

	data = FilterBySensorId(data, int(s))
	return data, err
}

func (m *DataDAO) Create(data Data) error {
	err := db.C(COLLECTION).Insert(&data)
	return err
}

func (m *DataDAO) Delete(id string) error {
	err := db.C(COLLECTION).RemoveId(bson.ObjectIdHex(id))
	return err
}

func (m *DataDAO) Update(id string, data Data) error {
	err := db.C(COLLECTION).UpdateId(bson.ObjectIdHex(id), &data)
	return err
}

func FilterBySensorId(array []Data, sensorId int) []Data {
	newArray := make([]Data, 0)
	for _, v := range array {
		if v.Payload.SensorId == sensorId {
			newArray = append(newArray, v)
		}
	}
	return newArray
}
