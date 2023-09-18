package usecase

import (
	"effective_mobile/internal/domain"
	"time"

	"github.com/sirupsen/logrus"
)

func New(userRepo domain.UserRepository,
	enrichmentRepo domain.EnrichmentRepository,
	kafkaRepo domain.KafkaRepository,
	log *logrus.Logger) (domain.UseCase, error) {
	uc := useCase{
		userRepo:       userRepo,
		enrichmentRepo: enrichmentRepo,
		kafkaRepo:      kafkaRepo,
		log:            log,
	}

	go uc.backgroundKafka()
	return &uc, nil
}

type useCase struct {
	userRepo       domain.UserRepository
	enrichmentRepo domain.EnrichmentRepository
	kafkaRepo      domain.KafkaRepository
	log            *logrus.Logger
}

func (u *useCase) GetUsers() ([]domain.User, error) {
	users, err := u.userRepo.GetUsers()
	if err != nil {
		u.log.Errorln(err)
		return nil, err
	}
	return users, nil
}

func (u *useCase) CreateUser(name, surname, patronymic string) (domain.User, error) {
	age, err := u.enrichmentRepo.Age(name)
	if err != nil {
		u.log.Errorln(err)
		return domain.User{}, err
	}
	gender, err := u.enrichmentRepo.Gender(name)
	if err != nil {
		u.log.Errorln(err)
		return domain.User{}, err
	}

	country, err := u.enrichmentRepo.Country(name)
	if err != nil {
		u.log.Errorln(err)
		return domain.User{}, err
	}

	user, err := u.userRepo.CreateUser(name, surname, patronymic, gender, country, age)
	if err != nil {
		u.log.Errorln(err)
		return domain.User{}, err
	}
	return user, nil
}

func (u *useCase) DeleteUser(id uint64) error {
	err := u.userRepo.DeleteUser(id)
	if err != nil {
		u.log.Errorln(err)
		return err
	}
	return nil
}

func (u *useCase) UpdateUser(id uint64, user domain.User) (domain.User, error) {
	savedUser, err := u.userRepo.UpdateUser(id, user)
	if err != nil {
		u.log.Errorln(err)
		return domain.User{}, err
	}
	return savedUser, nil
}

func (u *useCase) GetUsersWithPagination(req domain.GetUsersReq) (domain.GetUsersResponse, error) {
	users, err := u.userRepo.GetUsersWithPagination(req)
	if err != nil {
		u.log.Errorln(err)
		return domain.GetUsersResponse{}, err
	}
	return users, nil
}

func (u *useCase) backgroundKafka() {
	for {
		u.log.Traceln("new iterate")
		user, err := u.kafkaRepo.Consume()
		if err != nil {
			if err == domain.ErrEmptyMessage {
				u.log.Traceln("empty")
				continue
			}
			u.log.WithError(err).Errorln("cannot consume from kafka error")
			err = u.kafkaRepo.Produce(err.Error())
			if err != nil {
				u.log.WithError(err).Errorln("cannot produce to kafka error")
			}
			continue
		}
		errs := user.Validate()
		if len(errs) > 0 {
			err = u.kafkaRepo.Produce(errs.String())
			if err != nil {
				u.log.WithError(err).Errorln("cannot produce to kafka user not validated")
			}
			continue
		}
		_, err = u.CreateUser(user.Name, user.Surname, user.Patronymic)
		if err != nil {
			u.log.WithError(err).Errorln("cannot save user")
			continue
		}
		time.Sleep(time.Second)
	}
}
