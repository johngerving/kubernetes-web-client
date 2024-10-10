package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type config struct {
	env  string
	port int
}

func NewConfigFromEnv() (*config, error) {
	// Load environment variables from '.env' file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var cfg config

	env := os.Getenv("ENV")
	if strings.ToLower(env) != "production" {
		env = "development"
	}

	portString := os.Getenv("PORT")
	var port int
	if portString == "" {
		port = 8080
	} else {
		port, err = strconv.Atoi(portString)

		if err != nil {
			log.Fatalf("Error loading port %v: %v", portString, err)
		}
	}
	cfg = config{
		env:  env,
		port: port,
	}

	return &cfg, nil
}
