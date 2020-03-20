package interactor

import (
	"strconv"
	"time"

	"github.com/CESARBR/knot-cloud-storage/pkg/data"
	. "github.com/CESARBR/knot-cloud-storage/pkg/entities"
)

type DataInteractor struct{}

var dataStore = data.NewDataStore()

func (d *DataInteractor) GetAll(order int, skip int, take int, startDate time.Time, finishDate time.Time) ([]Data, error) {

	selectOrder := "timestamp"
	if order == -1 {
		selectOrder = "-timestamp"
	}

	data, err := dataStore.GetAll(selectOrder, skip, take, startDate, finishDate)

	return data, err
}

func (d *DataInteractor) GetByID(id string, order int, skip int, take int, startDate time.Time, finishDate time.Time) ([]Data, error) {

	selectOrder := "timestamp"
	if order == -1 {
		selectOrder = "-timestamp"
	}

	s, _ := strconv.ParseInt(id, 10, 64)

	data, err := dataStore.GetByID(selectOrder, skip, take, startDate, finishDate)

	data = FilterBySensorId(data, int(s))
	return data, err
}

func (d *DataInteractor) Create(data Data) error {
	err := dataStore.Create(data)
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
