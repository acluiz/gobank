package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountByID(id int) (*Account, error)
	DeleteAccount(id int) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=admin password=admin dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Init() error {
	return s.createAccountTable()
}

func (s *PostgresStore) createAccountTable() error {
	query := `CREATE TABLE if not exists account (
		id	serial	primary	key,
		first_name	varchar(50),
		last_name	varchar(50),
		number		serial,
		balance		serial,
		created_at	timestamp
	)`

	_, err := s.db.Exec(query)
	
	return err
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `
		INSERT INTO account (first_name, last_name, number, balance, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := s.db.Query(
		query, 
		acc.FirstName, 
		acc.LastName, 
		acc.Number, 
		acc.Balance, 
		acc.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	query := `SELECT * FROM account WHERE id = $1`

	rows, err := s.db.Query(query, id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account not found")
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	query := `SELECT * FROM account`

	rows, err := s.db.Query(query)

	if err != nil {
		return nil, err
	}

	accounts := []*Account{}

	for rows.Next() {
		acc, err := scanIntoAccount(rows)

		if err != nil {
			return nil, err
		}

		accounts = append(accounts, acc)
	}

	return accounts, nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	query := `DELETE FROM account WHERE id = $1`
	_, err := s.db.Query(query, id)

	return err
}

func scanIntoAccount(rows *sql.Rows)(*Account, error) {
	acc := new(Account)

	err := rows.Scan(
		&acc.ID, 
		&acc.FirstName, 
		&acc.LastName, 
		&acc.Number, 
		&acc.Balance, 
		&acc.CreatedAt,
	)

	return acc, err
}