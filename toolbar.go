package main

import (
	"fmt"
	"investmentwatcher/repository"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func (app *Config) getToolBar() *widget.Toolbar {
	toolbar := widget.NewToolbar(
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			app.addAccountsDialog()
		}),
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			app.addAccountDetailDialog()
		}),
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {}),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {}),
	)
	return toolbar
}

func (app *Config) addAccountDetailDialog() dialog.Dialog {
	accounts, err := app.currentTabs()
	if err != nil {
		app.ErrorLog.Println(err)
	}
	var options []string
	for _, account := range accounts {
		options = append(options, fmt.Sprintf("%d.%s", account.ID, account.Name))
	}
	selectEntry := widget.NewSelect(options, func(s string) {
	})
	addBalanceEntry := widget.NewEntry()
	addAssetEntry := widget.NewEntry()
	addDepositEntry := widget.NewEntry()
	addWithdrawalEntry := widget.NewEntry()
	addRegDt := widget.NewEntry()

	app.AddSelectEntry = selectEntry
	app.AddBalanceEntry = addBalanceEntry
	app.AddAssetEntry = addAssetEntry
	app.AddDepositEntry = addDepositEntry
	app.AddWithdrawalEntry = addWithdrawalEntry
	app.addRegDt = addRegDt

	dateValidator := func(s string) error {
		if _, err := time.Parse("2006-01-02", s); err != nil {
			return err
		}
		return nil
	}

	addRegDt.Validator = dateValidator
	isFloatValidator := func(s string) error {
		_, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		return nil
	}
	addBalanceEntry.Validator = isFloatValidator
	addAssetEntry.Validator = isFloatValidator
	addDepositEntry.Validator = isFloatValidator
	addWithdrawalEntry.Validator = isFloatValidator

	addRegDt.PlaceHolder = "YYYY-MM-DD"

	addForm := dialog.NewForm(
		"Add Account Detail",
		"Add",
		"Cancel",
		[]*widget.FormItem{
			{Text: "Accounts?", Widget: selectEntry},
			{Text: "Balance", Widget: addBalanceEntry},
			{Text: "Assets", Widget: addAssetEntry},
			{Text: "Deposit", Widget: addDepositEntry},
			{Text: "Withdrawal", Widget: addWithdrawalEntry},
			{Text: "RegDt", Widget: addRegDt}},
		func(b bool) {
			ids := strings.Split(selectEntry.Selected, ".")
			referenceId, _ := strconv.ParseInt(ids[0], 10, 64)
			balance, _ := strconv.ParseFloat(addBalanceEntry.Text, 64)
			assets, _ := strconv.ParseFloat(addAssetEntry.Text, 64)
			deposit, _ := strconv.ParseFloat(addDepositEntry.Text, 64)
			withdrawal, _ := strconv.ParseFloat(addWithdrawalEntry.Text, 64)
			regDt, _ := time.Parse("2006-01-02", addRegDt.Text)
			_, err := app.DB.InsertAccountDetail(repository.AccountDetail{
				ReferenceId: referenceId,
				Balance:     balance,
				Asset:       assets,
				Deposit:     deposit,
				Withdrawal:  withdrawal,
				RegDt:       regDt,
			})

			if err != nil {
				app.ErrorLog.Println(err)
			}

			// refresh
			app.refreshTabTable()
		}, app.MainWindow)
	// selectAccountEntry := widget.NewSelect()
	// app.AddAccountDetailEntry = selectAccountEntry
	addForm.Resize(fyne.Size{Width: 400})
	addForm.Show()
	return addForm
}

func (app *Config) addAccountsDialog() dialog.Dialog {
	addNameEntry := widget.NewEntry()
	app.AddAccountEntry = addNameEntry

	addForm := dialog.NewForm(
		"Add Account",
		"ADD",
		"Cancel",
		[]*widget.FormItem{
			{Text: "Name", Widget: addNameEntry},
		},
		func(b bool) {
			if b {
				name := addNameEntry.Text
				_, err := app.DB.InsertAccount(repository.Accounts{
					Name: name,
				})
				if err != nil {
					app.ErrorLog.Println(err)
				}
				app.refreshTabTable()
			}
		}, app.MainWindow)
	addForm.Resize(fyne.Size{Width: 400})
	addForm.Show()
	return addForm
}
