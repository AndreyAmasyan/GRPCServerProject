package postgresStorage

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"test/internal/data"

	_ "github.com/jackc/pgx/v4/stdlib"
)

// PostgresStorage ...
type PostgreStorage struct {
	db        *sql.DB
	usersList []data.User
}

func New(ctx context.Context) (*PostgreStorage, error) {
	dsn := url.URL{
		Scheme: os.Getenv("POSTGRES_SCHEME"),
		Host:   os.Getenv("POSTGRES_HOST") + ":" + os.Getenv("POSTGRES_DBPORT"),
		User:   url.UserPassword(os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD")),
		Path:   os.Getenv("POSTGRES_NAME"),
	}

	q := dsn.Query()
	q.Add("sslmode", "disable")

	db, err := sql.Open("pgx", dsn.String())
	if err != nil {
		fmt.Println("PostgreStorage New sql.Open", err)
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		fmt.Println("PostgreStorage New db.PingContext", err)
		return nil, err
	}

	return &PostgreStorage{
		db: db,
	}, nil
}

func (s *PostgreStorage) AddUser(ctx context.Context, user data.User) error {
	q := `INSERT INTO users (id, name) VALUES ($1, $2)`

	if _, err := s.db.ExecContext(ctx, q, user.ID, user.UserName); err != nil {
		fmt.Println("PostgreStorage AddUser s.db.ExecContext", err)
		return err
	}

	return nil
}

func (s *PostgreStorage) FetchAllUsers(ctx context.Context) ([]data.User, error) {
	q := `SELECT id, name FROM users`

	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		fmt.Println("PostgreStorage FetchAllUsers s.db.QueryContext", err)
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	if err := rows.Err(); err != nil {
		fmt.Println("PostgreStorage FetchAllUsers rows.Err", err)
		return nil, err
	}

	for rows.Next() {
		var id int
		var name string

		if err := rows.Scan(&id, &name); err != nil {
			fmt.Println("PostgreStorage FetchAllUsers rows.Scan", err)
			return nil, err
		}

		user := data.User{
			ID:       id,
			UserName: name,
		}

		s.usersList = append(s.usersList, user)
	}

	return s.usersList, nil
}

func (s *PostgreStorage) DeleteUser(ctx context.Context, id int) error {
	q := `DELETE FROM users WHERE id=$1`

	if _, err := s.db.ExecContext(ctx, q, id); err != nil {
		fmt.Println("PostgreStorage DeleteUser s.db.ExecContext", err)
		return err
	}

	return nil
}

func (s *PostgreStorage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS users (id int NOT NULL UNIQUE, name varchar(255) NOT NULL)`

	if _, err := s.db.ExecContext(ctx, q); err != nil {
		fmt.Println("PostgreStorage Init s.db.ExecContext", err)
		return err
	}

	return nil
}

func (s *PostgreStorage) Close() error {
	if err := s.db.Close(); err != nil {
		fmt.Println("PostgreStorage Close s.db.Close", err)
		return err
	}

	return nil
}
