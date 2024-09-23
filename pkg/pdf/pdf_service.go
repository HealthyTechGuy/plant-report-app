package pdf

import (
	"bytes"
	"fmt"
	"log"

	models "github.com/HealthyTechGuy/plant-report-app/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jung-kurt/gofpdf"
)

// PDFService is a concrete implementation of the PDFGenerator interface
type PDFService struct{}

// GeneratePDF creates a nicely formatted PDF report for given plant information
func (s *PDFService) GeneratePDF(userLocation models.UserLocation, plantInfo models.PlantInfo) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")

	// Set margins
	pdf.SetMargins(15, 10, 15)

	// Add a new page
	pdf.AddPage()

	// Header
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 30, "Plant Growth Report") // Adjust y-position to fit after the image
	pdf.Ln(25)

	// Sub-header (Plant name, date, location)
	pdf.SetFont("Arial", "I", 12)
	pdf.Cell(190, 8, fmt.Sprintf("Report for: %s", plantInfo.Name))
	pdf.Ln(8)
	pdf.Cell(190, 8, fmt.Sprintf("Location latitude: %.6f", userLocation.UserLatitude)) // Corrected precision formatting for lat/long
	pdf.Ln(8)
	pdf.Cell(190, 8, fmt.Sprintf("Location longitude: %.6f", userLocation.UserLongitude))
	pdf.Ln(12)

	// Plant Information Section
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(190, 10, "Plant Information")
	pdf.Ln(10)

	// Table-style plant info
	pdf.SetFont("Arial", "", 11)
	pdf.Cell(50, 10, "Plant Name:")
	pdf.Cell(100, 10, plantInfo.Name)
	pdf.Ln(8)
	pdf.Cell(50, 10, "Growing Period:")
	pdf.Cell(100, 10, plantInfo.GrowingPeriod)
	pdf.Ln(8)
	pdf.Cell(50, 10, "Optimal Planting Time:")
	pdf.Cell(100, 10, plantInfo.OptimalPlanting)
	pdf.Ln(8)
	pdf.Cell(50, 10, "Hardiness Zone:")
	pdf.Cell(100, 10, plantInfo.HardinessZone)
	pdf.Ln(10)

	// Footer
	pdf.SetY(-15)
	pdf.SetFont("Arial", "I", 8)
	pdf.Cell(0, 10, fmt.Sprintf("Page %d", pdf.PageNo()))

	// Output the PDF to a buffer
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		log.Printf("Error generating PDF: %v", err)
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *PDFService) UploadToS3(data []byte, bucket, key string) (string, error) {
	sess := session.Must(session.NewSession())
	svc := s3.New(sess)

	_, err := svc.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String("application/pdf"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	s3URL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, key)
	log.Printf("File uploaded to: %s", s3URL)
	return s3URL, nil
}
