package main

import (
	"net/http"
	"testing"

	plant "github.com/HealthyTechGuy/plant-report-app/internal/plant-service"
	"github.com/HealthyTechGuy/plant-report-app/pkg/pdf"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockPDFGenerator struct{}

func (m *MockPDFGenerator) GeneratePDF(plantInfo *plant.PlantInfo) ([]byte, error) {
	return []byte("mock PDF content"), nil
}

func TestHandler_Success(t *testing.T) {
	// Arrange
	pdfGenerator = &MockPDFGenerator{}
	request := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"location": "your-location",
			"plant":    "Blueberry Bush",
		},
	}

	// Act
	response, err := handler(request)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Contains(t, response.Body, "PDF generated successfully")
}

func TestHandler_MissingParameters(t *testing.T) {
	// Arrange
	request := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{},
	}

	// Act
	response, err := handler(request)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "Missing location or plant name", response.Body)
}

func TestHandler_PDFGenerationError(t *testing.T) {
	// Arrange
	pdfGenerator = &pdf.ErrorGeneratingPDF{} // Assuming this is a custom type that returns an error
	request := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"location": "your-location",
			"plant":    "Blueberry Bush",
		},
	}

	// Act
	response, err := handler(request)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
	assert.Equal(t, "Error generating PDF", response.Body)
}
