package plantservice

import (
	"errors"
	"fmt"

	"github.com/HealthyTechGuy/plant-report-app/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	// Define the possible plant IDs
	ErrPlantNotFound = errors.New("plant not found")
)

// PlantServiceInterface defines the methods for interacting with plant data
type PlantServiceInterface interface {
	GetPlantInfo(plantID string) (models.PlantInfo, error)
}

// PlantService is a concrete implementation of PlantServiceInterface
type PlantService struct {
	dynamoDBClient dynamodbiface.DynamoDBAPI
	tableName      string
}

// NewPlantService creates a new PlantService
func NewPlantService(tableName string) *PlantService {
	sess := session.Must(session.NewSession())
	return &PlantService{
		dynamoDBClient: dynamodb.New(sess),
		tableName:      tableName,
	}
}

// GetPlantInfo retrieves plant information from DynamoDB
func (s *PlantService) GetPlantInfo(plantID string) (pi models.PlantInfo, err error) {
	result, err := s.dynamoDBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(s.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"PlantID": {
				S: aws.String(plantID),
			},
		},
	})

	if err != nil {
		return pi, fmt.Errorf("failed to get item from DynamoDB: %w", err)
	}

	if result.Item == nil {
		return pi, ErrPlantNotFound
	}

	plantInfo := models.PlantInfo{
		ID:              *result.Item["PlantID"].S,
		Name:            *result.Item["name"].S,
		GrowingPeriod:   *result.Item["growing_period"].S,
		OptimalPlanting: *result.Item["optimal_planting"].S,
		HardinessZone:   *result.Item["hardiness_zone"].S,
	}

	return plantInfo, nil
}
