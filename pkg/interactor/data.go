package interactor

import (
	"strconv"
	"time"

	"github.com/CESARBR/knot-cloud-storage/pkg/data"
	. "github.com/CESARBR/knot-cloud-storage/pkg/entities"
)

type DataInteractor struct {
	DataStore data.DataStore
}

func NewDataInteractor(dataStore data.DataStore) DataInteractor {
	return DataInteractor{dataStore}
}

func (d *DataInteractor) GetAll(order int, skip int, take int, startDate time.Time, finishDate time.Time) ([]Data, error) {

	selectOrder := "timestamp"
	if order == -1 {
		selectOrder = "-timestamp"
	}

	data, err := d.DataStore.Get(selectOrder, skip, take, startDate, finishDate)

	return data, err
}

func (d *DataInteractor) GetByID(id string, order int, skip int, take int, startDate time.Time, finishDate time.Time) ([]Data, error) {

	selectOrder := "timestamp"
	if order == -1 {
		selectOrder = "-timestamp"
	}

	s, _ := strconv.ParseInt(id, 10, 64)

	data, err := d.DataStore.Get(selectOrder, skip, take, startDate, finishDate)

	data = FilterDataBySensorID(data, int(s))
	return data, err
}

func (d *DataInteractor) Save(data Data) error {
	err := d.DataStore.Save(data)
	return err
}

func FilterDataBySensorID(data []Data, sensorId int) []Data {
	filteredData := make([]Data, 0)
	for _, v := range data {
		if v.Payload.SensorId == sensorId {
			filteredData = append(filteredData, v)
		}
	}
	return filteredData
}
