package repositories

import (
	"database/sql"
	"net/http"
	"runners-mysql/models"
	"strconv"
)

type ResultsRepository struct {
	dbHandler   *sql.DB
	transaction *sql.Tx
}

func NewResultsRepository(dbHAndler *sql.DB) *ResultsRepository {
	return &ResultsRepository{
		dbHandler: dbHAndler,
	}
}

func (rr ResultsRepository) CreateResult(result *models.Result) (*models.Result, *models.ResponseError) {
	query := `
		INSERT INTO results(runner_id, race_result, location, position, result_year)
		VALUES (?, ?, ?, ?, ?)`

	res, err := rr.transaction.Exec(query, result.RunnerID, result.RaceResult, result.Location, result.Position, result.Year)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	resultId, err := res.LastInsertId()
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return &models.Result{
		ID:         strconv.FormatInt(resultId, 10),
		RunnerID:   result.RunnerID,
		RaceResult: result.RaceResult,
		Location:   result.Location,
		Position:   result.Position,
		Year:       result.Year,
	}, nil
}

func (rr ResultsRepository) DeleteResult(resultId string) *models.ResponseError {
	query := `
		DELETE FROM results
		WHERE id = ?`

	res, err := rr.transaction.Exec(query, resultId)
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	if rowsAffected == 0 {
		return &models.ResponseError{
			Message: "Result not found",
			Status:  http.StatusNotFound,
		}
	}

	return nil
}

func (rr ResultsRepository) GetAllRunnersResults(runnerId string) ([]*models.Result, *models.ResponseError) {
	query := `
	SELECT id, race_result, location, position, result_year
	FROM results
	WHERE runner_id = ?`

	rows, err := rr.dbHandler.Query(query, runnerId)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	results := make([]*models.Result, 0)
	var id, raceResult, location string
	var position, year int

	for rows.Next() {
		err := rows.Scan(&id, &raceResult, &location, &position, &year)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}

		result := &models.Result{
			ID:         id,
			RunnerID:   runnerId,
			RaceResult: raceResult,
			Location:   location,
			Position:   position,
			Year:       year,
		}

		results = append(results, result)
	}

	if rows.Err() != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return results, nil
}

func (rr ResultsRepository) GetPersonalBestResults(runnerId string) (string, *models.ResponseError) {
	query := `
	SELECT MIN(race_result)
	FROM results
	WHERE runner_id = ?`

	rows, err := rr.dbHandler.Query(query, runnerId)
	if err != nil {
		return "", &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	var raceResult string

	for rows.Next() {
		err := rows.Scan(&raceResult)
		if err != nil {
			return "", &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
	}

	if rows.Err() != nil {
		return "", &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return raceResult, nil
}

func (rr ResultsRepository) GetSeasonBestResults(runnerId string, year int) (string, *models.ResponseError) {
	query := `
	SELECT MIN(race_result)
	FROM results
	WHERE runner_id = ? AND result_year = ?`

	rows, err := rr.dbHandler.Query(query, runnerId, year)
	if err != nil {
		return "", &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	var raceResult string

	for rows.Next() {
		err := rows.Scan(&raceResult)
		if err != nil {
			return "", &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
	}

	if rows.Err() != nil {
		return "", &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return raceResult, nil
}

func (rr ResultsRepository) GetResult(resultId string) (*models.Result, *models.ResponseError) {
	query := `
	SELECT *
	FROM results
	WHERE id = ?;`

	rows, err := rr.dbHandler.Query(query, resultId)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	var id, runnerId, raceResult, location string
	var position, year int

	for rows.Next() {
		err := rows.Scan(&id, &runnerId, &raceResult, &location, &position, &year)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
	}

	if rows.Err() != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return &models.Result{
		ID:         id,
		RunnerID:   runnerId,
		RaceResult: raceResult,
		Location:   location,
		Position:   position,
		Year:       year,
	}, nil
}
