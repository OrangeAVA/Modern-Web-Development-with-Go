package services

import (
	"net/http"
	"runners-dynamodb/models"
	"runners-dynamodb/repositories"
	"time"
)

type ResultsService struct {
	resultsRepository *repositories.ResultsRepository
	runnersRepository *repositories.RunnersRepository
}

func NewResultsService(resultsRepository *repositories.ResultsRepository,
	runnersRepository *repositories.RunnersRepository) *ResultsService {

	return &ResultsService{
		resultsRepository: resultsRepository,
		runnersRepository: runnersRepository,
	}
}

func (rs ResultsService) CreateResult(result *models.Result) (*models.Result, *models.ResponseError) {
	// Validation
	if result.RunnerID == "" {
		return nil, &models.ResponseError{
			Message: "Invalid runner ID",
			Status:  http.StatusBadRequest,
		}
	}

	if result.RaceResult == "" {
		return nil, &models.ResponseError{
			Message: "Invalid race result",
			Status:  http.StatusBadRequest,
		}
	}

	if result.Location == "" {
		return nil, &models.ResponseError{
			Message: "Invalid location",
			Status:  http.StatusBadRequest,
		}
	}

	if result.Position < 0 {
		return nil, &models.ResponseError{
			Message: "Invalid position",
			Status:  http.StatusBadRequest,
		}
	}

	currentYear := time.Now().Year()
	if result.Year < 0 || result.Year > currentYear {
		return nil, &models.ResponseError{
			Message: "Invalid year",
			Status:  http.StatusBadRequest,
		}
	}

	runner, responseErr := rs.runnersRepository.GetRunner(result.RunnerID)
	if responseErr != nil {
		return nil, responseErr
	}

	// update runners personal best
	if runner.PersonalBest == "" {
		runner.PersonalBest = result.RaceResult
	} else {
		if result.RaceResult < runner.PersonalBest {
			runner.PersonalBest = result.RaceResult
		}
	}

	// update runners seeason best
	if result.Year == currentYear {
		if runner.SeasonBest == "" {
			runner.SeasonBest = result.RaceResult
		} else {
			if result.RaceResult < runner.SeasonBest {
				runner.SeasonBest = result.RaceResult
			}
		}
	}

	response, responseErr := rs.resultsRepository.CreateResult(result, runner)
	if responseErr != nil {
		return nil, responseErr
	}

	return response, nil
}

func (rs ResultsService) DeleteResult(resultId string) *models.ResponseError {
	if resultId == "" {
		return &models.ResponseError{
			Message: "Invalid result ID",
			Status:  http.StatusBadRequest,
		}
	}

	result, responseErr := rs.resultsRepository.GetResult(resultId)
	if responseErr != nil {
		return responseErr
	}

	results, responseErr := rs.resultsRepository.GetAllRunnersResults(result.RunnerID)
	if responseErr != nil {
		return responseErr
	}

	runner, responseErr := rs.runnersRepository.GetRunner(result.RunnerID)
	if responseErr != nil {
		return responseErr
	}

	// Checking if the deleted result is personal best for the runner
	if runner.PersonalBest == result.RaceResult {
		runner.PersonalBest = getPersonalBestResult(results, resultId)
	}

	// Checking if the deleted result is season best for the runner
	currentYear := time.Now().Year()
	if runner.SeasonBest == result.RaceResult && result.Year == currentYear {
		runner.SeasonBest = getSeasonBestResult(results, result.Year, resultId)
	}

	responseErr = rs.resultsRepository.DeleteResult(result, runner)
	if responseErr != nil {
		return responseErr
	}

	return nil
}

func getPersonalBestResult(results []*models.Result, resultId string) string {
	var personalBest string
	for _, result := range results {
		if result.ID != resultId && (result.RaceResult < personalBest || personalBest == "") {
			personalBest = result.RaceResult
		}
	}

	return personalBest
}

func getSeasonBestResult(results []*models.Result, year int, resultId string) string {
	var seasonBest string
	for _, result := range results {
		if result.Year == year {
			if result.ID != resultId && (result.RaceResult < seasonBest || seasonBest == "") {
				seasonBest = result.RaceResult
			}
		}
	}

	return seasonBest
}
