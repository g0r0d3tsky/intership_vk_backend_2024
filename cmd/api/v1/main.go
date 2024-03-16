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

	_ "cinema_service/docs"

	"github.com/joho/godotenv"
	"github.com/swaggo/http-swagger"
)

// @title Cinema Service API Documentation
// @version 1.0
// @description This is the API documentation for the Cinema Service.
// @host localhost:8080
// @BasePath /api/v1/
// @securityDefinitions.apiKey apiKey
// @in header
// @name Authorization
// @securityDefinitions.jwt apiKey
// @tokenUrl http://localhost:8080/api/v1/signIn
// @security jwt
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
	mux.Handle("/swagger/", httpSwagger.WrapHandler)
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
