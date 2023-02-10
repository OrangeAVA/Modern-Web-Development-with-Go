package services

import (
	"net/http"
	"runners-postgresql/models"
	"runners-postgresql/repositories"
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

	raceResult, err := parseRaceResult(result.RaceResult)
	if err != nil {
		return nil, &models.ResponseError{
			Message: "Invalid race result",
			Status:  http.StatusBadRequest,
		}
	}

	err = repositories.BeginTransaction(rs.runnersRepository, rs.resultsRepository)
	if err != nil {
		return nil, &models.ResponseError{
			Message: "Failed to start transaction",
			Status:  http.StatusBadRequest,
		}
	}

	response, responseErr := rs.resultsRepository.CreateResult(result)
	if responseErr != nil {
		repositories.RollbackTransaction(rs.runnersRepository, rs.resultsRepository)
		return nil, responseErr
	}

	runner, responseErr := rs.runnersRepository.GetRunner(result.RunnerID)
	if responseErr != nil {
		repositories.RollbackTransaction(rs.runnersRepository, rs.resultsRepository)
		return nil, responseErr
	}

	if runner == nil {
		repositories.RollbackTransaction(rs.runnersRepository, rs.resultsRepository)
		return nil, &models.ResponseError{
			Message: "Runner not found",
			Status:  http.StatusNotFound,
		}
	}

	// update runners personal best
	if runner.PersonalBest == "" {
		runner.PersonalBest = result.RaceResult
	} else {
		personalBest, err := parseRaceResult(runner.PersonalBest)
		if err != nil {
			repositories.RollbackTransaction(rs.runnersRepository, rs.resultsRepository)
			return nil, &models.ResponseError{
				Message: "Failed to parse personal best",
				Status:  http.StatusInternalServerError,
			}
		}

		if raceResult < personalBest {
			runner.PersonalBest = result.RaceResult
		}
	}

	// update runners seeason best
	if result.Year == currentYear {
		if runner.SeasonBest == "" {
			runner.SeasonBest = result.RaceResult
		} else {
			seasonBest, err := parseRaceResult(runner.SeasonBest)
			if err != nil {
				repositories.RollbackTransaction(rs.runnersRepository, rs.resultsRepository)
				return nil, &models.ResponseError{
					Message: "Failed to parse season best",
					Status:  http.StatusInternalServerError,
				}
			}

			if raceResult < seasonBest {
				runner.SeasonBest = result.RaceResult
			}
		}
	}

	responseErr = rs.runnersRepository.UpdateRunnerResults(runner)
	if responseErr != nil {
		repositories.RollbackTransaction(rs.runnersRepository, rs.resultsRepository)
		return nil, responseErr
	}

	repositories.CommitTransaction(rs.runnersRepository, rs.resultsRepository)
	return response, nil
}

func (rs ResultsService) DeleteResult(resultId string) *models.ResponseError {
	if resultId == "" {
		return &models.ResponseError{
			Message: "Invalid result ID",
			Status:  http.StatusBadRequest,
		}
	}

	err := repositories.BeginTransaction(rs.runnersRepository, rs.resultsRepository)
	if err != nil {
		return &models.ResponseError{
			Message: "Failed to start transaction",
			Status:  http.StatusBadRequest,
		}
	}

	result, responseErr := rs.resultsRepository.DeleteResult(resultId)
	if responseErr != nil {
		return responseErr
	}

	runner, responseErr := rs.runnersRepository.GetRunner(result.RunnerID)
	if responseErr != nil {
		repositories.RollbackTransaction(rs.runnersRepository, rs.resultsRepository)
		return responseErr
	}

	// Checking if the deleted result is personal best for the runner
	if runner.PersonalBest == result.RaceResult {
		personalBest, responseErr := rs.resultsRepository.GetPersonalBestResults(result.RunnerID)
		if responseErr != nil {
			repositories.RollbackTransaction(rs.runnersRepository, rs.resultsRepository)
			return responseErr
		}
		runner.PersonalBest = personalBest
	}

	// Checking if the deleted result is season best for the runner
	currentYear := time.Now().Year()
	if runner.SeasonBest == result.RaceResult && result.Year == currentYear {
		seasonBest, responseErr := rs.resultsRepository.GetSeasonBestResults(result.RunnerID, result.Year)
		if responseErr != nil {
			repositories.RollbackTransaction(rs.runnersRepository, rs.resultsRepository)
			return responseErr
		}
		runner.SeasonBest = seasonBest
	}

	responseErr = rs.runnersRepository.UpdateRunnerResults(runner)
	if responseErr != nil {
		repositories.RollbackTransaction(rs.runnersRepository, rs.resultsRepository)
		return responseErr
	}

	repositories.CommitTransaction(rs.runnersRepository, rs.resultsRepository)

	return nil
}

func parseRaceResult(timeString string) (time.Duration, error) {
	return time.ParseDuration(
		timeString[0:2] + "h" +
			timeString[3:5] + "m" +
			timeString[6:8] + "s")
}
