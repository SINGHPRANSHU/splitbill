package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/singhpranshu/splitbill/common"
	"github.com/singhpranshu/splitbill/controller"
	db "github.com/singhpranshu/splitbill/repository"
	"github.com/singhpranshu/splitbill/service/handler"
	jwt "github.com/singhpranshu/splitbill/service/middleware"
)

func main() {
	config := common.LoadConfig()
	db := db.NewDB(config.DatabaseURL)
	handler := handler.NewHandler(db)
	jwt.InitJwt(config)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	controller.NewUserController(r, handler).RegisterRoutes()
	controller.NewGroupController(r, handler).RegisterRoutes()
	controller.NewSplitController(r, handler).RegisterRoutes()

	http.ListenAndServe(":"+config.Port, r)
}
