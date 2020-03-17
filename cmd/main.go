package main

import (
	"fmt"

	. "github.com/Igoraamc/knot-cloud-storage/pkg/config"
	"github.com/Igoraamc/knot-cloud-storage/pkg/controllers"
	. "github.com/Igoraamc/knot-cloud-storage/pkg/interactor"
	"github.com/Igoraamc/knot-cloud-storage/pkg/server"
)

var dao = DataDAO{}
var config = Config{}

func init() {
	config.Read()

	dao.Server = config.Server
	dao.Database = config.Database
	fmt.Println(dao.Server)
	dao.Connect()
}

func main() {
	dataController := controllers.NewDataController()

	http := server.NewServer(config.Port, dataController)
	http.Start()
}
