package main

import (
	"effective_mobile/internal/app"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	ageUrl := os.Getenv("age_url")
	if ageUrl == "" {
		log.Fatalf("set age_url in env file")
	}
	genderUrl := os.Getenv("gender_url")
	if genderUrl == "" {
		log.Fatalf("set gender_url in env file")
	}
	countryUrl := os.Getenv("country_url")
	if countryUrl == "" {
		log.Fatalf("set country_url in env file")
	}

	dbUrl := os.Getenv("database")
	if dbUrl == "" {
		log.Fatalf("set database in env file")
	}
	user := os.Getenv("user")
	if dbUrl == "" {
		log.Fatalf("set user in env file")
	}
	password := os.Getenv("password")
	if dbUrl == "" {
		log.Fatalf("set password in env file")
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s port=5432 sslmode=disable", dbUrl, user, password)

	host := os.Getenv("host")
	if dbUrl == "" {
		log.Fatalf("set host in env file")
	}

	port := os.Getenv("port")
	if dbUrl == "" {
		log.Fatalf("set port in env file")
	}
	logger := logrus.New()
	logger.SetLevel(logrus.TraceLevel)
	app.Run(ageUrl, genderUrl, countryUrl, dsn, host, port, logger)
}
