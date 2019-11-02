package main

import "testing"

func TestNewBoard(t *testing.T) {
	want := "This is the title"
	b, err := NewBoard(want)
	if err != nil {
		t.Errorf("%v", err)
	}
	if b.Title != want {
		t.Errorf("Want: %v, have: %v", want, b.Title)
	}

}

func TestNewColumn(t *testing.T) {
	want := "This is the title"
	co, err := NewColumn(want)
	if err != nil {
		t.Errorf("%v", err)
	}
	if co.Title != want {
		t.Errorf("Want: %v, have: %v", want, co.Title)
	}
}

func TestNewCard(t *testing.T) {
	wantTitle := "This is the title"
	wantBody := "This is the body"

	c, err := NewCard(wantTitle, wantBody)
	if err != nil {
		t.Errorf("%v", err)
	}
	if c.Title != wantTitle {
		t.Errorf("Want: %v, have: %v", wantTitle, c.Title)
	}
	if c.Body != wantBody {
		t.Errorf("Want: %v, have: %v", wantBody, c.Body)
	}
}

func TestAddCardToColumn(t *testing.T) {
	Title := "This is the title"
	Body := "This is the body"
	c, err := NewCard(Title, Body)
	if err != nil {
		t.Errorf("%v", err)
	}

	coTitle := "This is the title"
	co, err := NewColumn(coTitle)
	if err != nil {
		t.Errorf("%v", err)
	}

	co.AddCard(c)
}
