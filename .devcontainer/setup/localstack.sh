#!/bin/bash

# Create S3 buckets
awslocal s3 mb "s3://${S3_BUCKET}" --region "${AWS_DEFAULT_REGION}"
awslocal s3 ls
