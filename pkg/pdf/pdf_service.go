package pdf

import (
	"bytes"
	"log"

	"github.com/jung-kurt/gofpdf"
)

// PDFService is a concrete implementation of the PDFGenerator interface
type PDFService struct{}

// PlantInfo holds information about a plant
type PlantInfo struct {
	ID              string
	Name            string
	GrowingPeriod   string
	OptimalPlanting string
	HardinessZone   string
}

// GeneratePDF creates a PDF report for a given plant information
func (s *PDFService) GeneratePDF(plantInfo *PlantInfo) ([]byte, error) {
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
