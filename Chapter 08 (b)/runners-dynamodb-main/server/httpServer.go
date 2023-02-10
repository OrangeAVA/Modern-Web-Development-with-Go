package server

import (
	"log"
	"runners-dynamodb/controllers"
	"runners-dynamodb/repositories"
	"runners-dynamodb/services"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type HttpServer struct {
	config            *viper.Viper
	router            *gin.Engine
	runnersController *controllers.RunnersController
	resultsController *controllers.ResultsController
}

func InitHttpServer(config *viper.Viper, dynamoDB *dynamodb.DynamoDB) HttpServer {
	runnersRepository := repositories.NewRunnersRepository(dynamoDB)
	resultRepository := repositories.NewResultsRepository(dynamoDB)
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
