package user

import (
	"effective_mobile/internal/domain"
	"fmt"

	"gorm.io/gorm"
)

type User struct {
	ID         uint64 `gorm:"column:id;primaryKey;autoIncrement:true"`
	Name       string `gorm:"column:name"`
	Surname    string `gorm:"column:surname"`
	Patronymic string `gorm:"column:patronymic"`
	Age        uint8  `gorm:"column:age"`
	Gender     string `gorm:"column:gender"`
	Country    string `gorm:"column:country"`
}

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) (domain.UserRepository, error) {
	if !db.Migrator().HasTable(&User{}) {
		err := db.Migrator().AutoMigrate(&User{})
		if err != nil {
			fmt.Println(err)
		}
	}
	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) GetUsers() ([]domain.User, error) {
	var users []User

	err := r.db.Model(&User{}).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return toDomainUsers(users), nil
}

func toDomainUsers(users []User) []domain.User {
	res := make([]domain.User, 0, len(users))
	for _, user := range users {
		res = append(res, toDomain(user))
	}
	return res
}

func (r *Repository) CreateUser(name, surname, patronymic, gender, countryName string, age uint8) (domain.User, error) {
	user := User{
		Name:       name,
		Surname:    surname,
		Patronymic: patronymic,
		Age:        age,
		Gender:     gender,
		Country:    countryName,
	}
	tx := r.db.Create(&user)
	if tx.Error != nil {
		fmt.Printf("failed to add %v\n", user)
		return domain.User{}, tx.Error
	}
	fmt.Printf("%v successfully added\n", user)

	return domain.User{
		ID:         user.ID,
		Name:       user.Name,
		Surname:    user.Surname,
		Patronymic: user.Patronymic,
		Age:        user.Age,
		Gender:     user.Gender,
		Country:    user.Country,
	}, nil
}

func (r *Repository) DeleteUser(id uint64) error {
	user := User{
		ID: id,
	}
	tx := r.db.Delete(user)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (r *Repository) UpdateUser(id uint64, user domain.User) (domain.User, error) {
	tx := r.db.First(&user, id)
	if tx.Error != nil {
		return domain.User{}, tx.Error
	}
	savedUser := User{
		ID:         id,
		Name:       user.Name,
		Surname:    user.Surname,
		Patronymic: user.Patronymic,
		Age:        user.Age,
		Gender:     user.Gender,
		Country:    user.Country,
	}
	tx = r.db.Save(&savedUser)
	if tx.Error != nil {
		return domain.User{}, tx.Error
	}
	return toDomain(savedUser), nil
}

func toDomain(savedUser User) domain.User {
	return domain.User{
		ID:         savedUser.ID,
		Name:       savedUser.Name,
		Surname:    savedUser.Surname,
		Patronymic: savedUser.Patronymic,
		Age:        savedUser.Age,
		Gender:     savedUser.Gender,
		Country:    savedUser.Country,
	}
}

func (r *Repository) GetUsersWithPagination(page, perPage uint) ([]domain.User, error) {
	var users []User
	offset := (page - 1) * perPage
	err := r.db.Model(&User{}).Limit(int(perPage)).Offset(int(offset)).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return toDomainUsers(users), nil
}
