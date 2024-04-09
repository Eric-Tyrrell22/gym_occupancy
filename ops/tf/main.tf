
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.44"
    }
  }

  required_version = ">= 1.2.0"
}

provider "aws" {
  region  = "us-east-1"
}

data "aws_iam_policy_document" "assume_role" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "iam_for_lambda" {
  name               = "iam_for_lambda"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

resource "aws_iam_role_policy_attachment" "test-attach" {
  role       = aws_iam_role.iam_for_lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

data "archive_file" "gym_occupancy_lambda_zip" {
  type        = "zip"
  source_file = "../../bootstrap"
  output_path = "../../lambda.zip"
}

resource "aws_lambda_function" "gym_occupancy_lambda" {
  filename      = data.archive_file.gym_occupancy_lambda_zip.output_path
  function_name = "gym_occupancy"
  role          = aws_iam_role.iam_for_lambda.arn
  handler       = "bootstrap"

  source_code_hash = data.archive_file.gym_occupancy_lambda_zip.output_base64sha256

  architectures = ["arm64"]
  runtime = "provided.al2023"
}
