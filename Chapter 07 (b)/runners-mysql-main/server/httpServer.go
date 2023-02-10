package server

import (
	"database/sql"
	"log"
	"runners-mysql/controllers"
	"runners-mysql/repositories"
	"runners-mysql/services"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type HttpServer struct {
	config            *viper.Viper
	router            *gin.Engine
	runnersController *controllers.RunnersController
	resultsController *controllers.ResultsController
}

func InitHttpServer(config *viper.Viper, dbHandler *sql.DB) HttpServer {
	runnersRepository := repositories.NewRunnersRepository(dbHandler)
	resultRepository := repositories.NewResultsRepository(dbHandler)
	runnersService := services.NewRunnersService(runnersRepository, resultRepository)
	resultsService := services.NewResultsService(resultRepository, runnersRepository)
	runnersController := controllers.NewRunnersController(runnersService)
	resultsController := controllers.NewResultsController(resultsService)

	router := gin.Default()

	router.POST("/runner", runnersController.CreateRunner)
	router.PUT("/runner", runnersController.UpdateRunner)
	router.DELETE("/runner/:id", runnersController.DeleteRunner)
	router.GET("/runner/:id", runnersController.GetRunner)
	router.GET("/runner", runnersController.GetRunnersBatch)

	router.POST("/result", resultsController.CreateResult)
	router.DELETE("/result/:id", resultsController.DeleteResult)

	return HttpServer{
		config:            config,
		router:            router,
		runnersController: runnersController,
		resultsController: resultsController,
	}
}

func (hs HttpServer) Start() {
	err := hs.router.Run(hs.config.GetString("http.server_address"))
	if err != nil {
		log.Fatalf("Error while starting HTTP server: %v", err)
	}
}
