package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"go-chi/chi-master"
	"net/http"
	"net/http/httptest"
	"testing"
)

type PostPersonInfoRepoMock struct {
	err error
}

func (m *PostPersonInfoRepoMock) GetThemAll() error {
	return m.err
}

func TestPostPersonInfo(t *testing.T) {
	t.Run("sends bad gateway on file opening error", func(t *testing.T) {
		mux := chi.NewMux()
		PostPersonInfo(mux, &PostPersonInfoRepoMock{err: errors.New("No file opening:(((")})

		req := httptest.NewRequest(http.MethodPost, "/Info", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		res := rec.Result()
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.FailNow()
		}
	})

	t.Run("sends accepted on successful finish", func(t *testing.T) {
		mux := chi.NewMux()
		PostPersonInfo(mux, &PostPersonInfoRepoMock{})
		data := PersonInfo{
			Id:       "4",
			Name:     "Kira",
			LastName: "Yoshikage",
			Adress: []Adress{
				Adress{
					City:   "Morio",
					Street: "KillerQueen",
				},
			},
		}
		file, _ := json.MarshalIndent(data, "", " ")
		req := httptest.NewRequest(http.MethodPost, "/Info", bytes.NewBuffer(file))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		res := rec.Result()
		if res.StatusCode != http.StatusAccepted {
			t.FailNow()
		}
	})

}
