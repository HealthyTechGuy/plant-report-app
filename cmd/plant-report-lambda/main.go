package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	plant "github.com/HealthyTechGuy/plant-report-app/internal/plant-service"
	"github.com/HealthyTechGuy/plant-report-app/pkg/pdf"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var plantService plant.PlantServiceInterface
var pdfGenerator pdf.PDFGenerator

// Response represents the response returned by the Lambda function
type Response struct {
	Message string `json:"message"`
	PDFUrl  string `json:"pdf_url"`
	Error   string `json:"error,omitempty"`
}

// Request represents the expected input from the API Gateway
type Request struct {
	Location struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
	PlantID string `json:"plant_id"`
}

func init() {
	// Initialize the services (DynamoDB, PDF generator)
	plantService = plant.NewPlantService(os.Getenv("TABLE_NAME"))
	pdfGenerator = pdf.NewPDFGenerator(os.Getenv("BUCKET_NAME"))
}

// HandleRequest is the main Lambda function handler
func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req Request

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

	// Generate the PDF report using the PDFGenerator
	pdfURL, err := pdfGenerator.GenerateReport(req.Location.Latitude, req.Location.Longitude, plantInfo)
	if err != nil {
		log.Println("Error generating PDF report:", err)
		return responseWithError(500, "Failed to generate PDF report"), nil
	}

	// Return the success response with the PDF URL
	return responseWithSuccess(200, pdfURL), nil
}

// responseWithSuccess creates a successful HTTP response with the PDF URL
func responseWithSuccess(statusCode int, pdfURL string) events.APIGatewayProxyResponse {
	response := Response{
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
	response := Response{
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
