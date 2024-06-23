module github.com/offerni/imagenaerum/consumer

go 1.22.3

require (
	github.com/go-chi/chi v1.5.5
	github.com/joho/godotenv v1.5.1
	github.com/offerni/imagenaerum/worker v0.0.0
	github.com/rabbitmq/amqp091-go v1.10.0
)

// recvisit this later
replace github.com/offerni/imagenaerum/worker => ../worker
