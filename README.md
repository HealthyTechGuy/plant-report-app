# Plant Report Generator

A serverless application that generates PDF reports with detailed information about how to grow specific plants, fruits, or vegetables in your chosen area. This tool delivers valuable insights and actionable recommendations aimed at optimizing the growth and success of your chosen flora. Happy Growing!

## Features

- Generates PDF reports for growing plants, fruits, and vegetables.
- Includes information such as growing period, optimal planting times, and hardiness zones.
- Utilizes AWS Lambda for serverless execution, DynamoDB for plant information storage, and S3 for storing the generated PDF files.

## Supported Plants

- Blueberry Bush
- Orange Tree
- Kale

## Architecture

- **AWS Lambda**: Executes the PDF generation logic.
- **DynamoDB**: Stores plant information (e.g., growing hardness zone).
- **S3**: Stores the generated PDF reports.
- **API Gateway**: Exposes the Lambda function via HTTP API.

## Setup

### Prerequisites

- Go 1.18+
- AWS CLI
- AWS Account

## Running project unit tests

I've added a github workflow file which will automatically run the project tests on deployment but if you wish to run them locally you can by using either of these commands:

`go test ./internal/plant-service` & `go test ./pkg/pdf`
`go test ./...`
`make test` 

## Usage

- Use the API Gateway URL to make requests to your Lambda function.
  Example Request:

`curl -X POST https://your-api-id.execute-api.region.amazonaws.com/prod/generate-report \
-H "Content-Type: application/json" \
-d '{"location": "your-location", "plant": "Blueberry Bush"}'`

  
