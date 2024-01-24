package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/binaryty/enricher-service/internal/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// ConfigProvider an interface implements config.
type ConfigProvider interface {
	MakePGURL() string
}

// Storage type of postgres storage.
type Storage struct {
	db *sqlx.DB
}

// New creates a new instance of storage.
func New(cfg ConfigProvider) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sqlx.Open("postgres", cfg.MakePGURL())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		db: db,
	}, nil
}

// Create creates person in storage and return id of last inserted record.
func (s *Storage) Create(ctx context.Context, person *models.Person) (int64, error) {
	const op = "storage.postgres.Create"

	query := `INSERT INTO persons (name, surname, patronymic, age, gender, nationality)
		VALUES(:name, :surname, :patronymic, :age, :gender, :nationality) ON CONFLICT DO NOTHING RETURNING id`

	stmt, err := s.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	var id int64
	err = stmt.GetContext(ctx, &id, person)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// SelectByID returns a person by ID.
func (s *Storage) SelectByID(ctx context.Context, personID int64) (*models.Person, error) {
	const op = "storage.postgres.SelectByID"

	query := `SELECT id, name, surname, patronymic, age, gender, nationality FROM persons WHERE id=$1`

	var person models.Person
	if err := s.db.GetContext(
		ctx,
		&person,
		query,
		personID,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &person, nil
}

// Update updates person.
func (s *Storage) Update(ctx context.Context, person *models.Person) error {
	const op = "storage.postgres.Update"

	query := `UPDATE persons 
SET name = $1,
    surname = $2,
    patronymic = $3,
    age = $4,
    gender = $5,
    nationality = $6
WHERE id = $7
`
	_, err := s.db.ExecContext(
		ctx,
		query,
		person.Name,
		person.Surname,
		person.Patronymic,
		person.Age,
		person.Gender,
		person.Nationality,
		person.ID,
	)

	return err
}

// SelectAll returns slice of persons.
func (s *Storage) SelectAll(ctx context.Context, limit int, offset int) ([]models.Person, error) {
	const op = "storage.postgres.SelectAll"

	query := `SELECT * FROM persons ORDER BY id LIMIT $1 OFFSET $2`

	rows, err := s.db.QueryxContext(ctx, query, limit, offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func() { _ = rows.Close() }()

	persons := make([]models.Person, 0)

	for rows.Next() {
		var person models.Person
		if err := rows.StructScan(&person); err != nil {
			return nil, fmt.Errorf("%s: %w", err)
		}

		persons = append(persons, person)
	}

	return persons, nil
}

// DeleteByID delete person from storage by ID.
func (s *Storage) DeleteByID(ctx context.Context, personID int64) error {
	const op = "storage.postgres.DeleteByID"

	query := `DELETE FROM persons WHERE id = $1`

	_, err := s.db.ExecContext(ctx, query, personID)

	return err
}
