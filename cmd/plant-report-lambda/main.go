package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	plant "github.com/HealthyTechGuy/plant-report-app/internal/plant-service"
	models "github.com/HealthyTechGuy/plant-report-app/models" // Import shared models
	"github.com/HealthyTechGuy/plant-report-app/pkg/logger"
	"github.com/HealthyTechGuy/plant-report-app/pkg/pdf"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
)

var plantService plant.PlantServiceInterface
var pdfGenerator pdf.PDFGenerator

func init() {
	// Initialize the services (DynamoDB, PDF generator)
	plantService = plant.NewPlantService(os.Getenv("TABLE_NAME"))
	pdfGenerator = &pdf.PDFService{} // Updated to use the concrete implementation
}

// HandleRequest is the main Lambda function handler
func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req models.Request

	logger.InitLogger("debug")
	defer logger.SyncLogger()

	// Unmarshal the request body
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		log.Println("Error unmarshaling request body:", err)
		return responseWithError(400, "Invalid request body"), nil
	}

	// Validate the input (ensure PlantID and Location are provided)
	if req.PlantID == "" || req.Location.Latitude == 0 || req.Location.Longitude == 0 {
		log.Println("Invalid input: missing required fields")
		return responseWithError(400, "Missing required fields: plant_id, latitude, and longitude"), nil
	}

	// Get plant details from the PlantService (DynamoDB)
	plantInfo, err := plantService.GetPlantInfo(req.PlantID)
	if err != nil {
		log.Println("Error fetching plant info:", err)
		return responseWithError(500, "Failed to fetch plant information"), nil
	}

	usrLocation := models.UserLocation{
		UserLatitude:  req.Location.Latitude,
		UserLongitude: req.Location.Longitude,
	}

	// Generate the PDF report using the PDFGenerator
	pdfURL, err := pdfGenerator.GeneratePDF(usrLocation, plantInfo)
	if err != nil {
		log.Println("Error generating PDF report:", err)
		return responseWithError(500, "Failed to generate PDF report"), nil
	}

	bucket := "plant-report-bucket"
	key := "file.pdf"

	s3URL, err := pdfGenerator.UploadToS3(pdfURL, bucket, key)
	if err != nil {
		log.Fatalf("Error uploading to S3: %v", err)
	}

	log.Printf("Uploaded PDF available at: %s", s3URL)

	logger.Info("pdfURL: ", zap.Any("value:", pdfURL))

	// Return the success response with the PDF URL
	return responseWithSuccess(200, ""), nil
}

// responseWithSuccess creates a successful HTTP response with the PDF URL
func responseWithSuccess(statusCode int, pdfURL string) events.APIGatewayProxyResponse {
	response := models.Response{
		Message: "PDF report generated successfully",
		PDFUrl:  pdfURL,
	}
	body, _ := json.Marshal(response)
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       string(body),
	}
}

// responseWithError creates an HTTP error response
func responseWithError(statusCode int, message string) events.APIGatewayProxyResponse {
	response := models.Response{
		Message: message,
	}
	body, _ := json.Marshal(response)
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       string(body),
	}
}

func main() {
	lambda.Start(HandleRequest)
}
