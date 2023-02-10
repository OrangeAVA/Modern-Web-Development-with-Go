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

type RunnersRepository struct {
	db *dynamodb.DynamoDB
}

func NewRunnersRepository(db *dynamodb.DynamoDB) *RunnersRepository {
	return &RunnersRepository{
		db: db,
	}
}

func (rr RunnersRepository) CreateRunner(runner *models.Runner) (*models.Runner, *models.ResponseError) {
	runner.ID = uuid.New().String()
	runner.IsActive = true

	runnerAttrMap, err := dynamodbattribute.MarshalMap(runner)
	if err != nil {
		return nil, &models.ResponseError{
			Message: "Failed to marshal runner into atribute-value map",
			Status:  http.StatusBadRequest,
		}
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("Runners"),
		Item:      runnerAttrMap,
	}

	_, err = rr.db.PutItem(input)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return runner, nil
}

func (rr RunnersRepository) UpdateRunner(runner *models.Runner) *models.ResponseError {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("Runners"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(runner.ID),
			},
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":fn": {
				S: aws.String(runner.FirstName),
			},
			":ln": {
				S: aws.String(runner.LastName),
			},
			":a": {
				N: aws.String(strconv.Itoa(runner.Age)),
			},
			":c": {
				S: aws.String(runner.Country),
			},
		},
		UpdateExpression: aws.String("SET first_name = :fn, last_name = :ln, age = :a, country = :c"),
	}

	_, err := rr.db.UpdateItem(input)
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}

func (rr RunnersRepository) DeleteRunner(runnerId string) *models.ResponseError {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("Runners"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(runnerId),
			},
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":a": {
				BOOL: aws.Bool(false),
			},
		},
		UpdateExpression: aws.String("SET is_active = :a"),
	}

	_, err := rr.db.UpdateItem(input)
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}

func (rr RunnersRepository) GetRunner(runnerId string) (*models.Runner, *models.ResponseError) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String("Runners"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(runnerId),
			},
		},
	}

	output, err := rr.db.GetItem(input)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	var runner *models.Runner
	err = dynamodbattribute.UnmarshalMap(output.Item, &runner)
	if err != nil {
		return nil, &models.ResponseError{
			Message: "Failed to unmarshal atribute-value map into runner",
			Status:  http.StatusInternalServerError,
		}
	}

	return runner, nil
}

func (rr RunnersRepository) GetAllRunners() ([]*models.Runner, *models.ResponseError) {
	input := &dynamodb.ScanInput{
		TableName: aws.String("Runners"),
	}

	output, err := rr.db.Scan(input)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	var runners []*models.Runner
	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &runners)
	if err != nil {
		return nil, &models.ResponseError{
			Message: "Failed to unmarshal atribute-value map into runners",
			Status:  http.StatusInternalServerError,
		}
	}

	return runners, nil
}

func (rr RunnersRepository) GetRunnersByCountry(country string) ([]*models.Runner, *models.ResponseError) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String("Runners"),
		IndexName:              aws.String("runners_global_index"),
		KeyConditionExpression: aws.String("country = :c"),
		FilterExpression:       aws.String("is_active = :a"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":c": {
				S: aws.String(country),
			},
			":a": {
				BOOL: aws.Bool(true),
			},
		},
		ScanIndexForward: aws.Bool(true),
	}

	output, err := rr.db.Query(input)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	var runners []*models.Runner
	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &runners)
	if err != nil {
		return nil, &models.ResponseError{
			Message: "Failed to unmarshal atribute-value map into runners",
			Status:  http.StatusInternalServerError,
		}
	}

	if len(runners) > 10 {
		return runners[:10], nil
	}

	return runners, nil
}

func (rr RunnersRepository) GetRunnersByIdMap(runnerIds map[string]struct{}) (map[string]models.Runner, *models.ResponseError) {
	var keys []map[string]*dynamodb.AttributeValue
	for runnerId := range runnerIds {
		key := map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(runnerId),
			},
		}

		keys = append(keys, key)
	}

	input := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			"Runners": {
				Keys: keys,
			},
		},
	}

	output, err := rr.db.BatchGetItem(input)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	runnersMap := make(map[string]models.Runner)
	for _, table := range output.Responses {
		for _, item := range table {
			var runner models.Runner
			err = dynamodbattribute.UnmarshalMap(item, &runner)
			if err != nil {
				return nil, &models.ResponseError{
					Message: "Failed to unmarshal atribute-value map into runner",
					Status:  http.StatusInternalServerError,
				}
			}

			runnersMap[runner.ID] = runner
		}
	}

	return runnersMap, nil
}
