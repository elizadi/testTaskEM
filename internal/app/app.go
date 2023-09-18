package app

import (
	"effective_mobile/internal/delivery/graphql"
	delivery "effective_mobile/internal/delivery/http"
	"effective_mobile/internal/repository/enrichment"
	"effective_mobile/internal/repository/kafka"
	"effective_mobile/internal/repository/user"
	"effective_mobile/internal/usecase"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func Run(ageUrl, genderUrl, countryUrl, dsn, host, port string, log *logrus.Logger) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "effective_mobile_",
		},
	})
	if err != nil {
		log.WithError(err).Errorln("Cannot open gorm connect")
		return
	}

	userRepo, err := user.New(db, log)
	if err != nil {
		log.WithError(err).Errorln("Cannot create user repository")
	}

	enrichmentRepo, err := enrichment.New(ageUrl, genderUrl, countryUrl, log)
	if err != nil {
		log.WithError(err).Errorln("Cannot create enrichment repository")
	}

	kafkaRepo, err := kafka.New()
	if err != nil {
		log.WithError(err).Errorln("Cannot create enrichment repository")
	}

	useCase, err := usecase.New(userRepo, enrichmentRepo, kafkaRepo, log)
	if err != nil {
		log.WithError(err).Errorln("Cannot create usecase")
	}

	router := gin.Default()
	graphql.Register(useCase, router)
	delivery.Register(useCase, router)
	intPort, err := strconv.Atoi(port)
	if err != nil {
		log.WithError(err).Errorln("Cannot parse port")
	}
	err = router.Run(fmt.Sprintf("%s:%d", host, intPort))
	if err != nil {
		log.WithError(err).Errorln("server fall")
	}
}
