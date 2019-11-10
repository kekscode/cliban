package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	uuid "github.com/satori/go.uuid"
)

// Board holding a list of columns
type Board struct {
	ID      uuid.UUID
	Title   string
	Columns []*Column
}

func NewBoard(title string) (*Board, error) {
	id := NewUUID()
	return &Board{id, title, nil}, nil
}

func (b *Board) AddColumn(co *Column) {
	b.Columns = append(b.Columns, co)
}

func (b *Board) Save() error {
	file, _ := json.MarshalIndent(b, "", " ")
	err := ioutil.WriteFile("cliban.json", file, 0644)
	if err != nil {
		err := errors.New(fmt.Sprintf("Could not write to disk %v", file))
		return err
	}
	return nil
}

func (b *Board) Load() {}

func (b *Board) GetLongestStackOfCards() int {
	longest := 0
	for _, co := range b.Columns {
		if len(co.Cards) >= longest {
			longest = len(co.Cards)
		}

	}
	return longest
}

func (b *Board) MoveCard(c *Card, direction string) error {
	for coi, co := range b.Columns {
		for cai, ca := range co.Cards {
			if ca.ID == c.ID {
				switch direction {
				case "left":
					if len(b.Columns[:coi]) != 0 {
						// Copy *card to slice on the left
						b.Columns[coi-1].Cards = append(b.Columns[coi-1].Cards, c)

						// Delete old *card pointer. This is complicated code, but avoids a memory leak.
						// See: https://github.com/golang/go/wiki/SliceTricks
						if cai < len(b.Columns[coi].Cards)-1 {
							copy(b.Columns[coi].Cards[cai:], b.Columns[coi].Cards[cai+1:])
						}
						b.Columns[coi].Cards[len(b.Columns[coi].Cards)-1] = nil
						b.Columns[coi].Cards = b.Columns[coi].Cards[:len(b.Columns[coi].Cards)-1]
						return nil
					}
				case "right":
					if len(b.Columns[coi+1:]) != 0 {
						// Copy *card to slice on the right
						b.Columns[coi+1].Cards = append(b.Columns[coi+1].Cards, c)

						// Delete old *card pointer. This is complicated code, but avoids a memory leak.
						// See: https://github.com/golang/go/wiki/SliceTricks
						if cai < len(b.Columns[coi].Cards)-1 {
							// FIXME: Somehow, this duplicates cards, at least visually(?):
							copy(b.Columns[coi].Cards[cai:], b.Columns[coi].Cards[cai+1:])
						}
						b.Columns[coi].Cards[len(b.Columns[coi].Cards)-1] = nil
						b.Columns[coi].Cards = b.Columns[coi].Cards[:len(b.Columns[coi].Cards)-1]
						log.Printf("%v \n%v\n", b.Columns, b.Columns[coi].Cards)
						return nil
					}
				case "up":
					if len(b.Columns[coi].Cards[:cai]) >= 1 {
						// Save card to move in new slice
						cardToMove := []*Card{}
						cardToMove = append(cardToMove, b.Columns[coi].Cards[cai])

						// Delete card to be moved from slice
						if cai < len(b.Columns[coi].Cards)-1 {
							copy(b.Columns[coi].Cards[cai:], b.Columns[coi].Cards[cai+1:])
						}
						b.Columns[coi].Cards[len(b.Columns[coi].Cards)-1] = nil
						b.Columns[coi].Cards = b.Columns[coi].Cards[:len(b.Columns[coi].Cards)-1]

						log.Printf("Stack: %v", b.Columns[coi].Cards)
						log.Printf("cardToMove: %v", cardToMove)

						// Insert card ("move")
						b.Columns[coi].Cards = append(b.Columns[coi].Cards[:cai-1], append(cardToMove, b.Columns[coi].Cards[cai-1:]...)...)
						log.Printf("New Stack: %v", b.Columns[coi].Cards)
						return nil
					}
				case "down":
					if len(b.Columns[coi].Cards[cai:]) >= 2 {
						// Save card to move in new slice
						cardToMove := []*Card{}
						cardToMove = append(cardToMove, b.Columns[coi].Cards[cai])

						// Delete card to be moved from slice
						if cai < len(b.Columns[coi].Cards)-1 {
							copy(b.Columns[coi].Cards[cai:], b.Columns[coi].Cards[cai+1:])
						}
						b.Columns[coi].Cards[len(b.Columns[coi].Cards)-1] = nil
						b.Columns[coi].Cards = b.Columns[coi].Cards[:len(b.Columns[coi].Cards)-1]

						log.Printf("Stack: %v", b.Columns[coi].Cards)
						log.Printf("cardToMove: %v", cardToMove)

						// Insert card ("move")
						b.Columns[coi].Cards = append(b.Columns[coi].Cards[:cai+1], append(cardToMove, b.Columns[coi].Cards[cai+1:]...)...)
						log.Printf("New Stack: %v", b.Columns[coi].Cards)
						return nil
					}
				}
			}
		}
	}
	msg := "Cound not move card"
	log.Println(msg)
	return errors.New(msg)

}

func (b *Board) AddCard(colid uuid.UUID, title, body string) (bool, error) {
	for coi, co := range b.Columns {
		if co.ID == colid {
			card, err := NewCard(title, body)
			if err != nil {
				log.Fatalf("Cound not create card: %v", err)
				return false, err
			}
			b.Columns[coi].Cards = append(b.Columns[coi].Cards, card)
		}
	}
	return true, nil
}

func (b *Board) AddCardBefore(id uuid.UUID) (bool, error) { return true, nil }
func (b *Board) AddCardAfter(id uuid.UUID) (bool, error)  { return true, nil }

func (b *Board) DeleteCard(id uuid.UUID) (bool, error) {
	for coi, co := range b.Columns {
		for cid, ca := range co.Cards {
			if ca.ID == id {
				b.Columns[coi].Cards = append(b.Columns[coi].Cards[:cid], b.Columns[coi].Cards[cid+1:]...)
				return true, nil
			}
		}
	}
	err := errors.New(fmt.Sprintf("Could not find card with %v", id))
	return false, err
}
func (b *Board) String() string {
	return b.Title
}
