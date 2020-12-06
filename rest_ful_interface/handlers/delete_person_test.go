package handlers

import (
	"errors"
	"go-chi/chi-master"
	"net/http"
	"net/http/httptest"
	"testing"
)

type DeletePersonInfoRepoMock struct {
	err error
}

func (m *DeletePersonInfoRepoMock) GetThemAll() error {
	return m.err
}

func TestDeletePersonInfo(t *testing.T) {
	t.Run("sends bad gateway on file opening error", func(t *testing.T) {
		mux := chi.NewMux()
		DeletePersonInfo(mux, &DeletePersonInfoRepoMock{err: errors.New("No file opening:(((")})

		req := httptest.NewRequest(http.MethodDelete, "/Info/{Id}", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		res := rec.Result()
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.FailNow()
		}
	})

	t.Run("sends accepted on successful finish", func(t *testing.T) {
		mux := chi.NewMux()
		DeletePersonInfo(mux, &DeletePersonInfoRepoMock{})

		req := httptest.NewRequest(http.MethodDelete, "/Info/1", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		res := rec.Result()
		if res.StatusCode != http.StatusAccepted {
			t.FailNow()
		}
	})

}
