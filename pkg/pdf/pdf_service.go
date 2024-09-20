package pdf

import (
	"bytes"
	"fmt"
	"log"

	models "github.com/HealthyTechGuy/plant-report-app/models" // Import shared models
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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

func (s *PDFService) UploadToS3(data []byte, bucket, key string) (string, error) {
	sess := session.Must(session.NewSession())
	svc := s3.New(sess)

	// Create an S3 PutObjectInput with the byte data
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
