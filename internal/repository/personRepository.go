package repository

import (
	"context"

	"github.com/ivanjabrony/personApi/internal/model"
)

type PersonRepository interface {
	Create(context.Context, *model.Person) (int, error)
	GetById(context.Context, int) (*model.Person, error)
	GetAll(context.Context) ([]model.Person, error)
	GetFiltered(context.Context, *model.PersonFilter) ([]model.Person, error)
	Update(context.Context, *model.Person) error
	DeleteById(context.Context, int) error
}
