package dto

type UpdatePersonDto struct {
	Id         int     `json:"id" example:"1" binding:"required" `
	Name       *string `json:"name" example:"Ivan" `
	Surname    *string `json:"surname" example:"Zabrodin" `
	Patronymic *string `json:"patronymic" example:"Vladimirovich"`
}
