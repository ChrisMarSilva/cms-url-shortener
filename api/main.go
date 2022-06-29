package main

import (
	"log"
	"os"

	"github.com/ChrisMarSilva/cms-url-shortener/route"
	"github.com/joho/godotenv"
)

// go mod init github.com/ChrisMarSilva/cms-url-shortener
// go get -u github.com/gofiber/fiber/v2
// go get -u github.com/gofiber/utils
// go get -u github.com/gofiber/fiber/middleware
// go get -u github.com/joho/godotenv
// go get -u github.com/asaskevich/govalidator
// go get -u github.com/go-redis/redis/v8
// go get -u github.com/google/uuid
// go get -u github.com/streadway/amqp
// go mod tidy

// http://localhost:3000/efa14d  // erro
// http://127.0.0.1:3000/efa14d  // ok

// go run main.go

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := route.NewRoutes()
	log.Fatal(app.Listen(os.Getenv("API_PORT")))

}
