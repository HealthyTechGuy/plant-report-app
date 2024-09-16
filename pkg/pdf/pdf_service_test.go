package pdf

import (
	"testing"

	models "github.com/HealthyTechGuy/plant-report-app/models" // Import shared models
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockPDFGenerator is a mock implementation of PDFGenerator for testing
type MockPDFGenerator struct{}

// GeneratePDF mocks the GeneratePDF method
func (m *MockPDFGenerator) GeneratePDF(plantInfo models.PlantInfo) ([]byte, error) {
	// Provide a simple mock implementation
	return []byte("mock PDF content"), nil
}

func TestGeneratePDF_Success(t *testing.T) {
	pdfService := &PDFService{}
	plantInfo := models.PlantInfo{
		ID:              "1",
		Name:            "Blueberry Bush",
		GrowingPeriod:   "May to August",
		OptimalPlanting: "Spring",
		HardinessZone:   "3-7",
	}
	userLocation := models.UserLocation{
		UserLatitude:  999.999,
		UserLongitude: 999.99,
	}

	pdfBytes, err := pdfService.GeneratePDF(userLocation, plantInfo)
	require.NoError(t, err)
	assert.NotEmpty(t, pdfBytes)
}

func TestMockGeneratePDF(t *testing.T) {
	mockPDFGenerator := &MockPDFGenerator{}
	plantInfo := models.PlantInfo{
		ID:              "1",
		Name:            "Blueberry Bush",
		GrowingPeriod:   "May to August",
		OptimalPlanting: "Spring",
		HardinessZone:   "3-7",
	}

	pdfBytes, err := mockPDFGenerator.GeneratePDF(plantInfo)
	require.NoError(t, err)
	assert.Equal(t, "mock PDF content", string(pdfBytes))
}
