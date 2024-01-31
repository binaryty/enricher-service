package models

const (
	Male   = "male"
	Female = "female"
)

type AddPersonRequest struct {
	Name        string `json:"name" db:"name"`
	Surname     string `json:"surname" db:"surname"`
	Patronymic  string `json:"patronymic" db:"patronymic"`
	Age         uint   `json:"age" db:"age"`
	Gender      string `json:"gender" db:"gender"`
	Nationality string `json:"nationality" db:"nationality"`
}

type SelectPersonByIDRequest struct {
	ID int64 `json:"id"`
}

type SelectPersonByIDResponse struct {
	Name        string `json:"name" db:"name"`
	Surname     string `json:"surname" db:"surname"`
	Patronymic  string `json:"patronymic,omitempty" db:"patronymic"`
	Age         uint   `json:"age" db:"age"`
	Gender      string `json:"gender" db:"gender"`
	Nationality string `json:"nationality" db:"nationality"`
}

type SelectPersonsRequest struct {
	Limit  int `json:"limit" db:"limit"`
	Offset int `json:"offset" db:"offset"`
}

type SelectPersonsResponse struct {
	Name        string `json:"name" db:"name"`
	Surname     string `json:"surname" db:"surname"`
	Patronymic  string `json:"patronymic,omitempty" db:"patronymic"`
	Age         uint   `json:"age" db:"age"`
	Gender      string `json:"gender" db:"gender"`
	Nationality string `json:"nationality" db:"nationality"`
}

type UpdatePersonRequest struct {
	ID          int64  `json:"id" db:"id"`
	Surname     string `json:"surname" db:"surname"`
	Patronymic  string `json:"patronymic,omitempty" db:"patronymic"`
	Age         uint   `json:"age" db:"age"`
	Gender      string `json:"gender" db:"gender"`
	Nationality string `json:"nationality" db:"nationality"`
}

type UpdatePersonResponse struct {
	ID          int64  `json:"id" db:"id"`
	Surname     string `json:"surname" db:"surname"`
	Patronymic  string `json:"patronymic,omitempty" db:"patronymic"`
	Age         uint   `json:"age" db:"age"`
	Gender      string `json:"gender" db:"gender"`
	Nationality string `json:"nationality" db:"nationality"`
}

type EnreichRequest struct {
	Age         uint
	Gender      string
	Nationality string
}

type AgeResponse struct {
	Age uint `json:"age"`
}

type GenderResponse struct {
	Gender string `json:"gender"`
}

type NationalityResponse struct {
	CountryID   string  `json:"country_id"`
	Probability float32 `json:"probability"`
}

type NationalitysResponse struct {
	Country []NationalityResponse `json:"country"`
}
