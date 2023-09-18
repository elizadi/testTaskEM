package graph

import (
	"effective_mobile/internal/domain"
	"effective_mobile/internal/graph/model"

)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	UseCase domain.UseCase
}

func toModelUsers(users []domain.User) []*model.User {
	res := make([]*model.User, 0, len(users))
	for _, user := range users {
		user := user
		
		age := int(user.Age)
		res = append(res, &model.User{
			ID:         int(user.ID),
			Name:       user.Name,
			Surname:    user.Surname,
			Patronymic: user.Patronymic,
			Age:        age,
			Gender:     user.Gender,
			Country:    user.Country,
		})
	}
	return res
}


func toDomainUser(user model.UserIn) domain.User {
	var domainUser domain.User
	domainUser.ID = uint64(user.ID)
	if user.Age != nil {
		domainUser.Age = uint8(*user.Age)
	}
	if user.Name != nil {
		domainUser.Name = *user.Name
	}
	if user.Surname != nil {
		domainUser.Surname = *user.Surname
	}
	if user.Patronymic != nil {
		domainUser.Patronymic = *user.Patronymic
	}
	if user.Gender != nil {
		domainUser.Gender = *user.Gender
	}
	if user.Country != nil {
		domainUser.Country = *user.Country
	}
	return domainUser
}