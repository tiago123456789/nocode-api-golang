package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/tiago123456789/nocode-api-golang/internal/config"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := config.StartDB()

	if err != nil {
		log.Fatal(err)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("USER_PASSWORD")), 14)
	if err != nil {
		log.Fatal(err)
	}

	var id int64
	err = db.QueryRow(
		"INSERT INTO auth(name, email, password) VALUES ($1, $2, $3) RETURNING id",
		os.Getenv("USER_NAME"),
		os.Getenv("USER_EMAIL"),
		string(password),
	).Scan(&id)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("User created!!!")
}
