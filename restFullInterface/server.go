package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"mux/mux-master"
	"net/http"
	"strconv"
)

type Information struct {
	ID           string    `json:"id"`
	Data1        string    `json:"data1"`
	Data2        string    `json:"data2"`
	DetailedData *MoreData `json:"detailedData"`
}

type MoreData struct {
	MoreData1 string `json:"moreData1"`
	MoreData2 string `json:"moreData2"`
}

var Info []Information

func getAllInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "json")
	json.NewEncoder(w).Encode(Info)
}

func getInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "json")
	params := mux.Vars(r)
	for _, someData := range Info {
		if someData.ID == params["id"] {
			json.NewEncoder(w).Encode(someData)
			return
		}
		json.NewEncoder(w).Encode(&Information{})
	}
}

func postInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "json")
	var newData Information
	json.NewDecoder(r.Body).Decode(&newData)
	newData.ID = strconv.Itoa(rand.Intn(1000000))
	Info = append(Info, newData)
	json.NewEncoder(w).Encode(newData)
}

func putInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "json")
	params := mux.Vars(r)
	for index, someData := range Info {
		if someData.ID == params["id"] {
			Info = append(Info[:index], Info[index+1:]...)
			var updatedData Information
			json.NewDecoder(r.Body).Decode(&updatedData)
			someData.ID = params["id"]
			Info = append(Info, updatedData)
			json.NewEncoder(w).Encode(updatedData)
			return
		}
	}
	json.NewEncoder(w).Encode(Info)
}

func deleteInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "json")
	params := mux.Vars(r)
	for index, someData := range Info {
		if someData.ID == params["id"] {
			Info = append(Info[:index], Info[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(Info)
}

func main() {
	router := mux.NewRouter()

	Info = append(Info, Information{ID: "1", Data1: "Intresting", Data2: "Golang", DetailedData: &MoreData{MoreData1: "1", MoreData2: "0"}})
	Info = append(Info, Information{ID: "2", Data1: "Fun", Data2: "Ruby", DetailedData: &MoreData{MoreData1: "0", MoreData2: "0"}})
	Info = append(Info, Information{ID: "3", Data1: "Terrific", Data2: "Azure", DetailedData: &MoreData{MoreData1: "1", MoreData2: "1"}})

	router.HandleFunc("/Info", getAllInfo).Methods("GET")
	router.HandleFunc("/Info/{id}", getInfo).Methods("GET")
	router.HandleFunc("/Info", postInfo).Methods("POST")
	router.HandleFunc("/Info/{id}", putInfo).Methods("PUT")
	router.HandleFunc("/Indo/{id}", deleteInfo).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}
