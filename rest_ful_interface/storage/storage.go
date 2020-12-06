package storage

import (
	"fmt"
)

type Storer struct{}

func (s *Storer) GetThemAll() error {
	fmt.Printf("You got them!")
	return nil
}

func (s *Storer) GetHim(id string) error {
	fmt.Printf("You got him!")
	return nil
}

func (s *Storer) DeleteHim(id string) error {
	fmt.Printf("You DELETED HIM!")
	return nil
}

func (s *Storer) UpdateHim(id string) error {
	fmt.Printf("You  updated him!")
	return nil
}
