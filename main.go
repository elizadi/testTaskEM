package main

import (
	"effective_mobile/internal/app"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// type server struct {
// 	uc domain.UseCase
// }

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
	dsn := fmt.Sprintf("host=%s user=postgres password=123456789Lis port=5432 sslmode=disable", dbUrl)

	host:= os.Getenv("host")
	if dbUrl == "" {
		log.Fatalf("set host in env file")
	}

	port:= os.Getenv("port")
	if dbUrl == "" {
		log.Fatalf("set port in env file")
	}
	app.Run(ageUrl, genderUrl, countryUrl, dsn, host, port)
}
