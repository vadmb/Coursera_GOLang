package handlers

import (
	"encoding/json"
	"go-chi/chi-master"
	"io"
	"net/http"
	"os"
)

type Adress struct {
	City   string `json:"city"`
	Street string `json:"street"`
}

type PersonInfo struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"lastName"`
	Adress   []Adress
	// 	City   string `json:"city"`
	// 	Street string `json:"street"`
	// } `json:"adress"`
}

type GetAllPersonsInfoRepo interface {
	GetThemAll() error
}

func GetAllPersonsInfo(mux chi.Router, repo GetAllPersonsInfoRepo) {
	mux.Get("/Info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "json")

		data, err := os.Open("DB.txt")
		er := repo.GetThemAll()
		if er != nil {
			http.Error(w, er.Error(), http.StatusUnprocessableEntity)
			return
		}
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
		w.WriteHeader(http.StatusAccepted)
		io.WriteString(w, string(jsonResponse))
	})
}
