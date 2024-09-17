#!/usr/bin/env node
import * as cdk from '@aws-cdk/core';
import { PlantReportStack } from '../lib/plant-report-stack';

const app = new cdk.App();
new PlantReportStack(app, 'PlantReportStack');
