terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      version = "~> 3.0"
    }
  }
}

provider "aws" {
  region = "eu-central-1"
  shared_credentials_file = "/Users/mobile32/.aws/credentials"
}

resource "aws_iam_role" "iam-for-lambda" {
  name = "iam_for_lambda"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": "AllowS3Access"
    }
  ]
}
EOF
}

resource "aws_iam_policy" "lambda-policy" {
  name = "lambda-policy"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:*"
      ],
      "Resource": [
        "${aws_s3_bucket.photos-bucket.arn}",
        "${aws_s3_bucket.photos-bucket.arn}/*"
      ]
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "lambda-policy-attach" {
  role = aws_iam_role.iam-for-lambda.name
  policy_arn = aws_iam_policy.lambda-policy.arn
}

resource "aws_iam_policy" "lambda-logging" {
  name        = "lambda_logging"
  path        = "/"
  description = "IAM policy for logging from a lambda"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "arn:aws:logs:*:*:*",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "lambda-logs-attach" {
  role       = aws_iam_role.iam-for-lambda.name
  policy_arn = aws_iam_policy.lambda-logging.arn
}

resource "aws_cloudwatch_log_group" "lambda-log-group" {
  name              = "/aws/lambda/scaler"
  retention_in_days = 14
}

resource "aws_lambda_function" "lambda-function" {
  filename = "../src/main.zip"
  function_name = "scaler"
  role = aws_iam_role.iam-for-lambda.arn
  handler = "main"
  source_code_hash = filebase64sha256("../src/main.zip")
  runtime = "go1.x"
  timeout = 300
  memory_size = 256

  environment {
    variables = {
      BUCKET_NAME = aws_s3_bucket.photos-bucket.bucket
      SOURCE_PATH = "originals"
      TARGET_PATH = "thumbs"
    }
  }

  depends_on = [
    aws_iam_role_policy_attachment.lambda-policy-attach,
    aws_iam_role_policy_attachment.lambda-logs-attach,

    aws_cloudwatch_log_group.lambda-log-group,
  ]
}

resource "aws_s3_bucket" "photos-bucket" {
  bucket = "scaling-photos-bucket"
}

resource "aws_lambda_permission" "allow-bucket" {
  statement_id  = "AllowExecutionFromS3Bucket"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.lambda-function.arn
  principal     = "s3.amazonaws.com"
  source_arn    = aws_s3_bucket.photos-bucket.arn
}

resource "aws_s3_bucket_notification" "bucket-notification" {
  bucket = aws_s3_bucket.photos-bucket.id

  lambda_function {
    lambda_function_arn = aws_lambda_function.lambda-function.arn
    events              = ["s3:ObjectCreated:*"]
    filter_prefix       = "originals/"
  }

  depends_on = [aws_lambda_permission.allow-bucket]
}

resource "aws_s3_bucket_object" "images" {
  for_each = fileset("../test_images/", "*")
  bucket = aws_s3_bucket.photos-bucket.id
  key = "originals/${each.value}"
  source = "../test_images/${each.value}"
  etag = filemd5("../test_images/${each.value}")
}
