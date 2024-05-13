package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func (app *Config) makeUI() {
	// get overview
	balance, previous, current, change := app.getOverviewText()

	overviewContent := container.NewGridWithColumns(4,
		balance,
		previous,
		current,
		change,
	)

	app.OverviewContainer = overviewContent

	// get toolbar
	toolbar := app.getToolBar()
	app.ToolBar = toolbar

	// get tabs
	tabs, err := app.getTabs()
	if err != nil {
		app.ErrorLog.Println(err)
	}

	tabs.SetTabLocation(container.TabLocationTop)
	tabs.Resize(fyne.Size{
		Width:  770,
		Height: 410,
	})
	app.Tabs = tabs

	finalContent := container.NewVBox(overviewContent, toolbar, tabs)
	app.MainWindow.SetContent(finalContent)
}

func (app *Config) refreshTabTable() {
	app.InfoLog.Println("Refresh TabTable")
	app.TabsList = app.getTabSlice()
	app.TabsTable = app.getTabTable()
	// get tabs
	tabs, err := app.getTabs()
	if err != nil {
		app.ErrorLog.Println(err)
	}

	app.Tabs.SetItems(tabs.Items)
	app.InfoLog.Println("Refresh Tabs:", len(app.Tabs.Items))
	finalContent := container.NewVBox(app.OverviewContainer, app.ToolBar, app.Tabs)
	app.MainWindow.SetContent(finalContent)
	app.MainWindow.Content().Refresh()
}
