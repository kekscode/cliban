package main

import (
	"log"
	"os"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	uuid "github.com/satori/go.uuid"
)

func NewUUID() uuid.UUID {
	id := uuid.NewV4()
	return id
}

func getSelectedCell(table *tview.Table) (r, c int) {
	r, c = table.GetSelection()
	return r, c
}

// Render the data in a tview table
func renderTableView(table *tview.Table, b *Board) {
	table.Clear()
	// rows
	for r := 0; r <= b.GetLongestStackOfCards(); r++ { // TODO: Try if GetLongestStack... can be substituted by b.GetRowCount method
		// columns
		for c := 0; c < len(b.Columns); c++ {

			// Create table with empty cells with non-selectable fields
			table.SetCell(r, c,
				tview.NewTableCell("").
					SetTextColor(tcell.ColorWhite).
					SetAlign(tview.AlignCenter).
					SetSelectable(false))

			// Set column titles for row with index == 0
			table.SetCell(0, c,
				tview.NewTableCell(b.Columns[c].Title).
					SetTextColor(tcell.ColorGreen).
					SetAlign(tview.AlignCenter).
					SetSelectable(false))

			// Read cards from Board/Columns/Cards struct and insert into table
			for idx, card := range b.Columns[c].Cards {
				table.SetCell(idx+1, c,
					tview.NewTableCell(card.Title).
						SetTextColor(tcell.ColorWhite).
						SetAlign(tview.AlignCenter).
						SetSelectable(true))
			}

		}
	}
}

func main() {

	f, err := os.OpenFile("cliban.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("Welcome to cliban")

	b, err := NewBoard("TestBoard")
	if err != nil {
		log.Fatalf("%v", err)
	}

	//FIXME: Move the column creation method to the board.go like i did
	//       for the creation of individual cards (see below)
	co1, err := NewColumn("TestColumn 1")
	if err != nil {
		log.Fatalf("%v", err)
	}

	co2, err := NewColumn("TestColumn 2")
	if err != nil {
		log.Fatalf("%v", err)
	}

	co3, err := NewColumn("TestColumn 3")
	if err != nil {
		log.Fatalf("%v", err)
	}

	co4, err := NewColumn("TestColumn 4")
	if err != nil {
		log.Fatalf("%v", err)
	}

	co5, err := NewColumn("TestColumn 5")
	if err != nil {
		log.Fatalf("%v", err)
	}

	b.AddColumn(co1)
	b.AddColumn(co2)
	b.AddColumn(co3)
	b.AddColumn(co4)
	b.AddColumn(co5)

	b.AddCard(co1.ID, "Title 1", "Body 1")
	b.AddCard(co2.ID, "Title 2", "Body 2")
	b.AddCard(co2.ID, "Title 3", "Body 3")
	b.AddCard(co3.ID, "Title 4", "Body 4")
	b.AddCard(co4.ID, "Title 5", "Body 5")
	b.AddCard(co5.ID, "Title 6", "Body 6")
	b.AddCard(co5.ID, "Title 7", "Body 7")
	b.AddCard(co5.ID, "Title 8", "Body 8")

	app := tview.NewApplication()

	table := tview.NewTable().SetBorders(true)
	table.SetTitle(b.Title)
	table.SetBorders(true)
	table.SetSelectable(true, true)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		key := event.Key()
		if key == tcell.KeyESC {
			app.Stop()
		}

		if key == tcell.KeyCtrlC {
			app.Stop()
		}

		// D deletes the selected card
		if key == tcell.KeyRune {
			key = tcell.Key(event.Rune())
			if key == 68 { // 68 == ascii("D")
				ro, co := getSelectedCell(table)

				// -1 because table head is ro[0]
				_, err := b.DeleteCard(b.Columns[co].Cards[ro-1].ID)
				if err != nil {
					log.Printf("Error: %v", err)
				}
			} else if key == 99 { // 99 == ascii("c")
				_, co := getSelectedCell(table)

				// -1 because table head is ro[0]
				_, err := b.AddCard(b.Columns[co].ID, "AddCardTest", "AddCardBodyTest")
				if err != nil {
					log.Printf("Error: %v", err)
				}
			} else if key == 72 { // 72 == ascii("H")
				ro, co := getSelectedCell(table)
				log.Printf("got cell: %v, %v", ro, co)

				// -1 because table head is ro[0]
				err := b.MoveCard(b.Columns[co].Cards[ro-1], "left")
				if err != nil {
					log.Printf("Error: %v", err)
				} else {
					log.Printf("moving cursor")
					if co != 0 {
						destRowLen := len(b.Columns[co-1].Cards)
						log.Printf("destRowLen: %v", destRowLen)
						table.Select(destRowLen, co-1)
						log.Printf("done moving cursor")
					}
				}
			} else if key == 76 { // 76 == ascii("L")
				ro, co := getSelectedCell(table)

				// -1 because table head is ro[0]
				err := b.MoveCard(b.Columns[co].Cards[ro-1], "right")
				if err != nil {
					log.Printf("Error: %v", err)
				} else {
					log.Printf("moving cursor")
					if co != len(b.Columns) {
						destRowLen := len(b.Columns[co+1].Cards)
						log.Printf("destRowLen: %v", destRowLen)
						table.Select(destRowLen, co+1)
						log.Printf("done moving cursor")
					}
				}
			}
		}

		// Always save before redrawing UI
		func(app *tview.Application) {
			app.QueueUpdateDraw(func() {
				//				ro, co := getSelectedCell(table)
				//				log.Printf("%v, %v", ro, co)
			})

			err := b.Save()
			if err != nil {
				log.Fatalf("Error saving: %v", err)
			}

		}(app)

		// Render table to show changes
		renderTableView(table, b)

		// Pass on event to framework
		return event
	})

	if err := app.SetRoot(table, true).Run(); err != nil {
		panic(err)
	}
}
