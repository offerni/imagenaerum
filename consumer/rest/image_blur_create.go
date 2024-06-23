package rest

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/offerni/imagenaerum/consumer/rabbitmq"
)

func (srv *Server) ImageBlurCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("hello!!")

	_ = r.MultipartForm.File["files"]
	sigma := r.FormValue("sigma")
	_, err := strconv.ParseFloat(sigma, 64)
	if err != nil {
		log.Println(err)
	}

	err = srv.RabbitMQSvc.Publish(rabbitmq.PublishOpts{
		Ch:           srv.RabbitMQSvc.Channel,
		QueueName:    "files_to_convert",
		ExchangeName: "file_exchange",
		RoutingKey:   "to_convert",
		Body:         []byte("Hello World"),
	})
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusNoContent)
}
