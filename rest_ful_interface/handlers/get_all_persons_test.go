package handlers

import (
	"errors"
	"go-chi/chi-master"
	"net/http"
	"net/http/httptest"
	"testing"
)

type GetAllPersonsInfoRepoMock struct {
	err error
}

func (m *GetAllPersonsInfoRepoMock) GetThemAll() error {
	return m.err
}

func TestGetAllPersonsInfo(t *testing.T) {
	t.Run("sends bad gateway on file opening error", func(t *testing.T) {
		mux := chi.NewMux()
		GetAllPersonsInfo(mux, &GetAllPersonsInfoRepoMock{err: errors.New("No file opening:(((")})

		req := httptest.NewRequest(http.MethodGet, "/Info", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		res := rec.Result()
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.FailNow()
		}
	})

	t.Run("sends accepted on successful finish", func(t *testing.T) {
		mux := chi.NewMux()
		GetAllPersonsInfo(mux, &GetAllPersonsInfoRepoMock{})

		req := httptest.NewRequest(http.MethodGet, "/Info", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		res := rec.Result()
		if res.StatusCode != http.StatusAccepted {
			t.FailNow()
		}
	})

}
