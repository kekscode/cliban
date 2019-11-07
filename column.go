package main

import uuid "github.com/satori/go.uuid"

// Column holding a list of cards
type Column struct {
	ID    uuid.UUID
	Title string
	Cards []*Card
}

func NewColumn(title string) (*Column, error) {
	id := NewUUID()
	co := Column{id, title, nil}
	return &co, nil
}

func (co *Column) String() string {
	return co.Title
}
