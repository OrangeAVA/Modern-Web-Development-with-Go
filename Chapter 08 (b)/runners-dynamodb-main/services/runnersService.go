package services

import (
	"net/http"
	"runners-dynamodb/models"
	"runners-dynamodb/repositories"
	"strconv"
	"time"
)

type RunnersService struct {
	runnersRepository *repositories.RunnersRepository
	resultsRepository *repositories.ResultsRepository
}

func NewRunnersService(runnersRepository *repositories.RunnersRepository, resultsRepository *repositories.ResultsRepository) *RunnersService {
	return &RunnersService{
		runnersRepository: runnersRepository,
		resultsRepository: resultsRepository,
	}
}

func (rs RunnersService) CreateRunner(runner *models.Runner) (*models.Runner, *models.ResponseError) {
	responseErr := validateRunner(runner)
	if responseErr != nil {
		return nil, responseErr
	}

	return rs.runnersRepository.CreateRunner(runner)
}

func (rs RunnersService) UpdateRunner(runner *models.Runner) *models.ResponseError {
	responseErr := validateRunnerId(runner.ID)
	if responseErr != nil {
		return responseErr
	}

	responseErr = validateRunner(runner)
	if responseErr != nil {
		return responseErr
	}

	return rs.runnersRepository.UpdateRunner(runner)
}

func (rs RunnersService) DeleteRunner(runnerId string) *models.ResponseError {
	responseErr := validateRunnerId(runnerId)
	if responseErr != nil {
		return responseErr
	}

	return rs.runnersRepository.DeleteRunner(runnerId)
}

func (rs RunnersService) GetRunner(runnerId string) (*models.Runner, *models.ResponseError) {
	responseErr := validateRunnerId(runnerId)
	if responseErr != nil {
		return nil, responseErr
	}

	runner, responseErr := rs.runnersRepository.GetRunner(runnerId)
	if responseErr != nil {
		return nil, responseErr
	}

	results, responseErr := rs.resultsRepository.GetAllRunnersResults(runnerId)
	if responseErr != nil {
		return nil, responseErr
	}

	runner.Results = results

	return runner, nil
}

func (rs RunnersService) GetRunnersBatch(country string, year string) ([]*models.Runner, *models.ResponseError) {
	if country != "" && year != "" {
		return nil, &models.ResponseError{
			Message: "Only one parameter, country or year, can be passed",
			Status:  http.StatusBadRequest,
		}
	}

	if country != "" {
		return rs.runnersRepository.GetRunnersByCountry(country)
	}

	if year != "" {
		intYear, err := strconv.Atoi(year)
		if err != nil {
			return nil, &models.ResponseError{
				Message: "Invalid year",
				Status:  http.StatusBadRequest,
			}
		}

		currentYear := time.Now().Year()
		if intYear < 0 || intYear > currentYear {
			return nil, &models.ResponseError{
				Message: "Invalid year",
				Status:  http.StatusBadRequest,
			}
		}

		results, responseErr := rs.resultsRepository.GetTenBestResultsFromYear(intYear)
		if responseErr != nil {
			return nil, responseErr
		}

		runnerIdMap := make(map[string]struct{})
		for _, result := range results {
			runnerIdMap[result.RunnerID] = struct{}{}
		}

		runnersMap, responseErr := rs.runnersRepository.GetRunnersByIdMap(runnerIdMap)
		if responseErr != nil {
			return nil, responseErr
		}

		// Join
		var runners []*models.Runner
		for _, result := range results {
			runner := runnersMap[result.RunnerID]
			runner.SeasonBest = result.RaceResult
			runners = append(runners, &runner)
		}

		return runners, nil
	}

	return rs.runnersRepository.GetAllRunners()
}

func validateRunner(runner *models.Runner) *models.ResponseError {
	if runner.FirstName == "" {
		return &models.ResponseError{
			Message: "Invalid first name",
			Status:  http.StatusBadRequest,
		}
	}

	if runner.LastName == "" {
		return &models.ResponseError{
			Message: "Invalid last name",
			Status:  http.StatusBadRequest}
	}

	if runner.Age < 0 || runner.Age > 125 {
		return &models.ResponseError{
			Message: "Invalid age",
			Status:  http.StatusBadRequest,
		}
	}

	if runner.Country == "" {
		return &models.ResponseError{
			Message: "Invalid country",
			Status:  http.StatusBadRequest,
		}
	}

	return nil
}

func validateRunnerId(runnerId string) *models.ResponseError {
	if runnerId == "" {
		return &models.ResponseError{
			Message: "Invalid runner ID",
			Status:  http.StatusBadRequest,
		}
	}

	return nil
}
