package main

import uuid "github.com/satori/go.uuid"

// Card holding information
type Card struct {
	ID    uuid.UUID
	Title string
	Body  string
}

func NewCard(title, body string) (*Card, error) {
	id := NewUUID()
	return &Card{id, title, body}, nil
}

func (c *Card) Update() {}
func (c *Card) String() string {
	return c.Title
}
