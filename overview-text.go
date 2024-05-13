package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"github.com/dustin/go-humanize"
)

func (app *Config) getOverviewText() (*canvas.Text, *canvas.Text, *canvas.Text, *canvas.Text) {

	var balance, current, rate, change *canvas.Text
	overview, err := app.GetOverview()
	if err != nil {
		grey := color.NRGBA{R: 155, G: 155, B: 155, A: 255}
		balance = canvas.NewText("Balance : Unknown", grey)
		current = canvas.NewText("Asset : Unknown", grey)
		rate = canvas.NewText("Rate : Unknown", grey)
		change = canvas.NewText("Change : Unknown", grey)
	} else {
		changeColor := color.NRGBA{R: 0, G: 180, B: 0, A: 255}
		displayColor := color.NRGBA{R: 0, G: 180, B: 0, A: 255}
		fmt.Println(overview.Change)
		if overview.Change < 0 {
			changeColor = color.NRGBA{R: 188, G: 0, B: 0, A: 255}
		}

		if overview.Rate < 0 {
			displayColor = color.NRGBA{R: 188, G: 0, B: 0, A: 255}
		}

		balanceTxt := fmt.Sprintf("Balance : %s", humanize.Comma(int64(overview.OriginalBalance)))
		currentTxt := fmt.Sprintf("Asset : %s", humanize.Comma(int64(overview.CurrentAsset)))
		rateTxt := fmt.Sprintf("Rate: %.2f%%", overview.Rate)
		changeTxt := fmt.Sprintf("Change: %.2f%%", overview.Change)

		balance = canvas.NewText(balanceTxt, color.Black)
		current = canvas.NewText(currentTxt, displayColor)
		rate = canvas.NewText(rateTxt, displayColor)
		change = canvas.NewText(changeTxt, changeColor)
	}

	balance.Alignment = fyne.TextAlignLeading
	current.Alignment = fyne.TextAlignCenter
	rate.Alignment = fyne.TextAlignCenter
	change.Alignment = fyne.TextAlignTrailing
	return balance, current, rate, change
}
