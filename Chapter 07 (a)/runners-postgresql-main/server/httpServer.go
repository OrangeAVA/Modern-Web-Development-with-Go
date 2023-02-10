package server

import (
	"database/sql"
	"log"
	"runners-postgresql/controllers"
	"runners-postgresql/repositories"
	"runners-postgresql/services"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type HttpServer struct {
	config            *viper.Viper
	router            *gin.Engine
	runnersController *controllers.RunnersController
	resultsController *controllers.ResultsController
	usersController   *controllers.UsersController
}

func InitHttpServer(config *viper.Viper, dbHandler *sql.DB) HttpServer {
	runnersRepository := repositories.NewRunnersRepository(dbHandler)
	resultRepository := repositories.NewResultsRepository(dbHandler)
	usersRepository := repositories.NewUsersRepository(dbHandler)
	runnersService := services.NewRunnersService(runnersRepository, resultRepository)
	resultsService := services.NewResultsService(resultRepository, runnersRepository)
	usersService := services.NewUsersService(usersRepository)
	runnersController := controllers.NewRunnersController(runnersService, usersService)
	resultsController := controllers.NewResultsController(resultsService, usersService)
	usersController := controllers.NewUsersController(usersService)

	router := gin.Default()

	router.POST("/runner", runnersController.CreateRunner)
	router.PUT("/runner", runnersController.UpdateRunner)
	router.DELETE("/runner/:id", runnersController.DeleteRunner)
	router.GET("/runner/:id", runnersController.GetRunner)
	router.GET("/runner", runnersController.GetRunnersBatch)

	router.POST("/result", resultsController.CreateResult)
	router.DELETE("/result/:id", resultsController.DeleteResult)

	router.POST("/login", usersController.Login)
	router.POST("/logout", usersController.Logout)

	return HttpServer{
		config:            config,
		router:            router,
		runnersController: runnersController,
		resultsController: resultsController,
		usersController:   usersController,
	}
}

func (hs HttpServer) Start() {
	err := hs.router.Run(hs.config.GetString("http.server_address"))
	if err != nil {
		log.Fatalf("Error while starting HTTP server: %v", err)
	}
}
