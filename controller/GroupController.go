package controller

import (
	"github.com/go-chi/chi/v5"
	"github.com/singhpranshu/splitbill/service/handler"
)

type GroupController struct {
	r       *chi.Mux
	handler *handler.Handler
}

func NewGroupController(r *chi.Mux, h *handler.Handler) *GroupController {
	return &GroupController{
		r:       r,
		handler: h,
	}
}

func (groupController *GroupController) RegisterRoutes() {
	groupController.r.Route("/group", func(r chi.Router) {
		r.Get("/{id}/", groupController.handler.GetGroup)
		r.Post("/", groupController.handler.CreateGroup)
		r.Post("/member", groupController.handler.Addmember)
	})
}
