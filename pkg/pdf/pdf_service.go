package pdf

import (
	"bytes"
	"log"

	models "github.com/HealthyTechGuy/plant-report-app/models" // Import shared models
	"github.com/jung-kurt/gofpdf"
)

// PDFService is a concrete implementation of the PDFGenerator interface
type PDFService struct{}

// GeneratePDF creates a PDF report for a given plant information
func (s *PDFService) GeneratePDF(userLocation models.UserLocation, plantInfo models.PlantInfo) ([]byte, error) {
	// Create a new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")

	// Add a new page
	pdf.AddPage()

	// Set font
	pdf.SetFont("Arial", "B", 16)

	// Add a title
	pdf.Cell(40, 10, "Plant Report")

	// Add plant information
	pdf.Ln(10) // Line break
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, "Plant Name: "+plantInfo.Name)
	pdf.Cell(40, 10, "Growing Period: "+plantInfo.GrowingPeriod)
	pdf.Cell(40, 10, "Optimal Planting Time: "+plantInfo.OptimalPlanting)
	pdf.Cell(40, 10, "Hardiness Zone: "+plantInfo.HardinessZone)

	// Output the PDF to a buffer
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		log.Printf("Error generating PDF: %v", err)
		return nil, err
	}

	return buf.Bytes(), nil
}
