package server

import (
	"go-chi/chi-master"
	"net/http"
	"rest_ful_interface/handlers"
	"rest_ful_interface/storage"
)

type Server struct {
	Storer *storage.Storer
}

func (s *Server) Start() error {

	// r := chi.NewRouter()

	// r.Route("/Info", func(r chi.Router) {
	// 	r.Get("/", handlers.GetAllPersonsInfo)
	// 	r.Post("/", handlers.PostPersonInfo)

	// 	r.Route("/{Id}", func(r chi.Router) {
	// 		r.Get("/", handlers.GetPersonInfo)
	// 		r.Put("/", handlers.UpdatePersonInfo)
	// 		r.Delete("/", handlers.DeletePersonInfo)
	// 	})
	// })

	mux := chi.NewMux()

	handlers.GetAllPersonsInfo(mux, s.Storer)
	handlers.GetPersonInfo(mux, s.Storer)
	handlers.DeletePersonInfo(mux, s.Storer)
	handlers.PostPersonInfo(mux, s.Storer)
	handlers.UpdatePersonInfo(mux, s.Storer)

	return http.ListenAndServe(":8000", mux)
}
