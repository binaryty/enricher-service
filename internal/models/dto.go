package models

const (
	Male   = "male"
	Female = "female"
)

type RawPerson struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
}

type PersonsRequest struct {
	Name        string `json:"name" db:"name"`
	Surname     string `json:"surname" db:"surname"`
	Patronymic  string `json:"patronymic,omitempty" db:"patronymic"`
	Age         int    `json:"age" db:"age"`
	Gender      string `json:"gender" db:"gender"`
	Nationality string `json:"nationality" db:"nationality"`
	Limit       int    `json:"limit" db:"limit"`
	Offset      int    `json:"offset" db:"offset"`
}

type UpdatePersonRequest struct {
	ID          int64  `json:"id" db:"id"`
	Surname     string `json:"surname" db:"surname"`
	Patronymic  string `json:"patronymic,omitempty" db:"patronymic"`
	Age         int    `json:"age" db:"age"`
	Gender      string `json:"gender" db:"gender"`
	Nationality string `json:"nationality" db:"nationality"`
}

type UpdatePersonResponse struct {
	ID          int64  `json:"id" db:"id"`
	Surname     string `json:"surname" db:"surname"`
	Patronymic  string `json:"patronymic,omitempty" db:"patronymic"`
	Age         int    `json:"age" db:"age"`
	Gender      string `json:"gender" db:"gender"`
	Nationality string `json:"nationality" db:"nationality"`
}
