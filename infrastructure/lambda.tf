terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.0"
    }
  }
}

provider "aws" {
  region                  = "eu-central-1"
  shared_credentials_file = "/Users/mobile32/.aws/credentials"
}

resource "aws_iam_role" "iam_for_lambda" {
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
      "Sid": ""
    }
  ]
}
EOF
}

# resource "aws_lambda_function" "test_lambda" {
#   filename      = "../main.zip"
#   function_name = "main"
#   role          = aws_iam_role.iam_for_lambda.arn
#   handler       = "main"

#   source_code_hash = filebase64sha256("../main.zip")

#   runtime = "go1.x"

#   environment {
#     variables = {
#     }
#   }
# }

resource "aws_s3_bucket" "original-photos-bucket" {
  bucket = "owl-original-photos-bucket"
  acl    = "public-read"

  tags = {
    Name        = "Owl original photos"
    Environment = "Dev"
  }
}

resource "aws_s3_bucket" "owl-resized-photos-bucket" {
  bucket = "owl-resized-photos-bucket"
  acl    = "public-read"

  tags = {
    Name        = "Owl resized photos"
    Environment = "Dev"
  }
}

resource "aws_s3_bucket_object" "images" {
  for_each = fileset("../test_images/", "*")
  bucket   = aws_s3_bucket.original-photos-bucket.id
  key      = each.value
  source   = "../test_images/${each.value}"
  etag     = filemd5("../test_images/${each.value}")
}
