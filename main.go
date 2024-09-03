package main

import (
	"log"

	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()


	if err := app.SetRoot(layout.GetGrid(), true).EnableMouse(true).Run(); err != nil {
			log.Fatalf("Error running application: %v", err)
		}
}
