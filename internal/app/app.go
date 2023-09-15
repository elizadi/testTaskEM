package app

import (
	"effective_mobile/internal/delivery/graphql"
	delivery "effective_mobile/internal/delivery/http"
	"effective_mobile/internal/repository/enrichment"
	"effective_mobile/internal/repository/user"
	"effective_mobile/internal/usecase"
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func Run(ageUrl, genderUrl, countryUrl, dsn, host, port string) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "effective_mobile_",
		},
	})
	if err != nil {
		log.Fatalf("Cannot open gorm connect")
	}

	userRepo, err := user.New(db)
	if err != nil {
		log.Fatalf("Cannot create user repository")
	}

	enrichmentRepo, err := enrichment.New(ageUrl, genderUrl, countryUrl)
	if err != nil {
		log.Fatalf("Cannot create enrichment repository")
	}

	useCase, err := usecase.New(userRepo, enrichmentRepo)
	if err != nil {
		log.Fatalf("Cannot create usecase")
	}

	router := gin.Default()
	graphql.Register(useCase, router)
	delivery.Register(useCase, router)
	intPort, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalf("Cannot parse port")
	}
	err = router.Run(fmt.Sprintf("%s:%d", host, intPort))
	if err != nil {
		log.Fatalf("server fall")
	}
}
