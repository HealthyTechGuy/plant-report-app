package pdf

import (
	models "github.com/HealthyTechGuy/plant-report-app/models" // Import shared models
)

// PDFGenerator defines the methods for generating PDF reports
type PDFGenerator interface {
	GeneratePDF(userLocation models.UserLocation, plantInfo models.PlantInfo) ([]byte, error)
}
