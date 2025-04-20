package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/singhpranshu/splitbill/controller"
	db "github.com/singhpranshu/splitbill/repository"
	"github.com/singhpranshu/splitbill/service/handler"
)

func main() {
	db := db.NewDB()
	handler := handler.NewHandler(db)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	controller.NewUserController(r, handler).RegisterRoutes()

	http.ListenAndServe(":3000", r)
}
