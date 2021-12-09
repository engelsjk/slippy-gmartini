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
	s, err := slippy.New()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("listening on port %s\n", s.Port())
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
