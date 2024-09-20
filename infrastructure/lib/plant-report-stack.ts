import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as apigateway from 'aws-cdk-lib/aws-apigateway';
import * as s3 from 'aws-cdk-lib/aws-s3';
import * as iam from 'aws-cdk-lib/aws-iam';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';  // Import DynamoDB

export class PlantReportStack extends cdk.Stack {
    constructor(scope: Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        // Create S3 bucket (for PDFs, etc.)
        const reportBucket = new s3.Bucket(this, 'PlantReportBucket', {
            bucketName: 'plant-report-bucket',
            removalPolicy: cdk.RemovalPolicy.DESTROY,
            versioned: true,
            publicReadAccess: false,
        });

        // Create DynamoDB table
        const plantReportTable = new dynamodb.Table(this, 'PlantReportTable', {
            tableName: 'plant-report-app',
            partitionKey: { name: 'PlantID', type: dynamodb.AttributeType.STRING },
            billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,  // Adjust based on your needs
            removalPolicy: cdk.RemovalPolicy.DESTROY,  // Remove the table when the stack is deleted
        });

        // Define Lambda function for handling the requests
        const plantReportLambda = new lambda.Function(this, 'PlantReportLambda', {
            runtime: lambda.Runtime.PROVIDED_AL2,
            code: lambda.Code.fromAsset('../dist/plant-report-lambda.zip'),
            handler: 'bootstrap',  // Dummy value, not used by custom runtime
            environment: {
                TABLE_NAME: plantReportTable.tableName,  // Use table name from the table resource
                BUCKET_NAME: reportBucket.bucketName,
            },
        });

        // Create IAM policy for DynamoDB access
        const dynamoPolicy = new iam.PolicyStatement({
            actions: ['dynamodb:GetItem'],
            resources: [
                plantReportTable.tableArn  // Use the ARN of the table created in the stack
            ],
        });

        // Attach the policy to the Lambda function's role
        plantReportLambda.addToRolePolicy(dynamoPolicy);

        // Define API Gateway to trigger the Lambda
        const api = new apigateway.LambdaRestApi(this, 'PlantReportApi', {
            handler: plantReportLambda,
            proxy: false,
        });

        const plantResource = api.root.addResource('report');
        plantResource.addMethod('POST');
    }
}
