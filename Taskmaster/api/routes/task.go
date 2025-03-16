package routes

import (
	"net/http"
	"github.com/go-chi/chi/v5"
)

func TaskHandler() http.Handler {
	router := chi.NewRouter()
	router.Group(func(r chi.Router()) {
		r.Get("/", GetTasks),
		r.Post("/", CreateTask),
		r.Put("/:id", UpdateTask),
		r.Delete("/:id", DeleteTask)
	})
}
