package main

import (
	"cinema_service/config"
	"cinema_service/internal/api/handlers"
	"cinema_service/internal/api/middleware"
	"cinema_service/internal/repository"
	"cinema_service/internal/usecase"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	c, err := config.Read()

	if err != nil {
		log.Println("failed to read config:", err.Error())
		return
	}
	dbPool, err := repository.Connect(c)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer func() {
		if dbPool != nil {
			dbPool.Close()
		}
	}()

	storageActor := repository.NewStorageActor(dbPool)
	storageMovie := repository.NewStorageMovie(dbPool)
	storageUser := repository.NewUserStorage(dbPool)

	serviceActor := usecase.NewActorsService(&storageActor)
	serviceMovie := usecase.NewMovieService(&storageMovie)
	serviceUser := usecase.NewUserService(&storageUser)

	handlerActor := handlers.NewActorHandler(serviceActor)
	handlerMovie := handlers.NewMovieHandler(serviceMovie)
	handlerUser := handlers.NewUserHandler(serviceUser)

	middlewareUser := middleware.NewUserMiddleware(serviceUser)
	//authentication := middlewareUser.Authenticate()

	mux := http.NewServeMux()

	mux = handlerActor.RegisterActor(mux, middlewareUser.Authenticate, middlewareUser.RequireAdmin)
	mux = handlerMovie.RegisterMovie(mux, middlewareUser.Authenticate, middlewareUser.RequireAdmin)
	mux = handlerUser.RegisterUser(mux)
	server := &http.Server{
		Addr:    net.JoinHostPort(c.Host, c.Port),
		Handler: mux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Starting server on port %v...\n", c.Port)
		err = server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-stop

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Server gracefully stopped")
}
