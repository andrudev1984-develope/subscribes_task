package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"strconv"
	"subscribes/internal/adapters/database"
	"subscribes/internal/adapters/protocol"
	"subscribes/internal/config"
	"subscribes/internal/usecase"
	"subscribes/openapi/subscribe"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	mainCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	appConfig := config.NewConfig()

	pool, err := pgxpool.New(ctx, fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s",
		appConfig.Db.User, appConfig.Db.Password, appConfig.Db.Host, appConfig.Db.Port, appConfig.Db.Name))

	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()

	dbAdapter := database.NewPgrep(pool)
	httpAdapter := protocol.NewRouter()
	usecaseAdapter := usecase.NewUseCase(dbAdapter)

	var h = subscribe.NewStrictHandler(usecaseAdapter, []subscribe.StrictMiddlewareFunc{})

	subscribe.HandlerFromMux(h, httpAdapter)

	var server = &http.Server{
		Addr:    ":" + strconv.Itoa(appConfig.Server.Port),
		Handler: httpAdapter,
	}

	go func() {
		log.Fatal(server.ListenAndServe())
	}()

	<-mainCtx.Done()

	log.Println("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal(err)
	}
}
