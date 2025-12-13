package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (s *Server) registerRoutes() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // Разрешить все источники
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	}))

	// HTML рендер
	r.Get("/", s.IndexHandler)
	r.Get("/users", s.UsersHandler)
	r.Get("/dataset", s.DatasetHandler)

	r.Get("/api/get_users", s.GetUsersHandler)
	r.Get("/api/get_dataset", s.GetDatasetHandler)
	r.Post("/api/connect", s.ConnectUserHandler)
	r.Post("/api/disconnect", s.DisconnectUserHandler)
	r.Post("/api/send_data", s.SendDataHandler)

	// статика
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	s.router = r
}
