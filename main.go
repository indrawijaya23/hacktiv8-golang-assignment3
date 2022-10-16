package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type fileStatus struct {
	Status status `json:"status"`
}

type status struct {
	Water       int    `json:"water"`
	Wind        int    `json:"wind"`
	WaterStatus string `json:"-"`
	WindStatus  string `json:"-"`
}

var s status
var fs fileStatus

const PORT = ":8080"

func main() {
	refreshStatus()
	go refreshStatusForever()

	http.HandleFunc("/", showStatus)
	http.ListenAndServe(PORT, nil)
}

func refreshStatus() {
	s.Water = rand.Intn(99) + 1
	s.Wind = rand.Intn(99) + 1
	fs.Status = s

	file, _ := json.MarshalIndent(fs, " ", " ")
	_ = ioutil.WriteFile("data.json", file, 0644)
	time.Sleep(time.Second * 15)
}

func refreshStatusForever() {
	for {
		refreshStatus()
	}
}

func showStatus(w http.ResponseWriter, r *http.Request) {
	file, _ := os.Open("data.json")
	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)

	var fsRender fileStatus

	json.Unmarshal(byteValue, &fsRender)

	if fsRender.Status.Water < 5 {
		fsRender.Status.WaterStatus = "aman"
	} else if fsRender.Status.Water <= 8 {
		fsRender.Status.WaterStatus = "siaga"
	} else {
		fsRender.Status.WaterStatus = "bahaya"
	}

	if fsRender.Status.Wind < 6 {
		fsRender.Status.WindStatus = "aman"
	} else if fsRender.Status.Wind <= 15 {
		fsRender.Status.WindStatus = "siaga"
	} else {
		fsRender.Status.WindStatus = "bahaya"
	}

	tpl, _ := template.ParseFiles("template.html")
	tpl.Execute(w, fsRender)
}
