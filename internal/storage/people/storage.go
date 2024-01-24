package peoplec

import (
	"context"
	"github.com/binaryty/enricher-service/internal/models"
)

type PeopleStorage interface {
	Create(context.Context, *models.Person) (int64, error)
	SelectByID(context.Context, int64) (*models.Person, error)
	Update(context.Context, *models.Person) error
	SelectAll(context.Context, int, int) ([]models.Person, error)
	DeleteByID(context.Context, int64) error
}
