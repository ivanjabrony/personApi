package mapper

import (
	"github.com/ivanjabrony/personApi/internal/model"
	"github.com/ivanjabrony/personApi/internal/model/dto"
)

func MapFromNewPersonDto(dto *dto.NewPersonDto) *model.Person {
	if dto != nil {
		return &model.Person{
			Name:       dto.Name,
			Surname:    dto.Surname,
			Patronymic: dto.Patronymic,
		}
	}

	return nil
}

func MapFromPersonDto(dto *dto.PersonDto) *model.Person {
	if dto != nil {
		return &model.Person{
			Id:          dto.Id,
			Name:        dto.Name,
			Surname:     dto.Surname,
			Patronymic:  dto.Patronymic,
			Age:         dto.Age,
			Gender:      dto.Gender,
			Nationality: dto.Nationality,
		}
	}

	return nil
}

func MapFromUpdatePersonDto(dto *dto.UpdatePersonDto) *model.Person {
	var updName, updSurname string
	if dto.Name != nil {
		updName = *dto.Name
	}
	if dto.Surname != nil {
		updSurname = *dto.Surname
	}
	if dto != nil {
		return &model.Person{
			Id:         dto.Id,
			Name:       updName,
			Surname:    updSurname,
			Patronymic: dto.Patronymic,
		}
	}

	return nil
}

func MapToPersonDto(model *model.Person) *dto.PersonDto {
	if model != nil {
		return &dto.PersonDto{
			Id:          model.Id,
			Name:        model.Name,
			Surname:     model.Surname,
			Patronymic:  model.Patronymic,
			Age:         model.Age,
			Gender:      model.Gender,
			Nationality: model.Nationality,
		}
	}

	return nil
}

func MapToManyPersonDto(models ...model.Person) []dto.PersonDto {
	dtos := make([]dto.PersonDto, len(models))
	for i, v := range models {
		dtos[i] = *MapToPersonDto(&v)
	}

	return dtos
}

func MapFromManyPersonDto(dtos ...dto.PersonDto) []model.Person {
	models := make([]model.Person, len(dtos))
	for i, v := range dtos {
		models[i] = *MapFromPersonDto(&v)
	}

	return models
}
