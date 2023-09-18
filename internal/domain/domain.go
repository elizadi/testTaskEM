package domain

type Server interface {
	Run() error
}

type UseCase interface {
	GetUsers() ([]User, error)
	CreateUser(name, surname, patronymic string) (User, error)
	DeleteUser(id uint64) error
	UpdateUser(id uint64, user User) (User, error)
	GetUsersWithPagination(req GetUsersReq) (GetUsersResponse, error)
}

type UserRepository interface {
	GetUsers() ([]User, error)
	CreateUser(name, surname, patronymic, gender, countryName string, age uint8) (User, error)
	DeleteUser(id uint64) error
	UpdateUser(id uint64, user User) (User, error)
	GetUsersWithPagination(req GetUsersReq) (GetUsersResponse, error)
}

type EnrichmentRepository interface {
	Age(name string) (uint8, error)
	Gender(name string) (string, error)
	Country(name string) (string, error)
}

type KafkaRepository interface {
	Produce(message string) error
	Consume() (FIOUser, error)
}
