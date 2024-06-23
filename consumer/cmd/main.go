package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/offerni/imagenaerum/consumer/img"
	"github.com/offerni/imagenaerum/consumer/rabbitmq"
	"github.com/offerni/imagenaerum/consumer/rest"
	"github.com/offerni/imagenaerum/worker/utils"
)

func main() {
	utils.EnsureDirectories()

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}
	rabbitmqUrl := os.Getenv("RABBITMQ_URL")

	rmqSvc := rabbitmq.Start(rabbitmqUrl)

	// Dumper so we can see the messages coming back from the worker
	go func() {
		defer rmqSvc.Close()
		if err := rmqSvc.ConsumeProcessedFile(rabbitmq.ConsumeProcessedFileOpts{
			QueueName:    "processed_files",
			ExchangeName: "file_exchange",
			RoutingKey:   "processed",
		}); err != nil {
			panic(err)
		}
	}()

	imgSvc, err := img.NewService(img.NewServiceOpts{
		RabbitMQSvc: rmqSvc,
	})
	if err != nil {
		panic(err)
	}

	go func() {
		rest.InitializeServer(rest.ServerDependecies{
			RabbitMQSvc: *rmqSvc,
			ImgSvc:      *imgSvc,
		})
	}()

}
