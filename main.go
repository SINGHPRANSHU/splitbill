package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/singhpranshu/splitbill/common"
	"github.com/singhpranshu/splitbill/controller"
	"github.com/singhpranshu/splitbill/kafka" // Import Kafka package
	db "github.com/singhpranshu/splitbill/repository"
	"github.com/singhpranshu/splitbill/service/handler"
)

func main() {
	config := common.LoadConfig()
	db := db.NewDB(config.DatabaseURL)
	handler := handler.NewHandler(db)
	log.Println("Database connection established")
	go func() {
		log.Println("Starting Kafka consumer...")

		err := kafka.ConsumeMessages(config.Broker, config.GroupID, config.Topic)
		if err != nil {
			log.Fatalf("Kafka consumer error: %v", err)
		}
		log.Println("Kafka consumer is listening...")

	}()
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	controller.NewUserController(r, handler).RegisterRoutes()
	controller.NewGroupController(r, handler).RegisterRoutes()

	http.ListenAndServe(":"+config.Port, r)
}
