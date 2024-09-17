import * as cdk from '@aws-cdk/core';
import * as lambda from '@aws-cdk/aws-lambda';
import * as apigateway from '@aws-cdk/aws-apigateway';
import * as dynamodb from '@aws-cdk/aws-dynamodb';
import * as s3 from '@aws-cdk/aws-s3';

export class PlantReportStack extends cdk.Stack {
    constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        // S3 Bucket for storing PDF reports
        const reportBucket = new s3.Bucket(this, 'PlantReportBucket', {
            versioned: true
        });

        // DynamoDB Table for storing plant information
        const plantTable = new dynamodb.Table(this, 'PlantTable', {
            partitionKey: { name: 'PlantID', type: dynamodb.AttributeType.STRING },
            billingMode: dynamodb.BillingMode.PAY_PER_REQUEST
        });

        // Lambda function for generating plant reports
        const plantReportLambda = new lambda.Function(this, 'PlantReportLambda', {
            runtime: lambda.Runtime.GO_1_X,
            handler: 'main',
            code: lambda.Code.fromAsset('../cmd/plant-report-lambda'),  // Path to the compiled Lambda Go binary
            environment: {
                TABLE_NAME: plantTable.tableName,
                BUCKET_NAME: reportBucket.bucketName
            }
        });

        // Grant the Lambda function permissions to interact with DynamoDB and S3
        plantTable.grantReadData(plantReportLambda);
        reportBucket.grantWrite(plantReportLambda);

        // API Gateway for Lambda invocation
        const api = new apigateway.LambdaRestApi(this, 'PlantReportAPI', {
            handler: plantReportLambda,
            proxy: false
        });

        const plant = api.root.addResource('plant');
        plant.addMethod('POST');  // POST /plant
    }
}
