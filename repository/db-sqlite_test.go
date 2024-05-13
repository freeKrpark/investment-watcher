package repository

import (
	"testing"
	"time"
)

func TestSQLiteRepository_Migrate(t *testing.T) {
	err := testRepo.Migrate()
	if err != nil {
		t.Error("migrate failed:", err)
	}
}

func TestSQLiteRepository_InsertAccounts(t *testing.T) {
	a := Accounts{
		Name: "test",
	}

	result, err := testRepo.InsertAccount(a)
	if err != nil {
		t.Error("insert failed:", err)
	}

	if result.ID < 0 {
		t.Error("failed id sent back:", result.ID)
	}
}

func TestSQLiteRepository_InsertAccountDetail(t *testing.T) {
	accounts := AccountDetail{
		ReferenceId: 1,
		Balance:     18857320,
		Asset:       17648090,
		Deposit:     0,
		Withdrawal:  0,
		RegDt:       time.Now(),
	}

	result, err := testRepo.InsertAccountDetail(accounts)
	if err != nil {
		t.Error("insert failed:", err)
	}

	if result.ID < 0 {
		t.Error("failed id sent back:", result.ID)
	}
}

func TestSQLiterRepository_SelectAllAccountDetailsByReferenceId(t *testing.T) {
	accounts, err := testRepo.SelectAllAccountDetailsByRefernecId(1)
	if err != nil {
		t.Error("get by referenceId failed:", err)
	}

	if len(accounts) != 1 {
		t.Error("expected for 1 but got", len(accounts))
	}
	if accounts[0].Asset != 17648090 {
		t.Error("wrong assets returned; expected 17648090, but got", accounts[0].Asset)
	}
}

func TestSQLiterRepository_SelectAccountDetailById(t *testing.T) {
	accounts, err := testRepo.SelectAccountDetailById(1)
	if err != nil {
		t.Error("get by id failed:", err)
	}

	if accounts.Asset != 17648090 {
		t.Error("wrong assets returned; expected 17648090, but got", accounts.Asset)
	}

	_, err = testRepo.SelectAccountDetailById(2)
	if err == nil {
		t.Error("get one returned value for non-existence")
	}
}

func TestSQLiterRepository_SelectTotalDetailGropByRegDt(t *testing.T) {
	totals, err := testRepo.SelectTotalDetailGropByRegDt()
	if err != nil {
		t.Error("group by failed:", err)
	}

	if len(totals) != 1 {
		t.Error("expected for 1 but got", len(totals))
	}
	if totals[0].Asset != 17648090 {
		t.Error("wrong assets returned; expected 17648090, but got", totals[0].Asset)
	}
}

func TestSQLiteRepository_SelectLatestAssetGroupbyAccount(t *testing.T) {
	assets, err := testRepo.SelectLatestAssetGroupbyAccount()
	if err != nil {
		t.Error("group by failed:", err)
	}
	if len(assets) != 1 {
		t.Error("expected for 1 but got", len(assets))
	}
	if assets[0].Asset != 17648090 {
		t.Error("wrong assets returned; expected 17648090, but got", assets[0].Asset)
	}
}

func TestSQLiterRepository_DeleteAccountDetails(t *testing.T) {
	err := testRepo.DeleteAccountDetail(1)
	if err != nil {
		t.Error("failed to delete Accounts", err)
	}

	err = testRepo.DeleteAccountDetail(2)
	if err == nil {
		t.Error("get one returned value for non-existence")
	}
}

func TestSQliteRepository_SelectAllAccounts(t *testing.T) {
	a, err := testRepo.SelectAllAccounts()
	if err != nil {
		t.Error("get all failed:", err)
	}

	if len(a) != 1 {
		t.Error("wrong number of rows returned; expected 1, but got:", len(a))
	}
}

func TestSQLiteRepository_SelectAccountsById(t *testing.T) {
	a, err := testRepo.SelectAccountsById(1)
	if err != nil {
		t.Error("get by id failed:", err)
	}

	if a.Name != "test" {
		t.Error("wrong name returned; expected test, but got", a.Name)
	}

	_, err = testRepo.SelectAccountsById(2)
	if err == nil {
		t.Error("got one returned value for non-existent id")
	}
}

func TestSQLiteRepository_DeleteAccountsById(t *testing.T) {
	err := testRepo.DeleteAccounts(1)
	if err != nil {
		t.Error("failed to delete account", err)
		if err != errDeleteFailed {
			t.Error("wrong error returned")
		}
	}

	err = testRepo.DeleteAccounts(2)
	if err == nil {
		t.Error("delete non-existent account")
	}
}
