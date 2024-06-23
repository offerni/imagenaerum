package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/offerni/imagenaerum/worker/img"
	"github.com/offerni/imagenaerum/worker/rabbitmq"
	"github.com/offerni/imagenaerum/worker/utils"
)

func main() {
	utils.EnsureDirectories()

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}

	rabbitmqUrl := os.Getenv("RABBITMQ_URL")

	rabbitMqSvc := rabbitmq.Start(rabbitmqUrl)

	imgSvc, err := img.NewService(img.NewServiceOpts{
		RabbitmqSvc: rabbitMqSvc,
	})
	if err != nil {
		log.Println(err)
	}

	err = imgSvc.Process()
	if err != nil {
		log.Println(err)
	}

}
