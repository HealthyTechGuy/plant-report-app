package plantservice

import (
	"errors"
	"fmt"

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
	GetPlantInfo(plantID string) (*PlantInfo, error)
}

// PlantService is a concrete implementation of PlantServiceInterface
type PlantService struct {
	dynamoDBClient dynamodbiface.DynamoDBAPI
	tableName      string
}

// PlantInfo holds information about a plant
type PlantInfo struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	GrowingPeriod   string `json:"growing_period"`
	OptimalPlanting string `json:"optimal_planting"`
	HardinessZone   string `json:"hardiness_zone"`
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
func (s *PlantService) GetPlantInfo(plantID string) (*PlantInfo, error) {
	result, err := s.dynamoDBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(s.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(plantID),
			},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get item from DynamoDB: %w", err)
	}

	if result.Item == nil {
		return nil, ErrPlantNotFound
	}

	plantInfo := &PlantInfo{
		ID:              *result.Item["id"].S,
		Name:            *result.Item["name"].S,
		GrowingPeriod:   *result.Item["growing_period"].S,
		OptimalPlanting: *result.Item["optimal_planting"].S,
		HardinessZone:   *result.Item["hardiness_zone"].S,
	}

	return plantInfo, nil
}
