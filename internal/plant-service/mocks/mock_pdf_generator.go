// mocks/mock_pdf_generator.go
package mocks

import (
	"github.com/HealthyTechGuy/plant-report-app/models"
	"github.com/stretchr/testify/mock"
)

type MockPDFGenerator struct {
	mock.Mock
}

// GeneratePDF generates a mock PDF and returns it as a byte slice
func (m *MockPDFGenerator) GeneratePDF(location models.UserLocation, plantInfo models.PlantInfo) ([]byte, error) {
	args := m.Called(location, plantInfo)
	return args.Get(0).([]byte), args.Error(1)
}
