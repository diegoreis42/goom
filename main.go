package main

import (
	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"
  "renderer"
)

func main() {
	app := tview.NewApplication()

	presentation := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("Welcome to the Application\n\nPress 'E' to Enter")

	presentation.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'E', 'e':
      renderer.Run()
			app.Stop()
		case 'Q', 'q':
			app.Stop()
		}
		return event
	})

	flex := tview.NewFlex().
		AddItem(presentation, 0, 1, true)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}

