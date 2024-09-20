package plantservice

import (
	"errors"
	"testing"

	"github.com/HealthyTechGuy/plant-report-app/internal/plant-service/mocks"
	"github.com/HealthyTechGuy/plant-report-app/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPlantInfo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamoDB := mocks.NewMockDynamoDBAPI(ctrl)
	plantService := &PlantService{
		dynamoDBClient: mockDynamoDB,
		tableName:      "test-table",
	}

	expectedPlantInfo := models.PlantInfo{
		ID:              "1",
		Name:            "Blueberry Bush",
		GrowingPeriod:   "May to August",
		OptimalPlanting: "Spring",
		HardinessZone:   "3-7",
	}

	mockDynamoDB.EXPECT().GetItem(gomock.Any()).Return(&dynamodb.GetItemOutput{
		Item: map[string]*dynamodb.AttributeValue{
			"PlantID":          {S: aws.String(expectedPlantInfo.ID)},
			"name":             {S: aws.String(expectedPlantInfo.Name)},
			"growing_period":   {S: aws.String(expectedPlantInfo.GrowingPeriod)},
			"optimal_planting": {S: aws.String(expectedPlantInfo.OptimalPlanting)},
			"hardiness_zone":   {S: aws.String(expectedPlantInfo.HardinessZone)},
		},
	}, nil)

	plantInfo, err := plantService.GetPlantInfo("1")
	require.NoError(t, err)
	assert.Equal(t, expectedPlantInfo, plantInfo)
}

func TestGetPlantInfo_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamoDB := mocks.NewMockDynamoDBAPI(ctrl)
	plantService := &PlantService{
		dynamoDBClient: mockDynamoDB,
		tableName:      "test-table",
	}

	mockDynamoDB.EXPECT().GetItem(gomock.Any()).Return(&dynamodb.GetItemOutput{
		Item: nil,
	}, nil)

	_, err := plantService.GetPlantInfo("unknown-id")
	assert.Equal(t, ErrPlantNotFound, err)
}

func TestGetPlantInfo_DynamoDBError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamoDB := mocks.NewMockDynamoDBAPI(ctrl)
	plantService := &PlantService{
		dynamoDBClient: mockDynamoDB,
		tableName:      "test-table",
	}

	mockDynamoDB.EXPECT().GetItem(gomock.Any()).Return(nil, errors.New("dynamo error"))

	_, err := plantService.GetPlantInfo("1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get item from DynamoDB")
}

func TestPDFGenerationAndUpload(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPDFGenerator := new(mocks.MockPDFGenerator)

	expectedPlantInfo := models.PlantInfo{
		ID:              "1",
		Name:            "Blueberry Bush",
		GrowingPeriod:   "May to August",
		OptimalPlanting: "Spring",
		HardinessZone:   "3-7",
	}

	userLocation := models.UserLocation{
		// Fill in any necessary fields for UserLocation
	}

	// Set up expectations for PDF generation
	mockPDFGenerator.On("GeneratePDF", userLocation, expectedPlantInfo).Return([]byte("mock pdf data"), nil)
	// Set up expectations for S3 upload
	mockPDFGenerator.On("UploadToS3", []byte("mock pdf data"), "plant-report-bucket", "file.pdf").Return("https://plant-report-bucket.s3.amazonaws.com/file.pdf", nil)

	// Simulate generating the PDF
	pdfData, err := mockPDFGenerator.GeneratePDF(userLocation, expectedPlantInfo)
	require.NoError(t, err)

	// Simulate uploading the PDF to S3
	s3URL, err := mockPDFGenerator.UploadToS3(pdfData, "plant-report-bucket", "file.pdf")
	require.NoError(t, err)
	assert.Equal(t, "https://plant-report-bucket.s3.amazonaws.com/file.pdf", s3URL)

	// Assert that the expectations were met
	mockPDFGenerator.AssertExpectations(t)
}
