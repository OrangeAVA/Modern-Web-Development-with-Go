package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"runners-postgresql/models"
	"runners-postgresql/repositories"
	"runners-postgresql/services"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetRunnersResponse(t *testing.T) {
	dbHandler, mock, _ := sqlmock.New()
	defer dbHandler.Close()

	columnsUsers := []string{"user_role"}
	mock.ExpectQuery("SELECT user_role").WillReturnRows(
		sqlmock.NewRows(columnsUsers).AddRow("runner"),
	)

	columns := []string{"id", "first_name", "last_name", "age", "is_active", "country", "personal_best", "season_best"}
	mock.ExpectQuery("SELECT *").WillReturnRows(
		sqlmock.NewRows(columns).
			AddRow("1", "John", "Smith", 30, true, "United States", "02:00:41", "02:13:13").
			AddRow("2", "Marijana", "Komatinovic", 30, true, "Serbia", "01:18:28", "01:18:28"))

	router := initTestRouter(dbHandler)
	request, _ := http.NewRequest("GET", "/runner", nil)
	request.Header.Set("token", "token")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode)

	var runers []*models.Runner
	json.Unmarshal(recorder.Body.Bytes(), &runers)

	assert.NotEmpty(t, runers)
	assert.Equal(t, 2, len(runers))
}

func initTestRouter(dbHandler *sql.DB) *gin.Engine {
	runnersRepository := repositories.NewRunnersRepository(dbHandler)
	usersRepository := repositories.NewUsersRepository(dbHandler)
	runnersService := services.NewRunnersService(runnersRepository, nil)
	usersServices := services.NewUsersService(usersRepository)
	runnersController := NewRunnersController(runnersService, usersServices)

	router := gin.Default()

	router.GET("/runner", runnersController.GetRunnersBatch)

	return router
}
