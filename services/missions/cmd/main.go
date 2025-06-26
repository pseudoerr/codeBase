package main

import (
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"

	_ "github.com/lib/pq"
	"github.com/pseudoerr/mission-service/config"
	_ "github.com/pseudoerr/mission-service/docs"
	"github.com/pseudoerr/mission-service/internal/handler"
	"github.com/pseudoerr/mission-service/repository"
	"github.com/pseudoerr/mission-service/service"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	config.LoadEnv()
	db, err := sql.Open("postgres", config.GetDatabaseURL())
	if err != nil {
		logger.Warn("failed to connect to db", "postgres", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		logger.Warn("failed to ping db", "ping", err)
	}
	logger.Info("connected to psql!")

	repo := repository.NewPostgresRepository(db)
	svc := &service.MissionService{
		Store:  repo,
		Logger: logger,
	}

	newHandler := &handler.Handler{Service: svc}
	router := handler.NewRouter(newHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	go func() {
		slog.Info("pprof available at :6060/debug/pprof")
		log.Println(http.ListenAndServe(":6060", nil))
	}()

	logger.Info("starting http server on port 8080")
	http.ListenAndServe(":"+port, router)

}
