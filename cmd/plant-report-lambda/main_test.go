package main

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/HealthyTechGuy/plant-report-app/internal/plant-service/mocks"
	"github.com/HealthyTechGuy/plant-report-app/models"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleRequest_Success(t *testing.T) {
	mockPlantService := new(mocks.MockPlantService)
	mockPDFGenerator := new(mocks.MockPDFGenerator)

	// Mock plant info
	plantInfo := models.PlantInfo{
		ID:              "blueberry",
		Name:            "Blueberry Bush",
		GrowingPeriod:   "6-8 months",
		OptimalPlanting: "Spring",
		HardinessZone:   "5-7",
	}

	// Mock plant service response
	mockPlantService.On("GetPlantInfo", "blueberry").Return(plantInfo, nil)

	// Mock PDF generation with []byte return type
	mockPDFGenerator.On("GeneratePDF", mock.Anything, plantInfo).Return([]byte("PDF content"), nil)

	// Mock S3 upload
	mockPDFGenerator.On("UploadToS3", []byte("PDF content"), "plant-report-bucket", "file.pdf").Return("https://plant-report-bucket.s3.amazonaws.com/file.pdf", nil)

	// Set the global variables
	plantService = mockPlantService
	pdfGenerator = mockPDFGenerator

	// Create a sample request
	req := models.Request{
		Location: struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		}{
			Latitude:  999.9,
			Longitude: 999.9,
		},
		PlantID: "blueberry",
	}
	reqBody, _ := json.Marshal(req)

	// Call the Lambda handler
	response, err := HandleRequest(context.TODO(), events.APIGatewayProxyRequest{
		Body: string(reqBody),
	})

	// Assert no error occurred
	assert.NoError(t, err)

	// Assert response
	assert.Equal(t, 200, response.StatusCode)
	assert.Contains(t, response.Body, "PDF report generated successfully")

	// Assert the PDF URL in the response
	mockPlantService.AssertCalled(t, "GetPlantInfo", "blueberry")
	mockPDFGenerator.AssertCalled(t, "GeneratePDF", mock.Anything, plantInfo)
	mockPDFGenerator.AssertCalled(t, "UploadToS3", []byte("PDF content"), "plant-report-bucket", "file.pdf")
}

func TestHandleRequest_InvalidRequestBody(t *testing.T) {
	// Call the Lambda handler with an invalid body
	response, err := HandleRequest(context.TODO(), events.APIGatewayProxyRequest{
		Body: "invalid",
	})

	// Assert an error occurred
	assert.NoError(t, err)

	// Assert response
	assert.Equal(t, 400, response.StatusCode)
	assert.Contains(t, response.Body, "Invalid request body")
}

func TestHandleRequest_MissingFields(t *testing.T) {
	mockPlantService := new(mocks.MockPlantService)
	mockPDFGenerator := new(mocks.MockPDFGenerator)

	// Set the global variables
	plantService = mockPlantService
	pdfGenerator = mockPDFGenerator

	// Create a sample request with missing fields
	req := models.Request{
		Location: struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		}{
			Latitude:  0,
			Longitude: 0,
		},
		PlantID: "blueberry",
	}
	reqBody, _ := json.Marshal(req)

	// Call the Lambda handler
	response, err := HandleRequest(context.TODO(), events.APIGatewayProxyRequest{
		Body: string(reqBody),
	})

	// Assert an error occurred
	assert.NoError(t, err)

	// Assert response
	assert.Equal(t, 400, response.StatusCode)
	assert.Contains(t, response.Body, "Missing required fields: plant_id, latitude, and longitude")
}

func TestHandleRequest_FailedPlantInfo(t *testing.T) {
	mockPlantService := new(mocks.MockPlantService)
	mockPDFGenerator := new(mocks.MockPDFGenerator)
	plantInfo := models.PlantInfo{}
	// Mock plant service to return an error
	mockPlantService.On("GetPlantInfo", "1").Return(plantInfo, assert.AnError)

	// Set the global variables
	plantService = mockPlantService
	pdfGenerator = mockPDFGenerator

	// Create a sample request
	req := models.Request{
		Location: struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		}{
			Latitude:  40.7128,
			Longitude: -74.0060,
		},
		PlantID: "1",
	}
	reqBody, _ := json.Marshal(req)

	// Call the Lambda handler
	response, err := HandleRequest(context.TODO(), events.APIGatewayProxyRequest{
		Body: string(reqBody),
	})

	// Assert no error occurred
	assert.NoError(t, err)

	// Assert response
	assert.Equal(t, 500, response.StatusCode)
	assert.Contains(t, response.Body, "Failed to fetch plant information")

	// Assert that the plant service was called
	mockPlantService.AssertCalled(t, "GetPlantInfo", "1")
}
