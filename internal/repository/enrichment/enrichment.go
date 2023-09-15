package enrichment

import (
	"effective_mobile/internal/domain"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type AgeData struct {
	Count uint64 `json:"count"`
	Name  string `json:"name"`
	Age   uint8  `json:"age"`
}

type GenderData struct {
	Count             uint64  `json:"count"`
	Name              string  `json:"name"`
	Gender            string  `json:"gender"`
	ProbabilityGender float32 `json:"probability"`
}

type NationalizeData struct {
	Count   uint64        `json:"count"`
	Name    string        `json:"name"`
	Country []CountryData `json:"country"`
}

type CountryData struct {
	ID          string  `json:"country_id"`
	Probability float32 `json:"probability"`
}

type Repository struct {
	ageUrl         string
	genderUrl     string
	nationalityUrl string
}

func New(ageUrl, genderUrl, nationalityUrl string) (domain.EnrichmentRepository, error) {
	if ageUrl == "" || genderUrl == "" || nationalityUrl == "" {
		return nil, errors.New("empty parameter")
	}
	return &Repository{
		ageUrl:         ageUrl,
		genderUrl:      genderUrl,
		nationalityUrl: nationalityUrl,
	}, nil
}

func (r *Repository) Age(name string) (uint8, error) {
	var ageData AgeData

	age, err := httpGet(r.ageUrl, name)
	if err != nil {
		return 0, err
	}

	err = json.Unmarshal(age, &ageData)
	if err != nil {
		return 0, err
	}

	return ageData.Age, nil
}

func (r *Repository) Country(name string) (string, error) {
	var nationalityData NationalizeData

	nationality, err := httpGet(r.nationalityUrl, name)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(nationality, &nationalityData)
	if err != nil {
		return "", err
	}

	country := chooseCountryWithMaxProbability(nationalityData.Country)
	return country, nil
}

func (r *Repository) Gender(name string) (string, error) {
	var genderData GenderData

	gender, err := httpGet(r.genderUrl, name)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(gender, &genderData)
	if err != nil {
		return "", err
	}

	return genderData.Gender, nil
}

func chooseCountryWithMaxProbability(countires []CountryData) string {
	var res string
	var maxProbability float32
	for _, country := range countires {
		if maxProbability <= country.Probability {
			maxProbability = country.Probability
			res = country.ID
		}
	}
	return res
}

func httpGet(link, name string) ([]byte, error) {
	resp, err := http.Get(link + name)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
