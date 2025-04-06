package service

import (
	"context"

	"github.com/ivanjabrony/personApi/internal/model"
	"github.com/ivanjabrony/personApi/internal/model/dto"
)

type PersonService interface {
	CreatePerson(context.Context, *dto.NewPersonDto) (int, error)
	GetPersonById(context.Context, int) (*dto.PersonDto, error)
	GetAllPersons(context.Context) ([]dto.PersonDto, error)
	GetPersonsFiltered(context.Context, *model.PersonFilter) ([]dto.PersonDto, error)
	UpdatePersonById(context.Context, *dto.UpdatePersonDto) error
	DeletePersonById(context.Context, int) error
}
