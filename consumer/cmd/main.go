package main

import (
	"github.com/offerni/imagenaerum/consumer/rabbitmq"
	"github.com/offerni/imagenaerum/consumer/rest"
	"github.com/offerni/imagenaerum/worker/utils"
)

func main() {
	utils.EnsureDirectories()

	rmqSvc := rabbitmq.Start()
	go func() {
		defer rmqSvc.Close()
		if err := rmqSvc.ConsumeConvertedFile(rabbitmq.ConsumeConvertedFileOpts{
			QueueName:    "converted_files",
			ExchangeName: "files_exchange",
			RoutingKey:   "converted",
		}); err != nil {
			panic(err)
		}
	}()

	rest.InitializeServer(rest.ServerDependecies{
		RabbitMQSvc: *rmqSvc,
	})
}
