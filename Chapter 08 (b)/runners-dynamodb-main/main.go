package main

import (
	"log"
	"runners-dynamodb/config"
	"runners-dynamodb/server"
)

func main() {
	log.Println("Starting Runers App")

	log.Println("Initializig configuration")
	config := config.InitConfig("runners")

	log.Println("Initializig database")
	dynamoDB := server.InitDatabase(config)

	log.Println("Initializig HTTP sever")
	httpServer := server.InitHttpServer(config, dynamoDB)

	httpServer.Start()
}
