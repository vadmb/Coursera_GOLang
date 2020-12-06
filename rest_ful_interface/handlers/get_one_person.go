package handlers

import (
	"encoding/json"
	"go-chi/chi-master"
	"io"
	"net/http"
	"os"
)

type GetPersonInfoRepo interface {
	GetThemAll() error
}

// func GetPersonInfo(w http.ResponseWriter, r *http.Request) {
func GetPersonInfo(mux chi.Router, repo GetPersonInfoRepo) {
	mux.Get("/Info/{Id}", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "json")

		data, err := os.Open("DB.txt")
		if er := repo.GetThemAll(); er != nil {
			http.Error(w, er.Error(), http.StatusUnprocessableEntity)
			return
		}
		if err != nil {
			panic(err)
		}

		defer data.Close()

		personID := chi.URLParam(r, "Id")
		workData := &PersonInfo{}
		decoder := json.NewDecoder(data)
		for {
			if err := decoder.Decode(&workData); err == io.EOF {
				break
			} else if err != nil {
				panic(err)
			}

			if workData.Id == personID {
				jsonResponse, _ := json.Marshal(workData)
				w.WriteHeader(http.StatusAccepted)
				io.WriteString(w, string(jsonResponse))

			}
		}
	})
}
