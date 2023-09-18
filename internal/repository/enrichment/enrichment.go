package enrichment

import (
	"effective_mobile/internal/domain"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
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
	genderUrl      string
	nationalityUrl string
	log            *logrus.Logger
}

func New(ageUrl, genderUrl, nationalityUrl string, log *logrus.Logger) (domain.EnrichmentRepository, error) {
	if ageUrl == "" || genderUrl == "" || nationalityUrl == "" {
		log.Error("empty parameter")
		return nil, errors.New("empty parameter")
	}
	return &Repository{
		ageUrl:         ageUrl,
		genderUrl:      genderUrl,
		nationalityUrl: nationalityUrl,
		log:            log,
	}, nil
}

func (r *Repository) Age(name string) (uint8, error) {
	var ageData AgeData

	age, err := httpGet(r.ageUrl, name, r.log)
	if err != nil {
		r.log.Errorln(err)
		return 0, err
	}

	err = json.Unmarshal(age, &ageData)
	if err != nil {
		r.log.Errorln(err)
		return 0, err
	}

	return ageData.Age, nil
}

func (r *Repository) Country(name string) (string, error) {
	var nationalityData NationalizeData

	nationality, err := httpGet(r.nationalityUrl, name, r.log)
	if err != nil {
		r.log.Errorln(err)
		return "", err
	}

	err = json.Unmarshal(nationality, &nationalityData)
	if err != nil {
		r.log.Errorln(err)
		return "", err
	}

	country := chooseCountryWithMaxProbability(nationalityData.Country)
	return country, nil
}

func (r *Repository) Gender(name string) (string, error) {
	var genderData GenderData

	gender, err := httpGet(r.genderUrl, name, r.log)
	if err != nil {
		r.log.Errorln(err)
		return "", err
	}

	err = json.Unmarshal(gender, &genderData)
	if err != nil {
		r.log.Errorln(err)
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

func httpGet(link, name string, log *logrus.Logger) ([]byte, error) {
	resp, err := http.Get(link + name)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	return body, nil
}
