package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	app "github.com/fajrulaulia/ampqgotest"
)

func main() {
	var err error
	port := "3000"
	if len(os.Getenv("PORT")) > 0 {
		port = os.Getenv("PORT")
	}

	r := mux.NewRouter()

	var c app.Config
	c.Route = r
	c.Ampq, err = app.InitalizeRabbitMQ()
	if err != nil {
		log.Printf("Error: Rabbit MQ failed to connect to server with error message %s", err.Error())
		return
	}
	log.Println("[V] RabbitMQ Connected")

	c.Es, err = app.IntializeElasticSearch()
	if err != nil {
		log.Printf("Error: Rabbit MQ failed to connect to server with error message %s", err.Error())
		return
	}
	log.Println("[V] Elasticsearch Connected")

	c.Apply()

	log.Println("[V] Service running on port", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}
