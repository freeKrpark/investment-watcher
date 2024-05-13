package repository

type TestRepository struct{}

func NewTestRepository() *TestRepository {
	return &TestRepository{}
}

func (repo *TestRepository) Migrate() error {
	return nil
}
func (repo *TestRepository) InsertAccount(accounts Accounts) (*Accounts, error) {
	return &accounts, nil
}
func (repo *TestRepository) SelectAllAccounts() ([]Accounts, error) {
	var all []Accounts
	a := Accounts{
		Name: "test1",
	}
	all = append(all, a)

	a = Accounts{
		Name: "test2",
	}
	all = append(all, a)
	return all, nil

}
func (repo *TestRepository) SelectAccountsById(id int) (*Accounts, error) {
	a := Accounts{
		Name: "test1",
	}
	return &a, nil
}
func (repo *TestRepository) DeleteAccounts(id int64) error {
	return nil
}
