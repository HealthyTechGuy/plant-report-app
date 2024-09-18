#!/usr/bin/env node
import * as cdk from 'aws-cdk-lib';
import { PlantReportStack } from '../lib/plant-report-stack';
import { App } from 'aws-cdk-lib';  // Ensure 'App' is imported from 'aws-cdk-lib'

const app = new App();
new PlantReportStack(app, 'PlantReportStack');
