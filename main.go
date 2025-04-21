package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/singhpranshu/splitbill/common"
	"github.com/singhpranshu/splitbill/controller"
	db "github.com/singhpranshu/splitbill/repository"
	"github.com/singhpranshu/splitbill/service/handler"
)

func main() {
	config := common.LoadConfig()
	db := db.NewDB(config.DatabaseURL)
	handler := handler.NewHandler(db)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	controller.NewUserController(r, handler).RegisterRoutes()
	controller.NewGroupController(r, handler).RegisterRoutes()

	http.ListenAndServe(":" + config.Port, r)
}
