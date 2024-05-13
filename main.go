package main

import (
	"database/sql"
	"investmentwatcher/repository"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	_ "github.com/glebarez/go-sqlite"
)

type Config struct {
	App                   fyne.App
	InfoLog               *log.Logger
	ErrorLog              *log.Logger
	DB                    repository.Repository
	MainWindow            fyne.Window
	OverviewContainer     *fyne.Container
	ToolBar               *widget.Toolbar
	Tabs                  *container.AppTabs
	TabsList              [][]interface{}
	TabsTable             *widget.Table
	GraphContainer        *fyne.Container
	AddAccountEntry       *widget.Entry
	AddAccountDetailEntry *widget.Entry
	Total                 [][]interface{}
	TotalTable            *widget.Table

	AddSelectEntry     *widget.Select
	AddBalanceEntry    *widget.Entry
	AddAssetEntry      *widget.Entry
	AddDepositEntry    *widget.Entry
	AddWithdrawalEntry *widget.Entry
	addRegDt           *widget.Entry
}

func main() {
	var investApp Config

	// create a fyne application
	fyneApp := app.NewWithID("freeKraprk.investmentwatcher")
	investApp.App = fyneApp

	// create our loggers
	investApp.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	investApp.ErrorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// open a connection to the database
	sqlDB, err := investApp.connectSQL()
	if err != nil {
		log.Panic(err)
	}

	// create a database repository
	investApp.setupDB(sqlDB)

	// create and size a fyne window
	investApp.MainWindow = investApp.App.NewWindow("Investment Watcher 한글")
	investApp.MainWindow.Resize(fyne.NewSize(770, 410))
	investApp.MainWindow.SetFixedSize(true)
	investApp.makeUI()
	investApp.MainWindow.SetMaster()

	// show and run the application
	investApp.MainWindow.ShowAndRun()

}

func (app *Config) connectSQL() (*sql.DB, error) {
	path := ""
	if os.Getenv("DB_PATH") != "" {
		path = os.Getenv("DB_PATH")
	} else {
		path = app.App.Storage().RootURI().Path() + "/sql.db"
		app.InfoLog.Print("db in:", path)
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (app *Config) setupDB(sqlDB *sql.DB) {
	app.DB = repository.NewSQLiteRepository(sqlDB)
	err := app.DB.Migrate()
	if err != nil {
		app.ErrorLog.Println(err)
		log.Panic()
	}
}
