package repositories

import (
	"context"
	"net/http"
	"runners-mongodb/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ResultsRepository struct {
	client *mongo.Client
}

func NewResultsRepository(client *mongo.Client) *ResultsRepository {
	return &ResultsRepository{
		client: client,
	}
}

func (rr ResultsRepository) CreateResult(result *models.Result, currentYear int) (*models.Result, *models.ResponseError) {
	objectId, err := primitive.ObjectIDFromHex(result.RunnerID)
	if err != nil {
		return nil, &models.ResponseError{
			Message: "Invalid runner ID",
			Status:  http.StatusBadRequest,
		}
	}

	result.ObjectID = primitive.NewObjectID()

	collection := rr.client.Database("runners_db").Collection("runners")
	filter := bson.D{{Key: "_id", Value: objectId}}
	update := bson.D{
		{Key: "$push", Value: bson.D{{Key: "results", Value: result}}},
		{Key: "$min", Value: bson.D{{Key: "personalbest", Value: result.RaceResult}}},
	}

	if result.Year == currentYear {
		update = bson.D{
			{Key: "$push", Value: bson.D{{Key: "results", Value: result}}},
			{Key: "$min", Value: bson.D{{Key: "personalbest", Value: result.RaceResult}}},
			{Key: "$min", Value: bson.D{{Key: "seasonbest", Value: result.RaceResult}}},
		}
	}

	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return &models.Result{
		ID:         result.ObjectID.Hex(),
		RunnerID:   result.RunnerID,
		RaceResult: result.RaceResult,
		Location:   result.Location,
		Position:   result.Position,
		Year:       result.Year,
	}, nil
}

func (rr ResultsRepository) UpdateRunnerResults(runner *models.Runner, resultId string) *models.ResponseError {
	runnerObjectId, err := primitive.ObjectIDFromHex(runner.ID)
	if err != nil {
		return &models.ResponseError{
			Message: "Invalid runner ID",
			Status:  http.StatusBadRequest,
		}
	}

	resultObjectId, err := primitive.ObjectIDFromHex(resultId)
	if err != nil {
		return &models.ResponseError{
			Message: "Invalid result ID",
			Status:  http.StatusBadRequest,
		}
	}

	collection := rr.client.Database("runners_db").Collection("runners")
	filter := bson.D{{Key: "_id", Value: runnerObjectId}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "personalbest", Value: runner.PersonalBest},
			{Key: "seasonbest", Value: runner.SeasonBest}}},
		{Key: "$pull", Value: bson.D{
			{Key: "results", Value: bson.D{{Key: "_id", Value: resultObjectId}}},
		}},
	}

	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}

func (rr ResultsRepository) GetRunnerByResultId(resultId string) (*models.Runner, *models.ResponseError) {
	objectId, err := primitive.ObjectIDFromHex(resultId)
	if err != nil {
		return nil, &models.ResponseError{
			Message: "Invalid result ID",
			Status:  http.StatusBadRequest,
		}
	}

	collection := rr.client.Database("runners_db").Collection("runners")
	filter := bson.D{{Key: "results", Value: bson.D{
		{Key: "$elemMatch", Value: bson.D{{Key: "_id", Value: objectId}}},
	}}}

	var runner *models.Runner
	err = collection.FindOne(context.TODO(), filter).Decode(&runner)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	runner.ID = runner.ObjectID.Hex()
	for _, result := range runner.Results {
		result.ID = result.ObjectID.Hex()
	}

	return runner, nil
}
