package repository

import (
	"errors"
	"time"
)

var (
	errDeleteFailed = errors.New("delete failed")
)

type Repository interface {
	Migrate() error
	InsertAccount(accounts Accounts) (*Accounts, error)
	SelectAllAccounts() ([]Accounts, error)
	SelectAccountsById(id int) (*Accounts, error)
	DeleteAccounts(id int64) error
	InsertAccountDetail(accounts AccountDetail) (*AccountDetail, error)
	SelectAllAccountDetailsByRefernecId(referenceId int) ([]AccountDetail, error)
	SelectAccountDetailById(id int) (*AccountDetail, error)
	SelectTotalDetailGropByRegDt() ([]TotalDetail, error)
	DeleteAccountDetail(id int64) error
	SelectLatestAssetGroupbyAccount() ([]LastestAsset, error)
}

type LastestAsset struct {
	Name  string    `json:"string"`
	Asset float64   `json:"asset"`
	RegDt time.Time `json:"reg_dt"`
}

type Accounts struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type AccountDetail struct {
	ID          int64     `json:"id"`
	ReferenceId int64     `json:"reference_id"`
	Balance     float64   `json:"balance"`
	Asset       float64   `json:"asset"`
	Deposit     float64   `json:"deposit"`
	Withdrawal  float64   `json:"withdrawal"`
	RegDt       time.Time `json:"reg_dt"`
}

type TotalDetail struct {
	Balance    float64   `json:"balance"`
	Asset      float64   `json:"asset"`
	Deposit    float64   `json:"deposit"`
	Withdrawal float64   `json:"withdrawal"`
	RegDt      time.Time `json:"reg_dt"`
}

func (repo *SQLiteRepositry) InsertAccount(accounts Accounts) (*Accounts, error) {
	stmt := "insert into accounts (name) values(?)"

	res, err := repo.Conn.Exec(stmt, accounts.Name)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	accounts.ID = id
	return &accounts, nil
}

func (repo *SQLiteRepositry) SelectAllAccounts() ([]Accounts, error) {
	query := "select id, name from accounts order by id"
	rows, err := repo.Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []Accounts
	for rows.Next() {
		var accounts Accounts
		err := rows.Scan(
			&accounts.ID,
			&accounts.Name,
		)

		if err != nil {
			return nil, err
		}

		all = append(all, accounts)
	}
	return all, nil
}

func (repo *SQLiteRepositry) SelectAccountsById(id int) (*Accounts, error) {
	rows := repo.Conn.QueryRow("select id, name from accounts where id = ?", id)
	var accounts Accounts

	err := rows.Scan(
		&accounts.ID,
		&accounts.Name,
	)

	if err != nil {
		return nil, err
	}

	return &accounts, nil
}

func (repo *SQLiteRepositry) DeleteAccounts(id int64) error {
	res, err := repo.Conn.Exec("delete from accounts where id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errDeleteFailed
	}

	_, err = repo.Conn.Exec("delete from accounts_detail where reference_id = ?", id)
	if err != nil {
		return err
	}

	return nil
}
func (repo *SQLiteRepositry) InsertAccountDetail(accounts AccountDetail) (*AccountDetail, error) {
	stmt := "insert into accounts_detail(reference_id, balance, asset, deposit, withdrawal, reg_dt) values(?,?,?,?,?,?)"
	res, err := repo.Conn.Exec(stmt, accounts.ReferenceId, accounts.Balance, accounts.Asset, accounts.Deposit, accounts.Withdrawal, accounts.RegDt.Unix())
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	accounts.ID = id
	return &accounts, nil

}
func (repo *SQLiteRepositry) SelectAllAccountDetailsByRefernecId(referenceId int) ([]AccountDetail, error) {
	query := "select id, reference_id,balance, asset, deposit, withdrawal, reg_dt from accounts_detail where reference_id = ? order by id desc"
	rows, err := repo.Conn.Query(query, referenceId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []AccountDetail
	for rows.Next() {
		var accounts AccountDetail
		var unixTime int64
		err := rows.Scan(
			&accounts.ID,
			&accounts.ReferenceId,
			&accounts.Balance,
			&accounts.Asset,
			&accounts.Deposit,
			&accounts.Withdrawal,
			&unixTime,
		)
		if err != nil {
			return nil, err
		}
		accounts.RegDt = time.Unix(unixTime, 0)
		all = append(all, accounts)
	}
	return all, nil
}

func (repo *SQLiteRepositry) SelectTotalDetailGropByRegDt() ([]TotalDetail, error) {
	query := "select sum(balance), sum(asset), sum(deposit), sum(withdrawal), reg_dt from accounts_detail group by reg_dt order by reg_dt desc"
	rows, err := repo.Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []TotalDetail
	for rows.Next() {
		var totals TotalDetail
		var unixTime int64
		err := rows.Scan(
			&totals.Balance,
			&totals.Asset,
			&totals.Deposit,
			&totals.Withdrawal,
			&unixTime,
		)
		if err != nil {
			return nil, err
		}
		totals.RegDt = time.Unix(unixTime, 0)
		all = append(all, totals)
	}
	return all, nil
}

func (repo *SQLiteRepositry) SelectLatestAssetGroupbyAccount() ([]LastestAsset, error) {
	query := "SELECT a.name, ad.asset, ad.reg_dt FROM accounts_detail ad INNER JOIN accounts a ON ad.reference_id = a.id WHERE (a.name, ad.reg_dt) IN ( SELECT a.name, MAX(ad.reg_dt) FROM accounts_detail ad INNER JOIN accounts a ON ad.reference_id = a.id GROUP BY a.name);"

	rows, err := repo.Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []LastestAsset
	for rows.Next() {
		var LastestAsset LastestAsset
		var unixTime int64
		err := rows.Scan(
			&LastestAsset.Name,
			&LastestAsset.Asset,
			&unixTime,
		)

		if err != nil {
			return nil, err
		}
		LastestAsset.RegDt = time.Unix(unixTime, 0)
		all = append(all, LastestAsset)
	}
	return all, nil
}

func (repo *SQLiteRepositry) SelectAccountDetailById(id int) (*AccountDetail, error) {
	row := repo.Conn.QueryRow("select id, reference_id,balance, asset, deposit, withdrawal, reg_dt from accounts_detail where id = ?", id)
	var a AccountDetail
	var unixTime int64
	err := row.Scan(
		&a.ID,
		&a.ReferenceId,
		&a.Balance,
		&a.Asset,
		&a.Deposit,
		&a.Withdrawal,
		&unixTime,
	)
	if err != nil {
		return nil, err
	}
	a.RegDt = time.Unix(unixTime, 0)
	return &a, nil
}

func (repo *SQLiteRepositry) DeleteAccountDetail(id int64) error {
	res, err := repo.Conn.Exec("delete from accounts_detail where id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errDeleteFailed
	}
	return nil
}
