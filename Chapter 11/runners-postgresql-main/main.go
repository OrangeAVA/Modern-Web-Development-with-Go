package main

import (
	"log"
	"os"
	"runners-postgresql/config"
	"runners-postgresql/server"

	_ "github.com/lib/pq"
)

func main() {
	log.Println("Starting Runners App")

	log.Println("Initializing configuration")
	config := config.InitConfig(getConfigFileName())

	log.Println("Initializing database")
	dbHandler := server.InitDatabase(config)

	log.Println("Initializing Prometheus")
	go server.InitPrometheus()

	log.Println("Initializig HTTP sever")
	httpServer := server.InitHttpServer(config, dbHandler)

	httpServer.Start()
}

func getConfigFileName() string {
	env := os.Getenv("ENV")

	if env != "" {
		return "runners-" + env
	}

	return "runners"
}
