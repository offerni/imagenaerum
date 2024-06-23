package main

import (
	"github.com/offerni/imagenaerum/consumer/img"
	"github.com/offerni/imagenaerum/consumer/rabbitmq"
	"github.com/offerni/imagenaerum/consumer/rest"
	"github.com/offerni/imagenaerum/worker/utils"
)

func main() {
	utils.EnsureDirectories()

	rmqSvc := rabbitmq.Start()
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

	rest.InitializeServer(rest.ServerDependecies{
		RabbitMQSvc: *rmqSvc,
		ImgSvc:      *imgSvc,
	})
}
