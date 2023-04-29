package main

import (
	"encoding/json"
	"math"
	"net/http"
	"time"
)

var ParkingLots map[string][]vehicles

type vehicles struct {
	entryTime time.Time
	plate     string
}

type Entry struct {
	Plate      string `json:"plate"`
	ParkingLot string `json:"parking_lot"`
}

type exitReq struct {
	ParkingLot string `json:"parking_lot"`
	TicketId   int    `json:"ticket_id"`
}

func createEntry(w http.ResponseWriter, r *http.Request) {
	var entry Entry
	err := json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newVehicle := vehicles{
		entryTime: time.Now(),
		plate:     entry.Plate,
	}

	ParkingLots[entry.ParkingLot] = []vehicles{newVehicle}

	type Response struct {
		ParkingLot string `json:"parking_lot"`
		TicketId   int    `json:"ticket_id"`
	}
	response := Response{
		ParkingLot: entry.ParkingLot,
		TicketId:   len(ParkingLots[entry.ParkingLot]) - 1,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func exit(w http.ResponseWriter, r *http.Request) {
	var exitReq exitReq
	err := json.NewDecoder(r.Body).Decode(&exitReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vehicle := ParkingLots[exitReq.ParkingLot][exitReq.TicketId]

	totalTime := (time.Now().Sub(vehicle.entryTime)).Minutes()
	roundedMinutes := int(math.Ceil(float64(totalTime)/15)) * 15
	totalCharge := float64(roundedMinutes) / 15 * 2.5

	type Response struct {
		Charge float64 `json:"charge"`
	}
	response := Response{Charge: totalCharge}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	ParkingLots = make(map[string][]vehicles)
	http.HandleFunc("/entry", createEntry)
	http.HandleFunc("/exit", exit)
	http.ListenAndServe(":8080", nil)
}
