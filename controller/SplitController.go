package controller

import (
	"github.com/go-chi/chi/v5"
	"github.com/singhpranshu/splitbill/service/handler"
	jwt "github.com/singhpranshu/splitbill/service/middleware"
)

type SplitController struct {
	r       *chi.Mux
	handler *handler.Handler
}

func NewSplitController(r *chi.Mux, h *handler.Handler) *SplitController {
	return &SplitController{
		r:       r,
		handler: h,
	}
}

func (splitController *SplitController) RegisterRoutes() {
	splitController.r.Route("/split", func(r chi.Router) {
		r.Use(jwt.AuthenticateMiddlewareHandler)
		r.Get("/{id}/", splitController.handler.GetSplitData)
		r.Post("/", splitController.handler.AddExpense)
	})
}
