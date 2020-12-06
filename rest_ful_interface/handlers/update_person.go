package handlers

import (
	"encoding/json"
	"go-chi/chi-master"
	"io"
	"net/http"
	"os"
)

type UpdatePersonInfoRepo interface {
	GetThemAll() error
}

func UpdatePersonInfo(mux chi.Router, repo UpdatePersonInfoRepo) {
	mux.Put("/Info/{Id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "json")
		personID := chi.URLParam(r, "Id")

		data, err := os.Open("DB.txt")
		if er := repo.GetThemAll(); er != nil {
			http.Error(w, er.Error(), http.StatusUnprocessableEntity)
			return
		}
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
			if workData.Id == personID {
				json.NewDecoder(r.Body).Decode(&newData)
				newData.Id = personID
				allData = append(allData, *newData)
			} else {
				allData = append(allData, *workData)
			}
		}
		jsonResponse, _ := json.Marshal(allData)
		w.WriteHeader(http.StatusAccepted)
		io.WriteString(w, string(jsonResponse))
	})
}
