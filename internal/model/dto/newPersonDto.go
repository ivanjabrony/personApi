package dto

type NewPersonDto struct {
	Name       string  `json:"name" example:"Ivan" binding:"required"`
	Surname    string  `json:"surname" example:"Zabrodin" binding:"required"`
	Patronymic *string `json:"patronymic" example:"Vladimirovich"`
}
