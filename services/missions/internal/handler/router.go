package handler

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pseudoerr/mission-service/internal/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(handler *Handler) http.Handler {
	r := mux.NewRouter()

	// no auth
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	authRoutes := r.PathPrefix("/").Subrouter()

	authRoutes.HandleFunc("/missions", handler.GetMissions).Methods("GET")
	authRoutes.HandleFunc("/missions/{id:[0-9]+}", handler.GetMissionByID).Methods("GET")
	authRoutes.HandleFunc("/missions", handler.CreateMission).Methods("POST")
	authRoutes.HandleFunc("/missions/{id:[0-9]+}", handler.UpdateMission).Methods("PUT")
	authRoutes.HandleFunc("/missions/{id:[0-9]+}", handler.DeleteMission).Methods("DELETE")
	authRoutes.HandleFunc("/profile", handler.GetProfile).Methods("GET")

	authRoutes.Use(middleware.AuthMiddleware)

	rl := middleware.NewRateLimiter(10, time.Minute)

	var handlerWithMiddleware http.Handler = r
	handlerWithMiddleware = rl.MiddleWare(handlerWithMiddleware)
	handlerWithMiddleware = middleware.LoggingMiddleware(handlerWithMiddleware)
	handlerWithMiddleware = middleware.RecoverMiddleware(handlerWithMiddleware)
	handlerWithMiddleware = middleware.CORSMiddleware(handlerWithMiddleware)

	return handlerWithMiddleware
}
