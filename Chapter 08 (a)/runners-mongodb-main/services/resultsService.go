package services

import (
	"net/http"
	"runners-mongodb/models"
	"runners-mongodb/repositories"
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

	response, responseErr := rs.resultsRepository.CreateResult(result, currentYear)
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

	runner, responseErr := rs.resultsRepository.GetRunnerByResultId(resultId)
	if responseErr != nil {
		return responseErr
	}

	result := findResult(runner.Results, resultId)

	// Checking if the deleted result is personal best for the runner
	if runner.PersonalBest == result.RaceResult {
		runner.PersonalBest = getPersonalBestResult(runner.Results, resultId)
	}

	// Checking if the deleted result is season best for the runner
	currentYear := time.Now().Year()
	if runner.SeasonBest == result.RaceResult &&
		result.Year == currentYear {
		runner.SeasonBest = getSeasonBestResult(runner.Results, result.Year, resultId)
	}

	responseErr = rs.resultsRepository.UpdateRunnerResults(runner, resultId)
	if responseErr != nil {
		return responseErr
	}

	return nil
}

func findResult(results []*models.Result, resultId string) *models.Result {
	for _, result := range results {
		if result.ID == resultId {
			return result
		}
	}

	return nil
}

func getPersonalBestResult(results []*models.Result, resultId string) string {
	personalBest := results[0].RaceResult
	for _, result := range results {
		if result.ID != resultId && result.RaceResult < personalBest {
			personalBest = result.RaceResult
		}
	}

	return personalBest
}

func getSeasonBestResult(results []*models.Result, year int, resultId string) string {
	var seasonBest string
	for _, result := range results {
		if result.Year == year {
			if result.ID != resultId && result.RaceResult < seasonBest {
				seasonBest = result.RaceResult
			}
		}
	}

	return seasonBest
}
