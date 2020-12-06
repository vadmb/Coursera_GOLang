package main

import (
	"log"

	"rest_ful_interface/server"
	"rest_ful_interface/storage"
)

func main() {
	s := server.Server{Storer: &storage.Storer{}}
	if err := s.Start(); err != nil {
		log.Fatalln(err)
	}
}
