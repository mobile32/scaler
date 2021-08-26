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

resource "aws_iam_policy" "policy" {
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
        "${aws_s3_bucket.original-photos-bucket.arn}",
        "${aws_s3_bucket.original-photos-bucket.arn}/*",
        "${aws_s3_bucket.resized-photos-bucket.arn}",
        "${aws_s3_bucket.resized-photos-bucket.arn}/*"
      ]
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "lambda-policy-attach" {
  role = aws_iam_role.iam-for-lambda.name
  policy_arn = aws_iam_policy.policy.arn
}

//resource "null_resource" "build" {
//  provisioner "local-exec" {
//    interpreter = [
//      "/bin/bash",
//      "-c"]
//
//    command = <<-EOT
//      exec "GOOS=linux GOARCH=amd64 go build -o main ../src/main.go"
//      exec "zip main.zip main"
//    EOT
//  }
//}

resource "aws_lambda_function" "lambda-function" {
  filename = "../src/main.zip"
  function_name = "scaler"
  role = aws_iam_role.iam-for-lambda.arn
  handler = "main"
  source_code_hash = filebase64sha256("../src/main.zip")
  runtime = "go1.x"
  timeout = 300

  environment {
    variables = {
      originalBucketName = "owl-original-photos-bucket"
      resizedBucketName = "owl-resized-photos-bucket"
      tmpPath = "tmp_path"
    }
  }

  depends_on = [
    aws_iam_role_policy_attachment.lambda-policy-attach,
  ]
}

resource "aws_s3_bucket" "original-photos-bucket" {
  bucket = "owl-original-photos-bucket"
  acl = "public-read"

  tags = {
    Name = "Owl original photos"
    Environment = "Dev"
  }
}

resource "aws_s3_bucket" "resized-photos-bucket" {
  bucket = "owl-resized-photos-bucket"
  acl = "public-read"

  tags = {
    Name = "Owl resized photos"
    Environment = "Dev"
  }
}

resource "aws_s3_bucket_object" "images" {
  for_each = fileset("../test_images/", "*")
  bucket = aws_s3_bucket.original-photos-bucket.id
  key = each.value
  source = "../test_images/${each.value}"
  etag = filemd5("../test_images/${each.value}")
}
