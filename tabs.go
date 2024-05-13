package main

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

func (app *Config) getTabs() (*container.AppTabs, error) {
	tabs, err := app.makeTabs()
	if err != nil {
		app.ErrorLog.Println(err)
		return nil, err
	}
	app.Tabs = tabs
	return tabs, nil
}

func (app *Config) makeTabs() (*container.AppTabs, error) {
	tabs := app.makeDefaultTabs()
	accounts, err := app.DB.SelectAllAccounts()
	if err != nil {
		app.ErrorLog.Println(err)
		return nil, err
	}
	for _, x := range accounts {
		tab := app.getAccountTab(int(x.ID))
		tabs.Append(container.NewTabItemWithIcon(x.Name, theme.InfoIcon(), tab))
	}
	return tabs, nil
}

func (app *Config) makeDefaultTabs() *container.AppTabs {
	tabsTabContent := app.tabTab()
	graphTabContent := app.graphTab()
	totalTabContent := app.totalTab()
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Graph", theme.HomeIcon(), graphTabContent),
		container.NewTabItemWithIcon("Tabs", theme.ListIcon(), tabsTabContent),
		container.NewTabItemWithIcon("Total", theme.InfoIcon(), totalTabContent),
	)
	return tabs
}
