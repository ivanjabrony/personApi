package service_impl

import (
	"context"
	"log/slog"

	"github.com/ivanjabrony/personApi/internal/client"
	"github.com/ivanjabrony/personApi/internal/mapper"
	"github.com/ivanjabrony/personApi/internal/model"
	"github.com/ivanjabrony/personApi/internal/model/dto"
	"github.com/ivanjabrony/personApi/internal/repository"
)

type PersonService struct {
	personRepository  repository.PersonRepository
	logger            *slog.Logger
	ageclient         client.AgeClient
	genderClient      client.GenderClient
	nationalityClient client.NationalityClient
}

func NewPersonService(
	personRepository repository.PersonRepository,
	ageclient client.AgeClient,
	genderClient client.GenderClient,
	nationalityClient client.NationalityClient,
	logger *slog.Logger) *PersonService {
	return &PersonService{personRepository, logger, ageclient, genderClient, nationalityClient}
}

func (service *PersonService) CreatePerson(ctx context.Context, newPersonDto *dto.NewPersonDto) (int, error) {
	person := mapper.MapFromNewPersonDto(newPersonDto)
	service.logger.Debug("Start of person creation", slog.Any("data", *newPersonDto))

	age, err := service.ageclient.GetAgeByName(ctx, person.Name)
	if err != nil {
		service.logger.Error("Couldn't retrieve data from Age client", slog.String("Error", err.Error()))
	}
	gender, err := service.genderClient.GetGenderByName(ctx, person.Name)
	if err != nil {
		service.logger.Error("Couldn't retrieve data from Gender client", slog.String("Error", err.Error()))
	}
	nationality, err := service.nationalityClient.GetNationalityByName(ctx, person.Name)
	if err != nil {
		service.logger.Error("Couldn't retrieve data from Nationality client", slog.String("Error", err.Error()))
	}

	person.Age, person.Gender, person.Nationality = age, gender, nationality
	id, err := service.personRepository.Create(ctx, person)

	if err != nil {
		service.logger.Error("Repository error while creating", slog.String("Error", err.Error()))
		return -1, err
	}

	service.logger.Info("Person successfully created", slog.Int("ID", id))
	return id, nil
}

func (service *PersonService) GetPersonById(ctx context.Context, id int) (*dto.PersonDto, error) {
	service.logger.Debug("Start of reading person", slog.Int("ID", id))
	person, err := service.personRepository.GetById(ctx, id)

	if err != nil {
		service.logger.Error("Repository error while reading", slog.String("Error", err.Error()))
		return nil, err
	}

	service.logger.Info("Person successfully retrieved", slog.Int("ID", id))
	return mapper.MapToPersonDto(person), nil
}

func (service *PersonService) GetAllPersons(ctx context.Context) ([]dto.PersonDto, error) {
	service.logger.Debug("Start of reading all person")
	persons, err := service.personRepository.GetAll(ctx)

	if err != nil {
		service.logger.Error("Repository error while reading all", slog.String("Error", err.Error()))
		return nil, err
	}

	service.logger.Info("Persons successfully retrieved")
	return mapper.MapToManyPersonDto(persons...), nil
}

func (service *PersonService) GetPersonsFiltered(ctx context.Context, filter *model.PersonFilter) ([]dto.PersonDto, error) {
	service.logger.Debug("Start of person filtering", slog.Any("data", *filter))
	persons, err := service.personRepository.GetFiltered(ctx, filter)

	if err != nil {
		service.logger.Error("Repository error while filtering", slog.String("Error", err.Error()))
		return nil, err
	}

	service.logger.Info("Persons successfully filtered")
	return mapper.MapToManyPersonDto(persons...), nil
}

func (service *PersonService) UpdatePersonById(ctx context.Context, dto *dto.UpdatePersonDto) error {
	service.logger.Debug("Start of person updating", slog.Any("data", *dto))
	err := service.personRepository.Update(ctx, mapper.MapFromUpdatePersonDto(dto))

	if err != nil {
		service.logger.Error("Repository error while updating", slog.String("Error", err.Error()))
		return err
	}

	service.logger.Info("Person successfully updated")
	return nil
}

func (service *PersonService) DeletePersonById(ctx context.Context, id int) error {
	service.logger.Debug("Start of person deleting", slog.Int("ID", id))
	err := service.personRepository.DeleteById(ctx, id)

	if err != nil {
		service.logger.Error("Repository error while deleting", slog.String("Error", err.Error()))
		return err
	}

	service.logger.Info("Person successfully deleted")
	return nil
}
