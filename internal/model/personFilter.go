package model

type PersonFilter struct {
	Name          *string  `json:"name"`
	Surname       *string  `json:"surname"`
	Patronymic    *string  `json:"patronymic"`
	Nationalities []string `json:"nationalities"`
	Genders       []string `json:"countries"`

	NameLike       *string `json:"name_like"`
	SurnameLike    *string `json:"surname_like"`
	PatronymicLike *string `json:"patronymic_like"`
	AgeMin         *int    `json:"age_min"`
	AgeMax         *int    `json:"age_max"`
}
