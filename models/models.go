package models

// PlantInfo holds information about a plant
type PlantInfo struct {
	ID              string
	Name            string
	GrowingPeriod   string
	OptimalPlanting string
	HardinessZone   string
}

type UserLocation struct {
	UserLatitude  float64 `json:"latitude"`
	UserLongitude float64 `json:"longitude"`
}

// Request represents the expected input from the API Gateway
type Request struct {
	Location struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
	PlantID string `json:"plant_id"`
}

// Response represents the response returned by the Lambda function
type Response struct {
	Message string `json:"message"`
	PDFUrl  string `json:"pdf_url"`
	Error   string `json:"error,omitempty"`
}
