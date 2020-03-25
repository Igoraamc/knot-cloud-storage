package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	. "github.com/CESARBR/knot-cloud-storage/pkg/entities"
	"github.com/CESARBR/knot-cloud-storage/pkg/interactor"
	"github.com/gorilla/mux"
)

var dataInteractor = interactor.NewDataInteractor()

type DataController struct{}

type DataInterface interface {
	GetAll()
	GetByID()
	Create()
}

type errorMessage struct {
	error   bool
	message string
}

func NewDataController() *DataController {
	return &DataController{}
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"message": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (d *DataController) GetAll(w http.ResponseWriter, r *http.Request) {

	order, skip, take, startDate, finishDate, errUrl := getUrlQueryParams(r)
	if errUrl.error != false {
		respondWithError(w, http.StatusUnprocessableEntity, errUrl.message)
		return
	}

	things, err := dataInteractor.GetAll(order, skip, take, startDate, finishDate)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, things)
}

func (d *DataController) GetByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	order, skip, take, startDate, finishDate, errUrl := getUrlQueryParams(r)
	if errUrl.error != false {
		respondWithError(w, http.StatusUnprocessableEntity, errUrl.message)
		return
	}
	thing, err := dataInteractor.GetByID(params["id"], order, skip, take, startDate, finishDate)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Thing ID")
	}
	respondWithJson(w, http.StatusOK, thing)
}

func (d *DataController) Save(w http.ResponseWriter, r *http.Request) {
	var thing Data
	if err := json.NewDecoder(r.Body).Decode(&thing); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	// thing.ID = bson.NewObjectId()
	thing.Timestamp = time.Now()
	if err := dataInteractor.Save(thing); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, thing)
	defer r.Body.Close()

}

func getUrlQueryParams(r *http.Request) (order int, skip int, take int, startDate time.Time, finishDate time.Time, errorStatus errorMessage) {
	order = 1
	skip = 0
	take = 10
	startDate, _ = time.Parse("2006-1-02", "2006-1-02")
	finishDate = time.Now()
	var err error = nil

	if len(r.URL.Query()["order"]) > 0 {
		order, err = strconv.Atoi(r.URL.Query()["order"][0])
		if err != nil {
			errorStatus.error = true
			errorStatus.message = "order must be in the following format: 1 or -1"
		}
	}

	if len(r.URL.Query()["skip"]) > 0 {
		skip, err = strconv.Atoi(r.URL.Query()["skip"][0])
		if err != nil {
			errorStatus.error = true
			errorStatus.message = "skip must be an integer"
		}
	}

	if len(r.URL.Query()["take"]) > 0 {
		take, err = strconv.Atoi(r.URL.Query()["take"][0])
		if err != nil {
			errorStatus.error = true
			errorStatus.message = "take must be an integer (maximum 100)"
		}

		if take > 100 {
			take = 100
		}
	}

	if len(r.URL.Query()["startDate"]) > 0 {
		fmt.Println(r.URL.Query()["startDate"][0])

		startDate, err = time.Parse("2006-01-02 15:04:05", r.URL.Query()["startDate"][0])
		if err != nil {
			errorStatus.error = true
			errorStatus.message = "date must be in the following format: YYYY-MM-DD HH:MM:SS"
		}
	}

	if len(r.URL.Query()["finishDate"]) > 0 {
		finishDate, err = time.Parse("2006-01-02 15:04:05", r.URL.Query()["finishDate"][0])
		if err != nil {
			errorStatus.error = true
			errorStatus.message = "date must be in the following format: YYYY-MM-DD HH:MM:SS"
		}
	}

	return order, skip, take, startDate, finishDate, errorStatus
}
