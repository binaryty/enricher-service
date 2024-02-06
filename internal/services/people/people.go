package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/binaryty/enricher-service/internal/models"
	"log"
	"log/slog"
)

type Service struct {
	log            *slog.Logger
	personProvider PersonProvider
	enricher       Enricher
}

type Enricher interface {
	Process(context.Context, models.RawPerson) (*models.Person, error)
}

type PersonProvider interface {
	Create(context.Context, models.Person) (int64, error)
	SelectByID(context.Context, int64) (*models.Person, error)
	Update(context.Context, *models.Person) error
	SelectAll(context.Context, models.Params) ([]models.Person, error)
	DeleteByID(context.Context, int64) error
}

var (
	ErrPersonNotFound = errors.New("person not found")
)

// New returns a new instance of People service.
func New(log *slog.Logger, personProvider PersonProvider, enricher Enricher) *Service {
	return &Service{
		log:            log,
		personProvider: personProvider,
		enricher:       enricher,
	}
}

// AddPerson enrich raw person data and save to storage.
func (s *Service) AddPerson(ctx context.Context, rawData models.RawPerson) (int64, error) {
	const op = "services.people.AddPerson"
	logger := s.log.With("operation", op)

	enrichResponse, err := s.enricher.Process(ctx, rawData)
	log.Println(enrichResponse)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	person := models.Person{
		Name:        rawData.Name,
		Surname:     rawData.Surname,
		Patronymic:  rawData.Patronymic,
		Age:         enrichResponse.Age,
		Gender:      enrichResponse.Gender,
		Nationality: enrichResponse.Nationality,
	}

	id, err := s.personProvider.Create(ctx, person)
	if err != nil {
		logger.Debug("can't create person", slog.String("[ERROR]", err.Error()))

		return -1, fmt.Errorf("%s: %w", op, err)
	}

	logger.Info("person successfully created", slog.Int64("ID", id))

	return id, nil
}

// SelectByID ...
func (s *Service) SelectByID(ctx context.Context, id int64) (*models.Person, error) {
	const op = "services.people.SelectByID"

	logger := s.log.With("operation", op)

	person, err := s.personProvider.SelectByID(ctx, id)
	if err != nil {
		logger.Debug("can't query person", slog.String("[ERROR]", err.Error()))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.Debug("person find successfully")

	return person, nil
}

// SelectAll ...
func (s *Service) SelectAll(ctx context.Context, params models.Params) ([]models.Person, error) {
	const op = "services.people.SelectAll"
	logger := s.log.With("operation", op)

	persons, err := s.personProvider.SelectAll(ctx, params)
	if err != nil {
		if errors.Is(err, ErrPersonNotFound) {
			logger.Debug("not found", slog.String("[ERROR]", err.Error()))

			return nil, ErrPersonNotFound
		}
		logger.Debug("can't query persons", slog.String("[ERROR]", err.Error()))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.Debug("all persons get successfully")

	return persons, nil
}

// DeleteByID ...
func (s *Service) DeleteByID(ctx context.Context, id int64) error {
	const op = "services.person.DeleteByID"
	logger := s.log.With("operation", op)

	if err := s.personProvider.DeleteByID(ctx, id); err != nil {
		logger.Debug("can't delete person", slog.String("[ERROR]", err.Error()))

		return err
	}
	logger.Debug("person successfully deleted", slog.Int64("ID", id))

	return nil
}

// Update ...
func (s *Service) Update(ctx context.Context, params *models.Person) error {

	return s.personProvider.Update(ctx, params)
}
