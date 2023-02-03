package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
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

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "AgentQS homepage endpoint hit")
}

func receiveAckFromAgentQ(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: receiveAckFromAgentQ")

	responseData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(responseData))
}

func handleRequests() {
	http.HandleFunc("/", homePage)

	// ack route and map it to our receiveAckFromAgentQ function like so
	http.HandleFunc("/ackFromAgentQ", receiveAckFromAgentQ)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func readSchedulesAndPushSpls() {
	fmt.Println("Start of readSchedulesAndPushSpls")

	// reads the schedule from schedule log table
	Schedules = []Schedule{
		{Id: 12500, Name: "Andhra GAP 112206603 AP", OrderId: "QCS/011097/22-23", GovtRoNo: "ENE68-OPOM0COMM(NEWC)/1/2022-HRA"},
	}

	// creates the spl
	Spls = []Spl{
		{Name: "SPL1", ScheduleTag: ScheduleTag{Id: 12500, Name: "Andhra GAP 112206603 AP"}},
	}

	// push spl to agentq via receiveSpls api
	jsonValue, _ := json.Marshal(Spls)
	response, err := http.Post("http://localhost:8082/receiveSpls", "application/json", bytes.NewBuffer(jsonValue))

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
	// initialise agentqs apis
	go handleRequests()

	// end-less process
	for {
		rand.Seed(time.Now().UnixNano())
		time.Sleep(time.Second * time.Duration(rand.Intn(10)))

		// function to read shedule from db, create and send spl after a random wait to agentq
		readSchedulesAndPushSpls()
	}
}
