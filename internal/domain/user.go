package domain

import "errors"

type User struct {
	ID         uint64
	Name       string
	Surname    string
	Patronymic string
	Age        uint8
	Gender     string
	Country    string
}

type FIOUser struct {
	Name       string
	Surname    string
	Patronymic string
}

func (u *FIOUser) Validate() Errs {
	res := make([]error, 0)
	if u.Name == "" {
		res = append(res, errors.New("name is not set"))
	}
	if u.Surname == "" {
		res = append(res, errors.New("surname is not set"))
	}
	return res
}
