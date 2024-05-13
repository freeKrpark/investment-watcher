package main

import (
	"fmt"
	"investmentwatcher/repository"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/dustin/go-humanize"
)

func (app *Config) getAccountTab(id int) *fyne.Container {
	accountsTable := app.getAccountTable(id)
	// 그리드의 행과 열 수를 설정합니다.
	cols := 1 // 한 열에 하나의 테이블만 있으므로 1입니다.

	// AdaptiveGrid를 생성합니다.
	accountsContainer := container.NewBorder(
		nil,
		nil,
		nil,
		nil,
		container.NewAdaptiveGrid(cols, accountsTable),
	)

	return accountsContainer
}

func (app *Config) getAccountTable(id int) *widget.Table {
	accounts := app.getAccountSlice(id)
	t := widget.NewTable(
		func() (int, int) {
			return len(accounts), len(accounts[0])
		},
		func() fyne.CanvasObject {
			ctr := container.NewVBox(widget.NewLabel(""))
			return ctr
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			if i.Col == (len(accounts[0])-1) && i.Row != 0 {
				// last cell - put in a button
				w := widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), func() {
					dialog.ShowConfirm("Delete?", "", func(deleted bool) {
						if deleted {
							id, _ := strconv.Atoi(accounts[i.Row][0].(string))
							app.InfoLog.Println(id)
							err := app.DB.DeleteAccountDetail(int64(id))
							if err != nil {
								app.ErrorLog.Println(err)
							}
						}
						// refresh the hodlings table
						app.refreshTabTable()
					}, app.MainWindow)

				})
				w.Importance = widget.HighImportance

				o.(*fyne.Container).Objects = []fyne.CanvasObject{
					w,
				}
			} else {
				o.(*fyne.Container).Objects = []fyne.CanvasObject{
					widget.NewLabel(accounts[i.Row][i.Col].(string)),
				}
			}
		})

	colWidths := []float32{30, 100, 100, 50, 100, 100, 110, 110}
	for i := 0; i < len(colWidths); i++ {
		t.SetColumnWidth(i, colWidths[i])
	}
	return t
}

func (app *Config) getAccountSlice(id int) [][]interface{} {
	var slice [][]interface{}
	accounts, err := app.getAccont(id)
	if err != nil {
		app.ErrorLog.Println(err)
	}
	slice = append(slice, []interface{}{"ID", "Balance", "Asset", "Rate", "Deposit", "Withdrawal", "RegDt", "Delete?"})
	for _, x := range accounts {
		var currentRow []interface{}
		currentRow = append(currentRow, strconv.FormatInt(x.ID, 10))
		currentRow = append(currentRow, humanize.Comma(int64(x.Balance)))
		currentRow = append(currentRow, humanize.Comma(int64(x.Asset)))
		rateOfReturn := (x.Asset - x.Balance) / x.Balance * 100
		currentRow = append(currentRow, fmt.Sprintf("%.2f", rateOfReturn))
		currentRow = append(currentRow, humanize.Comma(int64(x.Deposit)))
		currentRow = append(currentRow, humanize.Comma(int64(x.Withdrawal)))
		currentRow = append(currentRow, x.RegDt.Format("2006-01-02"))
		currentRow = append(currentRow, widget.NewButton("Delete", func() {

		}))
		slice = append(slice, currentRow)

	}
	return slice
}

func (app *Config) getAccont(id int) ([]repository.AccountDetail, error) {
	accounts, err := app.DB.SelectAllAccountDetailsByRefernecId(id)
	if err != nil {
		app.ErrorLog.Println(err)
		return nil, err
	}
	return accounts, nil
}
