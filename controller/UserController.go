package controller

import (
	"github.com/go-chi/chi/v5"
	"github.com/singhpranshu/splitbill/service/handler"
)

type UserController struct {
	r       *chi.Mux
	handler *handler.Handler
}

func NewUserController(r *chi.Mux, h *handler.Handler) *UserController {
	return &UserController{
		r:       r,
		handler: h,
	}
}

func (userController *UserController) RegisterRoutes() {
	userController.r.Route("/user", func(r chi.Router) {
		r.Get("/{user_id}", userController.handler.GetUser)
		r.Post("/", userController.handler.CreateUser)
	})
}
