package rest

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
	"github.com/offerni/imagenaerum/img"
)

const defaultport string = "8080"

type Server struct {
	HttpSrv    *http.Server
	ImgService img.Service
}

type NewServerOpts struct {
	HttpSrv    *http.Server
	ImgService img.Service
}

func InitializeServer() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultport
	}

	// multiplexer
	mux := chi.NewRouter()

	// middlewares
	mux.Use(middleware.Logger)
	mux.Use(contextMiddleware)

	// new Server
	imgService := img.NewService()
	srv := NewServer(NewServerOpts{
		ImgService: *imgService,
		HttpSrv: &http.Server{
			Addr:    fmt.Sprintf(":%s", port),
			Handler: mux,
		},
	})

	// routes
	initializeRoutes(mux, *srv)

	go func() {
		fmt.Printf("HTTP Server started on port %s\n", port)
		err = srv.HttpSrv.ListenAndServe()
		if err != nil && http.ErrServerClosed != err {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	fmt.Println("Gracefully Shutting Down Server...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := srv.HttpSrv.Shutdown(ctx); err != nil {
		log.Fatalf("Could not shutdown server %v\n", err)
	}
	fmt.Println("Server Stopped")
}

func contextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Second*30)
		defer cancel()

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func NewServer(opts NewServerOpts) *Server {
	return &Server{
		HttpSrv:    opts.HttpSrv,
		ImgService: opts.ImgService,
	}
}
