package dto

type PersonDto struct {
	Id         int     `json:"id" example:"1"`
	Name       string  `json:"name" example:"Ivan"`
	Surname    string  `json:"surname" example:"Zabrodin"`
	Patronymic *string `json:"patronymic" example:"Vladimirovich"`

	Age         *int    `json:"age" example:"21"`
	Gender      *string `json:"gender" example:"male"`
	Nationality *string `json:"nationality" example:"russian"`
}
