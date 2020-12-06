package handlers

import (
	"errors"
	"go-chi/chi-master"
	"net/http"
	"net/http/httptest"
	"testing"
)

type UpdatePersonInfoRepoMock struct {
	err error
}

func (m *UpdatePersonInfoRepoMock) GetThemAll() error {
	return m.err
}

func TestUpdatePersonInfo(t *testing.T) {
	t.Run("sends bad gateway on file opening error", func(t *testing.T) {
		mux := chi.NewMux()
		UpdatePersonInfo(mux, &UpdatePersonInfoRepoMock{err: errors.New("No file opening:(((")})

		req := httptest.NewRequest(http.MethodPut, "/Info/1", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		res := rec.Result()
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.FailNow()
		}
	})

	t.Run("sends accepted on successful finish", func(t *testing.T) {
		mux := chi.NewMux()
		UpdatePersonInfo(mux, &UpdatePersonInfoRepoMock{})

		req := httptest.NewRequest(http.MethodPut, "/Info/1", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		res := rec.Result()
		if res.StatusCode != http.StatusAccepted {
			t.FailNow()
		}
	})

}
