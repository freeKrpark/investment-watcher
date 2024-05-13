package main

import (
	"fmt"
	"investmentwatcher/repository"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/dustin/go-humanize"
)

func (app *Config) totalTab() *fyne.Container {
	app.Total = app.getTotalSlice()
	app.TotalTable = app.getTotalTable()

	// 그리드의 행과 열 수를 설정합니다.
	cols := 1 // 한 열에 하나의 테이블만 있으므로 1입니다.

	// AdaptiveGrid를 생성합니다.
	accountsContainer := container.NewBorder(
		nil,
		nil,
		nil,
		nil,
		container.NewAdaptiveGrid(cols, app.TotalTable),
	)

	return accountsContainer
}

func (app *Config) getTotalTable() *widget.Table {
	t := widget.NewTable(
		func() (int, int) {
			return len(app.Total), len(app.Total[0])
		},
		func() fyne.CanvasObject {
			ctr := container.NewVBox(widget.NewLabel(""))
			return ctr
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {

			o.(*fyne.Container).Objects = []fyne.CanvasObject{
				widget.NewLabel(app.Total[i.Row][i.Col].(string)),
			}

		})

	colWidths := []float32{100, 100, 50, 100, 100, 110}
	for i := 0; i < len(colWidths); i++ {
		t.SetColumnWidth(i, colWidths[i])
	}
	return t
}

func (app *Config) getTotalSlice() [][]interface{} {
	var slice [][]interface{}
	accounts, err := app.currentTotal()
	if err != nil {
		app.ErrorLog.Println(err)
	}

	slice = append(slice, []interface{}{"Balance", "Asset", "Rate", "Deposit", "Withdrawal", "RegDt"})
	for _, x := range accounts {
		var currentRow []interface{}
		currentRow = append(currentRow, humanize.Comma(int64(x.Balance)))
		currentRow = append(currentRow, humanize.Comma(int64(x.Asset)))
		rateOfReturn := (x.Asset - x.Balance) / x.Balance * 100
		currentRow = append(currentRow, fmt.Sprintf("%.2f", rateOfReturn))
		currentRow = append(currentRow, humanize.Comma(int64(x.Deposit)))
		currentRow = append(currentRow, humanize.Comma(int64(x.Withdrawal)))
		currentRow = append(currentRow, x.RegDt.Format("2006-01-02"))
		slice = append(slice, currentRow)
	}
	return slice
}

func (app *Config) currentTotal() ([]repository.TotalDetail, error) {
	t, err := app.DB.SelectTotalDetailGropByRegDt()
	if err != nil {
		app.ErrorLog.Println(err)
		return nil, err
	}
	return t, nil
}
