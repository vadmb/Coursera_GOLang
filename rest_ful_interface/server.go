package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"mux/mux-master"
	"net/http"
	"os"
	"strconv"
)

type PersonInfo struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"lastName"`
	Adress   []struct {
		City   string `json:"city"`
		Street string `json:"street"`
	} `json:"adress"`
}

func getAllPersonsInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "json")

	data, err := os.Open("DB.txt")
	if err != nil {
		panic(err)
	}

	defer data.Close()

	workData := &PersonInfo{}
	var allData []PersonInfo
	decoder := json.NewDecoder(data)
	for {
		if err := decoder.Decode(&workData); err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		allData = append(allData, *workData)
	}
	jsonResponse, _ := json.Marshal(allData)
	io.WriteString(w, string(jsonResponse))
}

func getPersonInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "json")

	data, err := os.Open("DB.txt")
	if err != nil {
		panic(err)
	}

	defer data.Close()

	params := mux.Vars(r)
	workData := &PersonInfo{}
	decoder := json.NewDecoder(data)
	for {
		if err := decoder.Decode(&workData); err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		if workData.Id == params["id"] {
			jsonResponse, _ := json.Marshal(workData)
			io.WriteString(w, string(jsonResponse))

		}
	}
}

func postPersonInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "json")
	newData := &PersonInfo{}
	json.NewDecoder(r.Body).Decode(&newData)
	newData.Id = strconv.Itoa(rand.Intn(1000000))

	data, err := os.Open("DB.txt")
	if err != nil {
		panic(err)
	}

	workData := &PersonInfo{}
	var allData []PersonInfo
	decoder := json.NewDecoder(data)
	for {
		if err := decoder.Decode(&workData); err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		allData = append(allData, *workData)
	}

	allData = append(allData, *newData)

	jsonResponse, _ := json.Marshal(allData)
	ioutil.WriteFile("DB.txt", jsonResponse, 0644)
}

func updatePersonInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "json")
	params := mux.Vars(r)
	data, err := os.Open("DB.txt")
	if err != nil {
		panic(err)
	}
	defer data.Close()

	newData := &PersonInfo{}

	workData := &PersonInfo{}
	var allData []PersonInfo
	decoder := json.NewDecoder(data)
	for {
		if err := decoder.Decode(&workData); err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		if workData.Id == params["id"] {
			json.NewDecoder(r.Body).Decode(&newData)
			newData.Id = params["id"]
			allData = append(allData, *newData)
		} else {
			allData = append(allData, *workData)
		}
	}
	jsonResponse, _ := json.Marshal(allData)
	io.WriteString(w, string(jsonResponse))
}

func deletePersonInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "json")
	params := mux.Vars(r)
	data, err := os.Open("DB.txt")
	if err != nil {
		panic(err)
	}
	defer data.Close()

	workData := &PersonInfo{}
	var allData []PersonInfo
	decoder := json.NewDecoder(data)
	for {
		if err := decoder.Decode(&workData); err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		if workData.Id == params["id"] {

		} else {
			allData = append(allData, *workData)
		}
	}
	jsonResponse, _ := json.Marshal(allData)
	io.WriteString(w, string(jsonResponse))
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/allInfo", getAllPersonsInfo)
	router.HandleFunc("/info/{id}", getPersonInfo)
	router.HandleFunc("/postInfo", postPersonInfo).Methods("POST")
	router.HandleFunc("/updateInfo/{id}", updatePersonInfo).Methods("PUT")
	router.HandleFunc("/deleteInfo/{id}", deletePersonInfo).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}
