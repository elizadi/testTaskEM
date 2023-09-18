package user

import (
	"effective_mobile/internal/domain"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
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
	log *logrus.Logger
}

func New(db *gorm.DB, log *logrus.Logger) (domain.UserRepository, error) {
	if !db.Migrator().HasTable(&User{}) {
		err := db.Migrator().AutoMigrate(&User{})
		if err != nil {
			fmt.Println(err)
		}
	}
	return &Repository{
		db: db,
		log: log,
	}, nil
}

func (r *Repository) GetUsers() ([]domain.User, error) {
	var users []User

	err := r.db.Model(&User{}).Find(&users).Error
	if err != nil {
		r.log.Errorln(err)
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
	if name == "" || surname == "" {
		r.log.Errorln("empty parameter")
		return domain.User{}, errors.New("empty parameter")
	}
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
		r.log.Errorln(tx.Error)
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
		r.log.Errorln(tx.Error)
		return tx.Error
	}
	return nil
}

func (r *Repository) UpdateUser(id uint64, user domain.User) (domain.User, error) {
	checkUser := User{}
	tx := r.db.First(&checkUser, id)
	if tx.Error != nil {
		r.log.Errorln(tx.Error)
		return domain.User{}, tx.Error
	}
	user = checkNotNil(user, checkUser)
	savedUser := User{
		ID:         id,
		Name:       user.Name,
		Surname:    user.Surname,
		Patronymic: user.Patronymic,
		Age:        user.Age,
		Gender:     user.Gender,
		Country:    user.Country,
	}
	fmt.Print(savedUser)
	tx = r.db.Updates(&savedUser)
	if tx.Error != nil {
		r.log.Errorln(tx.Error)
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

func (r *Repository) GetUsersWithPagination(req domain.GetUsersReq) (domain.GetUsersResponse, error) {
	var users []User
	db := r.db.
		Model(&User{})
	if req.ID != 0 {
		db.Where("id = ?", req.ID)
	}
	if req.Name != "" {
		db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Surname != "" {
		db.Where("surname LIKE ?", "%"+req.Surname+"%")
	}
	if req.Patronymic != "" {
		db.Where("patronymic LIKE ?", "%"+req.Patronymic+"%")
	}
	if req.Age != 0 {
		db.Where("age = ?", req.Age)
	}
	if req.Gender != "" {
		db.Where("gender LIKE ?", "%"+req.Gender+"%")
	}
	if req.Country != "" {
		db.Where("country LIKE ?", "%"+req.Country+"%")
	}
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		r.log.Errorln(err)
		return domain.GetUsersResponse{}, err
	}
	err = db.
		Limit(req.Pag.Limit()).
		Offset(req.Pag.Offset()).
		Find(&users).
		Error
	if err != nil {
		r.log.Errorln(err)
		return domain.GetUsersResponse{}, err
	}

	pageCount := float64(total / int64(req.Pag.Limit()))
	odd := total % int64(req.Pag.Limit())
	if odd > 0 {
		pageCount++
	}
	response := domain.GetUsersResponse{
		Users: toDomainUsers(users),
		RespInfo: domain.RespInfo{
			Total:     total,
			PageCount: int64(pageCount),
		},
	}
	return response, nil
}

func checkNotNil(user domain.User, checkUser User) domain.User {
	if user.Name == "" {
		user.Name = checkUser.Name
	}
	if user.Surname == "" {
		user.Surname = checkUser.Surname
	}
	if user.Patronymic == "" {
		user.Patronymic = checkUser.Patronymic	
	}
	if user.Age == 0 {
		user.Age = checkUser.Age
	} 
	if user.Gender == "" {
		user.Gender = checkUser.Gender
	} 
	if user.Country == "" {
		user.Country = checkUser.Country
	} 
	return user
}