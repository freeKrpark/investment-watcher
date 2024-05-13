package main

import (
	"investmentwatcher/repository"
	"os"
	"testing"

	"fyne.io/fyne/v2/test"
)

var testApp Config

func TestMain(m *testing.M) {
	a := test.NewApp()
	testApp.App = a
	testApp.MainWindow = a.NewWindow("")
	testApp.DB = repository.NewTestRepository()

	os.Exit(m.Run())

}
