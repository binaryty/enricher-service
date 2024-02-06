package models

type Person struct {
	ID          int64  `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Surname     string `db:"surname" json:"surname"`
	Patronymic  string `db:"patronymic" json:"patronymic"`
	Age         uint   `db:"age" json:"age"`
	Gender      string `db:"gender" json:"gender"`
	Nationality string `db:"nationality" json:"nationality"`
}

type RawPerson struct {
	Name       string `json:"name" db:"name"`
	Surname    string `json:"surname" db:"surname"`
	Patronymic string `json:"patronymic" db:"patronymic"`
}

type Params struct {
	Limit  int `db:"limit" json:"limit"`
	Offset int `db:"offset" json:"offset"`
}
