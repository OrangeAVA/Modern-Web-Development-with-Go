package repositories

import (
	"net/http"
	"runners-dynamodb/models"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

type ResultsRepository struct {
	db *dynamodb.DynamoDB
}

func NewResultsRepository(db *dynamodb.DynamoDB) *ResultsRepository {
	return &ResultsRepository{
		db: db,
	}
}

func (rr ResultsRepository) CreateResult(result *models.Result, runner *models.Runner) (*models.Result, *models.ResponseError) {
	result.ID = uuid.New().String()

	resultAttrMap, err := dynamodbattribute.MarshalMap(result)
	if err != nil {
		return nil, &models.ResponseError{
			Message: "Failed to marshal result into atribute-value map",
			Status:  http.StatusBadRequest,
		}
	}

	input := &dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{
			{
				Put: &dynamodb.Put{
					TableName: aws.String("Results"),
					Item:      resultAttrMap,
				},
			},
			{
				Update: &dynamodb.Update{
					TableName: aws.String("Runners"),
					Key: map[string]*dynamodb.AttributeValue{
						"id": {
							S: aws.String(runner.ID),
						},
					},
					ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
						":pb": {
							S: aws.String(runner.PersonalBest),
						},
						":sb": {
							S: aws.String(runner.SeasonBest),
						},
					},
					UpdateExpression: aws.String("SET personal_best = :pb, season_best = :sb"),
				},
			},
		},
	}

	_, err = rr.db.TransactWriteItems(input)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return result, nil
}

func (rr ResultsRepository) DeleteResult(result *models.Result, runner *models.Runner) *models.ResponseError {

	input := &dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{
			{
				Delete: &dynamodb.Delete{
					TableName: aws.String("Results"),
					Key: map[string]*dynamodb.AttributeValue{
						"runner_id": {
							S: aws.String(result.RunnerID),
						},
						"race_result": {
							S: aws.String(result.RaceResult),
						},
					},
				},
			},
			{
				Update: &dynamodb.Update{
					TableName: aws.String("Runners"),
					Key: map[string]*dynamodb.AttributeValue{
						"id": {
							S: aws.String(runner.ID),
						},
					},
					ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
						":pb": {
							S: aws.String(runner.PersonalBest),
						},
						":sb": {
							S: aws.String(runner.SeasonBest),
						},
					},
					UpdateExpression: aws.String("SET personal_best = :pb, season_best = :sb"),
				},
			},
		},
	}

	_, err := rr.db.TransactWriteItems(input)
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}

func (rr ResultsRepository) GetResult(resultId string) (*models.Result, *models.ResponseError) {
	input := &dynamodb.ScanInput{
		TableName:        aws.String("Results"),
		FilterExpression: aws.String("id = :i"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":i": {
				S: aws.String(resultId),
			},
		},
	}

	output, err := rr.db.Scan(input)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	var results []*models.Result
	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &results)
	if err != nil {
		return nil, &models.ResponseError{
			Message: "Failed to unmarshal atribute-value map into result",
			Status:  http.StatusInternalServerError,
		}
	}

	return results[0], nil
}

func (rr ResultsRepository) GetAllRunnersResults(runnerId string) ([]*models.Result, *models.ResponseError) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String("Results"),
		KeyConditionExpression: aws.String("runner_id = :rid"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":rid": {
				S: aws.String(runnerId),
			},
		},
	}

	output, err := rr.db.Query(input)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	var results []*models.Result
	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &results)
	if err != nil {
		return nil, &models.ResponseError{
			Message: "Failed to unmarshal atribute-value map into results",
			Status:  http.StatusInternalServerError,
		}
	}

	return results, nil
}

func (rr ResultsRepository) GetTenBestResultsFromYear(year int) ([]*models.Result, *models.ResponseError) {
	input := &dynamodb.ScanInput{
		TableName:        aws.String("Results"),
		FilterExpression: aws.String("race_year = :ry"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":ry": {
				N: aws.String(strconv.Itoa(year)),
			},
		},
	}

	output, err := rr.db.Scan(input)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	var results []*models.Result
	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &results)
	if err != nil {
		return nil, &models.ResponseError{
			Message: "Failed to unmarshal atribute-value map into results",
			Status:  http.StatusInternalServerError,
		}
	}

	if len(results) > 10 {
		return results[:10], nil
	}

	return results, nil
}
