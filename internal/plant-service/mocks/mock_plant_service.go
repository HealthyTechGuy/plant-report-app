// mocks/mock_plant_service.go
package mocks

import (
	"github.com/HealthyTechGuy/plant-report-app/models"
	"github.com/stretchr/testify/mock"
)

type MockPlantService struct {
	mock.Mock
}

func (m *MockPlantService) GetPlantInfo(plantID string) (models.PlantInfo, error) {
	args := m.Called(plantID)
	return args.Get(0).(models.PlantInfo), args.Error(1)
}
