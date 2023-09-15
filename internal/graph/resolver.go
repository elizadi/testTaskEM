package graph

import (
	"effective_mobile/internal/domain"
	"effective_mobile/internal/graph/model"
	"strconv"
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
		id := strconv.Itoa(int(user.ID))
		age := int(user.Age)
		res = append(res, &model.User{
			ID:         &id,
			Name:       user.Name,
			Surname:    user.Surname,
			Patronymic: &user.Patronymic,
			Age:        &age,
			Gender:     &user.Gender,
			Country:    &user.Country,
		})
	}
	return res
}
