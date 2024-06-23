module github.com/offerni/imagenaerum/consumer

go 1.22.3

require (
	github.com/go-chi/chi v1.5.5
	github.com/joho/godotenv v1.5.1
	github.com/offerni/imagenaerum/worker v0.0.0
)

replace github.com/offerni/imagenaerum/worker => ../worker
