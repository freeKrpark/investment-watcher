package repository

import "database/sql"

type SQLiteRepositry struct {
	Conn *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepositry {
	return &SQLiteRepositry{
		Conn: db,
	}
}

func (repo *SQLiteRepositry) Migrate() error {
	query := `
		create table if not exists accounts(
			id integer primary key autoincrement,
			name varchar not null);

		create table if not exists accounts_detail(
			id integer primary key autoincrement,
			reference_id integer not null,
			balance real not null,
			asset real not null,
			deposit real not null default 0,
			withdrawal real not null default 0,
			reg_dt integer not null
		);
	`

	_, err := repo.Conn.Exec(query)
	return err

}
