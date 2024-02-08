package models

type Person struct {
	ID          int64  `db:"id" json:"id" example:"11"`
	Name        string `db:"name" json:"name" example:"Ivan"`
	Surname     string `db:"surname" json:"surname" example:"Ivanov"`
	Patronymic  string `db:"patronymic" json:"patronymic" example:"Ivanovich"`
	Age         uint   `db:"age" json:"age" example:"45"`
	Gender      string `db:"gender" json:"gender" example:"male"`
	Nationality string `db:"nationality" json:"nationality" example:"RU"`
}

type RawPerson struct {
	Name       string `json:"name" db:"name" example:"Petr"`
	Surname    string `json:"surname" db:"surname" example:"Petrov"`
	Patronymic string `json:"patronymic" db:"patronymic" example:"Petrovich"`
}

type Params struct {
	Limit  int
	Offset int
}
