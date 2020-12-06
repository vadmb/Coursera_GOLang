package handlers

import (
	"encoding/json"
	"go-chi/chi-master"
	"io"
	"io/ioutil"

	//	"math/rand"
	"net/http"
	"os"
	//	"strconv"
)

type PostPersonInfoRepo interface {
	GetThemAll() error
}

func PostPersonInfo(mux chi.Router, repo PostPersonInfoRepo) {
	mux.Post("/Info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "json")
		newData := &PersonInfo{}
		json.NewDecoder(r.Body).Decode(&newData)
		//newData.Id = strconv.Itoa(rand.Intn(1000000))

		data, err := os.Open("DB.txt")
		if er := repo.GetThemAll(); er != nil {
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

		allData = append(allData, *newData)

		jsonResponse, _ := json.Marshal(allData)
		ioutil.WriteFile("DB.txt", jsonResponse, 0644)
		w.WriteHeader(http.StatusAccepted)
	})
}
