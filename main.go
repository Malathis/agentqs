package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Schedule struct {
	Id       int    `json:"id"`
	Name     string `json:"schedule"`
	OrderId  string `json:"orderId"`
	GovtRoNo string `json:"govtRoNo"`
}

type ScheduleTag struct {
	Id   int    `json:"id"`
	Name string `json:"schedule"`
}

type Spl struct {
	Name        string      `json:"spl"`
	ScheduleTag ScheduleTag `json:"scheduleTag"`
}

// let's declare a global Schedules array
// that we can populate to simulate a database
var Schedules []Schedule

// let's declare a global Spls array
// that we can populate by creating spls from scheudles
var Spls []Spl

func readSchedulesAndPushSpls() {
	fmt.Println("Start of readSchedulesAndPushSpls")

	// reads the schedules from schedule log table
	Schedules = []Schedule{
		Schedule{Name: "Andhra GAP 112206603 AP", OrderId: "QCS/011097/22-23", GovtRoNo: "ENE68-OPOM0COMM(NEWC)/1/2022-HRA"},
		Schedule{Name: "Andhra GAP 112206603 AP AS", OrderId: "QCS/011254/22-23", GovtRoNo: "ENE68-OPOM0COMM(NEWC)/1/2022-HRA"},
	}

	// creates the spls
	Spls = []Spl{
		Spl{Name: "SPL1", ScheduleTag: ScheduleTag{Id: 12500, Name: "Andhra GAP 112206603 AP"}},
		Spl{Name: "SPL2", ScheduleTag: ScheduleTag{Id: 12501, Name: "Andhra GAP 112206603 AP AS"}},
	}

	// push spls to agentq via getSpl api
	jsonValue, _ := json.Marshal(Spls)
	response, err := http.Post("http://localhost:8082/getSpls", "application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(responseData))
}

func main() {
	// function to read shedules from db, create spls from schedules and push spls to agentq
	readSchedulesAndPushSpls()
}
