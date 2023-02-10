package repositories

import (
	"database/sql"
	"fmt"
	"net/http"
	"runners-mysql/models"
	"strconv"
)

type RunnersRepository struct {
	dbHandler   *sql.DB
	transaction *sql.Tx
}

func NewRunnersRepository(dbHandler *sql.DB) *RunnersRepository {
	return &RunnersRepository{
		dbHandler: dbHandler,
	}
}

func (rr RunnersRepository) CreateRunner(runner *models.Runner) (*models.Runner, *models.ResponseError) {
	query := `
		INSERT INTO runners(first_name, last_name, age, country)
		VALUES (?, ?, ?, ?)`

	res, err := rr.dbHandler.Exec(query, runner.FirstName, runner.LastName, runner.Age, runner.Country)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	runnerId, err := res.LastInsertId()
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return &models.Runner{
		ID:        strconv.FormatInt(runnerId, 10),
		FirstName: runner.FirstName,
		LastName:  runner.LastName,
		Age:       runner.Age,
		IsActive:  true,
		Country:   runner.Country,
	}, nil
}

func (rr RunnersRepository) UpdateRunner(runner *models.Runner) *models.ResponseError {
	query := `
		UPDATE runners
		SET
			first_name = ?,
			last_name = ?,
			age = ?,
			country = ?
		WHERE id = ?`

	res, err := rr.dbHandler.Exec(query, runner.FirstName, runner.LastName, runner.Age, runner.Country, runner.ID)
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
			Message: "Runner not found",
			Status:  http.StatusNotFound,
		}
	}

	return nil
}

func (rr RunnersRepository) UpdateRunnerResults(runner *models.Runner) *models.ResponseError {
	query := `
		UPDATE runners
		SET
			personal_best = ?,
			season_best = ?
		WHERE id = ?`

	_, err := rr.transaction.Exec(query, runner.PersonalBest, runner.SeasonBest, runner.ID)
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}

func (rr RunnersRepository) DeleteRunner(runnerId string) *models.ResponseError {
	query := `UPDATE runners SET is_active = FALSE WHERE id = ?`

	res, err := rr.dbHandler.Exec(query, runnerId)
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
			Message: "Runner not found",
			Status:  http.StatusNotFound,
		}
	}

	return nil
}

func (rr RunnersRepository) GetRunner(runnerId string) (*models.Runner, *models.ResponseError) {
	fmt.Println(runnerId)
	query := `
		SELECT *
		FROM runners
		WHERE id = ?`

	rows, err := rr.dbHandler.Query(query, runnerId)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	var id, firstName, lastName, country string
	var personalBest, seasonBest sql.NullString
	var age int
	var isActive bool
	for rows.Next() {
		err := rows.Scan(&id, &firstName, &lastName, &age, &isActive, &country, &personalBest, &seasonBest)
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

	return &models.Runner{
		ID:           id,
		FirstName:    firstName,
		LastName:     lastName,
		Age:          age,
		IsActive:     isActive,
		Country:      country,
		PersonalBest: personalBest.String,
		SeasonBest:   seasonBest.String,
	}, nil
}

func (rr RunnersRepository) GetAllRunners() ([]*models.Runner, *models.ResponseError) {
	query := `
	SELECT *
	FROM runners`

	rows, err := rr.dbHandler.Query(query)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	runners := make([]*models.Runner, 0)
	var id, firstName, lastName, country string
	var personalBest, seasonBest sql.NullString
	var age int
	var isActive bool

	for rows.Next() {
		err := rows.Scan(&id, &firstName, &lastName, &age, &isActive, &country, &personalBest, &seasonBest)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}

		runner := &models.Runner{
			ID:           id,
			FirstName:    firstName,
			LastName:     lastName,
			Age:          age,
			IsActive:     isActive,
			Country:      country,
			PersonalBest: personalBest.String,
			SeasonBest:   seasonBest.String,
		}

		runners = append(runners, runner)
	}

	if rows.Err() != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return runners, nil
}

func (rr RunnersRepository) GetRunnersByCountry(country string) ([]*models.Runner, *models.ResponseError) {
	query := `
	SELECT id, first_name, last_name, age, personal_best, season_best
	FROM runners
	WHERE country = ? AND is_active = TRUE
	ORDER BY personal_best
	LIMIT 10`

	rows, err := rr.dbHandler.Query(query, country)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	runners := make([]*models.Runner, 0)
	var id, firstName, lastName string
	var personalBest, seasonBest sql.NullString
	var age int

	for rows.Next() {
		err := rows.Scan(&id, &firstName, &lastName, &age, &personalBest, &seasonBest)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}

		runner := &models.Runner{
			ID:           id,
			FirstName:    firstName,
			LastName:     lastName,
			Age:          age,
			IsActive:     true,
			Country:      country,
			PersonalBest: personalBest.String,
			SeasonBest:   seasonBest.String,
		}

		runners = append(runners, runner)
	}

	if rows.Err() != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return runners, nil
}

func (rr RunnersRepository) GetRunnersByYear(year int) ([]*models.Runner, *models.ResponseError) {
	query := `
	SELECT runners.id, runners.first_name, runners.last_name, runners.age, runners.is_active, runners.country, runners.personal_best, results.race_result
	FROM runners
	INNER JOIN (
		SELECT runner_id, MIN(race_result) as race_result
		FROM results
		WHERE result_year = ?
		GROUP BY runner_id) results
	ON runners.id = results.runner_id
	ORDER BY results.race_result
	LIMIT 10`

	rows, err := rr.dbHandler.Query(query, year)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	runners := make([]*models.Runner, 0)
	var id, firstName, lastName, country string
	var personalBest, seasonBest sql.NullString
	var age int
	var isActive bool

	for rows.Next() {
		err := rows.Scan(&id, &firstName, &lastName, &age, &isActive, &country, &personalBest, &seasonBest)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}

		runner := &models.Runner{
			ID:           id,
			FirstName:    firstName,
			LastName:     lastName,
			Age:          age,
			IsActive:     isActive,
			Country:      country,
			PersonalBest: personalBest.String,
			SeasonBest:   seasonBest.String,
		}

		runners = append(runners, runner)
	}

	if rows.Err() != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return runners, nil
}
