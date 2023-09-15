package usecase

import (
	"effective_mobile/internal/domain"
)

func New(userRepo domain.UserRepository,
	enrichmentRepo domain.EnrichmentRepository) (domain.UseCase, error) {
	return &useCase{
		userRepo:       userRepo,
		enrichmentRepo: enrichmentRepo,
	}, nil
}

type useCase struct {
	userRepo       domain.UserRepository
	enrichmentRepo domain.EnrichmentRepository
}

func (u *useCase) GetUsers() ([]domain.User, error) {
	users, err := u.userRepo.GetUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *useCase) CreateUser(name, surname, patronymic string) (domain.User, error) {
	age, err := u.enrichmentRepo.Age(name)
	if err != nil {
		return domain.User{}, err
	}
	gender, err := u.enrichmentRepo.Gender(name)
	if err != nil {
		return domain.User{}, err
	}

	country, err := u.enrichmentRepo.Country(name)
	if err != nil {
		return domain.User{}, err
	}

	user, err := u.userRepo.CreateUser(name, surname, patronymic, gender, country, age)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (u *useCase) DeleteUser(id uint64) error {
	err := u.userRepo.DeleteUser(id)
	if err != nil {
		return err
	}
	return nil
}

func (u *useCase) UpdateUser(id uint64, user domain.User) (domain.User, error) {
	savedUser, err := u.userRepo.UpdateUser(id, user)
	if err != nil {
		return domain.User{}, err
	}
	return savedUser, nil
}

func (u *useCase) GetUsersWithPagination(page, perPage uint) ([]domain.User, error) {
	users, err := u.userRepo.GetUsersWithPagination(page, perPage)
	if err != nil{
		return nil, err
	}
	return users, nil
}