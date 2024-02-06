package enricher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/binaryty/enricher-service/internal/config"
	"github.com/binaryty/enricher-service/internal/models"
	"github.com/binaryty/enricher-service/internal/response"
	"net/http"
)

var (
	ErrHandleAge         = errors.New("can't handle age")
	ErrHandleGender      = errors.New("can't handle gender")
	ErrHandleNationality = errors.New("can't handle nationality")
)

type Enricher struct {
	ageAPI         string
	genderAPI      string
	nationalityAPI string
	client         *http.Client
}

// New create a new instance of Enricher.
func New(cfg *config.Config) *Enricher {
	return &Enricher{
		ageAPI:         cfg.API.Age,
		genderAPI:      cfg.API.Gender,
		nationalityAPI: cfg.API.Nationality,
		client:         &http.Client{},
	}
}

// Process processing enrich raw data.
func (e *Enricher) Process(ctx context.Context, rawData models.RawPerson) (*models.Person, error) {

	ageResp, err := e.handleAge(ctx, rawData.Name)
	if err != nil {
		return nil, err
	}

	genderResp, err := e.handleGender(ctx, rawData.Name)
	if err != nil {
		return nil, err
	}

	nationalityResp, err := e.handleNationality(ctx, rawData.Name)
	if err != nil {
		return nil, err
	}

	return &models.Person{
		Name:        rawData.Name,
		Surname:     rawData.Surname,
		Patronymic:  rawData.Patronymic,
		Age:         ageResp.Age,
		Gender:      genderResp.Gender,
		Nationality: nationalityResp.CountryID,
	}, nil
}

// handleAge get age from public API.
func (e *Enricher) handleAge(ctx context.Context, name string) (*response.AgeResponse, error) {
	const op = "services.enricher.handleAge"

	uri := fmt.Sprintf("%s?name=%s", e.ageAPI, name)

	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, ErrHandleAge
	}
	defer func() { _ = resp.Close }()

	ageResp := response.AgeResponse{}
	err = json.NewDecoder(resp.Body).Decode(&ageResp)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &ageResp, nil
}

// handleGender get gender from public API.
func (e *Enricher) handleGender(ctx context.Context, name string) (*response.GenderResponse, error) {
	const op = "services.enricher.handleGender"

	uri := fmt.Sprintf("%s?name=%s", e.genderAPI, name)

	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, ErrHandleGender
	}
	defer func() { _ = resp.Close }()

	genderResp := response.GenderResponse{}
	err = json.NewDecoder(resp.Body).Decode(&genderResp)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &genderResp, nil
}

// handleNationality get nationality from public API.
func (e *Enricher) handleNationality(ctx context.Context, name string) (*response.NationalityResponse, error) {
	const op = "services.enricher.handleNationality"

	uri := fmt.Sprintf("%s?name=%s", e.nationalityAPI, name)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, ErrHandleNationality
	}
	defer func() { _ = resp.Close }()

	nationResp := response.NationalitysResponse{}

	err = json.NewDecoder(resp.Body).Decode(&nationResp)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return getCountry(nationResp), nil
}

// getCountry ...
func getCountry(countries response.NationalitysResponse) *response.NationalityResponse {
	var maxProbability float32
	var idx int

	c := countries.Countries

	for i := 0; i < len(c); i++ {
		if c[i].Probability > maxProbability {
			maxProbability = c[i].Probability
			idx = i
		}
	}

	return &countries.Countries[idx]
}
