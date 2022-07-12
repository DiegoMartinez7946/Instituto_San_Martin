package main

import (
	"log"
	"net/http"

	"github.com/benjacifre10/san_martin_b/config"
	"github.com/benjacifre10/san_martin_b/routes"
	"github.com/benjacifre10/san_martin_b/utils"
	"github.com/rs/cors"
)

func main() {
	if config.CheckConnection() == 0 {
		log.Fatal("DB without connection")
		return
	}

	// find and set the port
	PORT := utils.GoDotEnvValue("PORT")
	if PORT == "" {
		PORT = "8081"
	}
	log.Println("Listening port: " + 	PORT)

	// set the router
	router := routes.HandlerRoutes()

	// set the permissions with cors
	handlerPermission := cors.AllowAll().Handler(router)

	// check the serve
	log.Fatal(http.ListenAndServe(":" + PORT, handlerPermission))

}
