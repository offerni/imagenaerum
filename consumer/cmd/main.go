package main

import (
	"github.com/offerni/imagenaerum/consumer/rabbitmq"
	"github.com/offerni/imagenaerum/consumer/rest"
	"github.com/offerni/imagenaerum/worker/utils"
)

func main() {
	utils.EnsureDirectories()

	rmqSvc := rabbitmq.Start()
	defer rmqSvc.Close()

	if err := rmqSvc.ConsumeConvertedFile("test"); err != nil {
		panic(err)
	}

	rest.InitializeServer(rest.ServerDependecies{
		RabbitMQSvc: *rmqSvc,
	})
}
