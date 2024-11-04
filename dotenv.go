package main

import (
	//"log"
	"os"
	//"github.com/joho/godotenv"
)

func goDotEnvVariable(key string) string {
	// This is only for development
	/*
		err := godotenv.Load()
		if err != nil {
			log.Print("Error loading .env file")
			return os.Getenv(key)
		}
	*/

	return os.Getenv(key)
}
