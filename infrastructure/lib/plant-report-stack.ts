import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';  // Construct is now imported from 'constructs'
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as apigateway from 'aws-cdk-lib/aws-apigateway';
import * as s3 from 'aws-cdk-lib/aws-s3';

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

        // Define Lambda function for handling the requests
        const plantReportLambda = new lambda.Function(this, 'PlantReportLambda', {
            runtime: lambda.Runtime.PROVIDED_AL2,  // Updated to go1.20 or latest
            code: lambda.Code.fromAsset('../dist/plant-report-lambda.zip'),
            handler: 'bootstrap',  // Dummy value, not used by custom runtime
            environment: {
                TABLE_NAME: 'dynamodb-table-name',
                BUCKET_NAME: reportBucket.bucketName,
            },
        });

        // Define API Gateway to trigger the Lambda
        const api = new apigateway.LambdaRestApi(this, 'PlantReportApi', {
            handler: plantReportLambda,
            proxy: false,
        });

        const plantResource = api.root.addResource('report');
        plantResource.addMethod('POST');
    }
}
