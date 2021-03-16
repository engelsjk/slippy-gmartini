package main

import (
	"log"

	slippy "github.com/engelsjk/slippy-gmartini"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	slippy.Run()
}
