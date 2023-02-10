package repositories

import (
	"context"
	"net/http"
	"runners-mongodb/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RunnersRepository struct {
	client *mongo.Client
}

func NewRunnersRepository(client *mongo.Client) *RunnersRepository {
	return &RunnersRepository{
		client: client,
	}
}

func (rr RunnersRepository) CreateRunner(runner *models.Runner) (*models.Runner, *models.ResponseError) {
	collection := rr.client.Database("runners_db").Collection("runners")

	runner.IsActive = true

	result, err := collection.InsertOne(context.TODO(), runner)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	id := result.InsertedID.(primitive.ObjectID)

	return &models.Runner{
		ID:        id.Hex(),
		FirstName: runner.FirstName,
		LastName:  runner.LastName,
		Age:       runner.Age,
		IsActive:  runner.IsActive,
		Country:   runner.Country,
	}, nil
}

func (rr RunnersRepository) UpdateRunner(runner *models.Runner) *models.ResponseError {
	objectId, err := primitive.ObjectIDFromHex(runner.ID)
	if err != nil {
		return &models.ResponseError{
			Message: "Invalid runner ID",
			Status:  http.StatusBadRequest,
		}
	}

	collection := rr.client.Database("runners_db").Collection("runners")
	filter := bson.D{{Key: "_id", Value: objectId}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "firstname", Value: runner.FirstName},
		{Key: "lastname", Value: runner.LastName},
		{Key: "age", Value: runner.Age},
		{Key: "country", Value: runner.Country},
	}}}

	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}

func (rr RunnersRepository) DeleteRunner(runnerId string) *models.ResponseError {
	objectId, err := primitive.ObjectIDFromHex(runnerId)
	if err != nil {
		return &models.ResponseError{
			Message: "Invalid runner ID",
			Status:  http.StatusBadRequest,
		}
	}

	collection := rr.client.Database("runners_db").Collection("runners")
	filter := bson.D{{Key: "_id", Value: objectId}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "isactive", Value: false}}}}

	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}

func (rr RunnersRepository) GetRunner(runnerId string) (*models.Runner, *models.ResponseError) {
	objectId, err := primitive.ObjectIDFromHex(runnerId)
	if err != nil {
		return nil, &models.ResponseError{
			Message: "Invalid runner ID",
			Status:  http.StatusBadRequest,
		}
	}

	collection := rr.client.Database("runners_db").Collection("runners")
	filter := bson.D{{Key: "_id", Value: objectId}}

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

func (rr RunnersRepository) GetAllRunners() ([]*models.Runner, *models.ResponseError) {
	collection := rr.client.Database("runners_db").Collection("runners")
	filter := bson.D{}
	options := options.Find().SetProjection(bson.D{{Key: "results", Value: 0}})

	cursor, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	var runners []*models.Runner
	for cursor.Next(context.TODO()) {
		var runner *models.Runner
		err = cursor.Decode(&runner)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
		runner.ID = runner.ObjectID.Hex()
		runners = append(runners, runner)
	}

	if err := cursor.Err(); err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return runners, nil
}

func (rr RunnersRepository) GetRunnersByCountry(country string) ([]*models.Runner, *models.ResponseError) {
	collection := rr.client.Database("runners_db").Collection("runners")
	filter := bson.D{{Key: "country", Value: country}, {Key: "isactive", Value: true}}
	options := options.Find().
		SetProjection(bson.D{{Key: "results", Value: 0}}).
		SetSort(bson.D{{Key: "personalbest", Value: 1}}).
		SetLimit(10)

	cursor, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	var runners []*models.Runner
	for cursor.Next(context.TODO()) {
		var runner *models.Runner
		err = cursor.Decode(&runner)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
		runner.ID = runner.ObjectID.Hex()
		runners = append(runners, runner)
	}

	if err := cursor.Err(); err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return runners, nil
}

type unwindRunner struct {
	ID           primitive.ObjectID `json:"-" bson:"_id"`
	FirstName    string
	LastName     string
	Age          int
	IsActive     bool
	Country      string
	PersonalBest string
	SeasonBest   string
	Results      *models.Result
}

func (rr RunnersRepository) GetRunnersByYear(year int) ([]*models.Runner, *models.ResponseError) {
	collection := rr.client.Database("runners_db").Collection("runners")
	unwindStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$results"}}}}
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "results.year", Value: year}}}}
	sortStage := bson.D{{Key: "$sort", Value: bson.D{{Key: "results.raceresult", Value: 1}}}}
	limitStage := bson.D{{Key: "$limit", Value: 10}}

	options := options.Aggregate().SetHint("results.year_1")

	pipeline := mongo.Pipeline{unwindStage, matchStage, sortStage, limitStage}

	cursor, err := collection.Aggregate(context.TODO(), pipeline, options)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	var unwindRunners []unwindRunner
	err = cursor.All(context.TODO(), &unwindRunners)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	runners := make([]*models.Runner, 0)
	for _, ur := range unwindRunners {

		runner := &models.Runner{
			ID:           ur.ID.Hex(),
			FirstName:    ur.FirstName,
			LastName:     ur.LastName,
			Age:          ur.Age,
			IsActive:     ur.IsActive,
			Country:      ur.Country,
			PersonalBest: ur.PersonalBest,
			SeasonBest:   ur.Results.RaceResult,
		}

		runners = append(runners, runner)
	}

	return runners, nil
}
