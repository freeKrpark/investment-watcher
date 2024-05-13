package main

import (
	"investmentwatcher/repository"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func (app *Config) tabTab() *fyne.Container {
	app.TabsList = app.getTabSlice()
	app.TabsTable = app.getTabTable()

	tabsContainer := container.NewBorder(
		nil,
		nil,
		nil,
		nil,

		container.NewAdaptiveGrid(1, app.TabsTable),
	)

	return tabsContainer
}

func (app *Config) getTabTable() *widget.Table {
	t := widget.NewTable(
		func() (int, int) {
			return len(app.TabsList), len(app.TabsList[0])
		},
		func() fyne.CanvasObject {
			ctr := container.NewVBox(widget.NewLabel(""))
			return ctr
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			if i.Col == (len(app.TabsList[0])-1) && i.Row != 0 {
				// last cell - put in a button
				w := widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), func() {
					dialog.ShowConfirm("Delete?", "", func(deleted bool) {
						if deleted {
							id, _ := strconv.Atoi(app.TabsList[i.Row][0].(string))
							app.InfoLog.Println(id)
							err := app.DB.DeleteAccounts(int64(id))
							if err != nil {
								app.ErrorLog.Println(err)
							}
						}
						// refresh the Tabs table
						app.refreshTabTable()
					}, app.MainWindow)

				})
				w.Importance = widget.HighImportance

				o.(*fyne.Container).Objects = []fyne.CanvasObject{
					w,
				}
			} else {
				o.(*fyne.Container).Objects = []fyne.CanvasObject{
					widget.NewLabel(app.TabsList[i.Row][i.Col].(string)),
				}
			}
		})

	colWidths := []float32{50, 600, 110}
	for i := 0; i < len(colWidths); i++ {
		t.SetColumnWidth(i, colWidths[i])

	}
	return t
}

func (app *Config) getTabSlice() [][]interface{} {
	var slice [][]interface{}
	tabs, err := app.currentTabs()
	if err != nil {
		app.ErrorLog.Println(err)
	}

	slice = append(slice, []interface{}{"ID", "Name", "Delete?"})
	for _, x := range tabs {
		var currentRow []interface{}
		currentRow = append(currentRow, strconv.FormatInt(x.ID, 10))
		currentRow = append(currentRow, x.Name)
		currentRow = append(currentRow, widget.NewButton("Delete", func() {}))
		slice = append(slice, currentRow)
	}
	return slice

}

func (app *Config) currentTabs() ([]repository.Accounts, error) {
	t, err := app.DB.SelectAllAccounts()
	if err != nil {
		app.ErrorLog.Println(err)
		return nil, err
	}
	return t, nil
}
