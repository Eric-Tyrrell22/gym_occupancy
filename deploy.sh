#!/usr/bin/env bash

GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap main.go

zip lambda.zip bootstrap

#!/bin/bash

FUNCTION_NAME="gym_occupancy"
ROLE_NAME="GymOccupancyLambdaRole"
TRUST_POLICY=$(cat <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
)


# Attempt to get the role
ROLE_ARN=$(aws iam get-role --role-name "$ROLE_NAME" --query 'Role.Arn' --output text 2>/dev/null)

if [[ $? -eq 0 ]]; then
  echo "Role $ROLE_NAME already exists. ARN: $ROLE_ARN"
else
  echo "Role $ROLE_NAME does not exist, creating..."


  # Create the role with the inline trust policy
  ROLE_ARN=$(aws iam create-role --role-name "$ROLE_NAME" --assume-role-policy-document "$TRUST_POLICY" --query 'Role.Arn' --output text)

  # Attach the AWSLambdaBasicExecutionRole policy
  aws iam attach-role-policy --role-name "$ROLE_NAME" --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole
  echo "Role $ROLE_NAME created and policy attached. ARN: $ROLE_ARN"
fi

# Check if the Lambda function exists
if aws lambda get-function --function-name "$FUNCTION_NAME" > /dev/null 2>&1; then
  echo "Lambda function $FUNCTION_NAME already exists, updating..."
  aws lambda update-function-code --function-name "$FUNCTION_NAME" --zip-file fileb://lambda.zip
else
  echo "Lambda function $FUNCTION_NAME does not exist, creating..."

  aws lambda create-function --function-name "$FUNCTION_NAME" \
  --runtime provided.al2023 --handler bootstrap \
  --architectures arm64 \
  --role "$ROLE_ARN" \
  --zip-file fileb://lambda.zip
fi
