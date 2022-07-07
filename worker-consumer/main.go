package main

import (
	"log"
	"os"
	"time"

	"github.com/ChrisMarSilva/cms-url-shortener-worker-consumer/databases"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

// go mod init github.com/ChrisMarSilva/cms-url-shortener-worker-consumer
// go get -u github.com/joho/godotenv
// go get -u github.com/asaskevich/govalidator
// go get -u github.com/go-redis/redis/v8
// go get -u github.com/google/uuid
// go get -u github.com/streadway/amqp
// go mod tidy

// go run main.go

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	// log.Println("amqp: ", os.Getenv("RABBIT_MQ_ADDR"))
	// log.Println("queue: ", os.Getenv("RABBIT_MQ_QUEUE"))

	conn, err := amqp.Dial(os.Getenv("RABBIT_MQ_ADDR"))
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel", err)
	}
	defer ch.Close()

	messages, err := ch.Consume(os.Getenv("RABBIT_MQ_QUEUE"), "", false, false, false, false, nil)
	if err != nil {
		log.Println("ch.Consume", err)
	}

	forever := make(chan bool)

	go func() {

		r := databases.CreateClient(0)
		defer r.Close()

		for message := range messages {

			log.Printf(" > Received message: %s\n", message.Body)

			// pegar a url, ip, id
			url := message.Body
			id := uuid.New().String()[:6]
			var expiry time.Duration = 24

			val, _ := r.Get(databases.Ctx, id).Result()
			if val == "" {
				err = r.Set(databases.Ctx, id, url, expiry*3600*time.Second).Err()
				if err != nil {
					log.Println("Unable to connext to server")
				}
			}

			message.Ack(true)

		}
	}()

	<-forever

}
