package main

import (
	"log"
	"runners-mongodb/config"
	"runners-mongodb/server"
)

func main() {
	log.Println("Starting Runers App")

	log.Println("Initializig configuration")
	config := config.InitConfig("runners")

	log.Println("Initializig database")
	client := server.InitDatabase(config)

	log.Println("Initializig HTTP sever")
	httpServer := server.InitHttpServer(config, client)

	httpServer.Start()
}
