package pdf

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockPDFGenerator is a mock implementation of PDFGenerator for testing
type MockPDFGenerator struct{}

// GeneratePDF mocks the GeneratePDF method
func (m *MockPDFGenerator) GeneratePDF(plantInfo *PlantInfo) ([]byte, error) {
	// Provide a simple mock implementation
	return []byte("mock PDF content"), nil
}

func TestGeneratePDF_Success(t *testing.T) {
	pdfService := &PDFService{}
	plantInfo := &PlantInfo{
		ID:              "1",
		Name:            "Blueberry Bush",
		GrowingPeriod:   "May to August",
		OptimalPlanting: "Spring",
		HardinessZone:   "3-7",
	}

	pdfBytes, err := pdfService.GeneratePDF(plantInfo)
	require.NoError(t, err)
	assert.NotEmpty(t, pdfBytes)
}

func TestMockGeneratePDF(t *testing.T) {
	mockPDFGenerator := &MockPDFGenerator{}
	plantInfo := &PlantInfo{
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
