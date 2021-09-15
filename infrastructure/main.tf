resource "aws_iam_role" "iam-for-lambda" {
  name = "iam-for-lambda-${var.function_name}"

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
  name = "lambda-policy-${var.function_name}"

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
  name = "lambda-logging-${var.function_name}"
  path = "/"
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
  role = aws_iam_role.iam-for-lambda.name
  policy_arn = aws_iam_policy.lambda-logging.arn
}

resource "aws_cloudwatch_log_group" "lambda-log-group" {
  name = "/aws/lambda/${var.function_name}"
  retention_in_days = 14
}

data "archive_file" "dummy" {
  type = "zip"
  output_path = "${path.module}/payload.zip"

  source {
    content = "empty lambda function"
    filename = "payload.txt"
  }
}

resource "aws_lambda_function" "lambda-function" {
  filename = data.archive_file.dummy.output_path
  function_name = var.function_name
  role = aws_iam_role.iam-for-lambda.arn
  handler = "main"
  runtime = "go1.x"
  timeout = 300
  memory_size = 256

  environment {
    variables = {
      BUCKET_NAME = aws_s3_bucket.photos-bucket.bucket
      SOURCE_PATH = var.source_path
      TARGET_PATH = var.target_path
      IMAGES_WIDTH = var.image_width
      IMAGES_HEIGHT = var.image_height
    }
  }

  depends_on = [
    aws_iam_role_policy_attachment.lambda-policy-attach,
    aws_iam_role_policy_attachment.lambda-logs-attach,

    aws_cloudwatch_log_group.lambda-log-group,
  ]
}

resource "aws_s3_bucket" "photos-bucket" {
  bucket = var.bucket_name
}

resource "aws_lambda_permission" "allow-bucket" {
  statement_id = "AllowExecutionFromS3Bucket"
  action = "lambda:InvokeFunction"
  function_name = aws_lambda_function.lambda-function.arn
  principal = "s3.amazonaws.com"
  source_arn = aws_s3_bucket.photos-bucket.arn
}

resource "aws_s3_bucket_notification" "bucket-notification" {
  bucket = aws_s3_bucket.photos-bucket.id

  lambda_function {
    lambda_function_arn = aws_lambda_function.lambda-function.arn
    events = [
      "s3:ObjectCreated:*"]
    filter_prefix = "${var.source_path}/"
  }

  depends_on = [
    aws_lambda_permission.allow-bucket]
}

resource "aws_s3_bucket_object" "images" {
  for_each = fileset("${var.test_images_path}/", "*")
  bucket = aws_s3_bucket.photos-bucket.id
  key = "${var.source_path}/${each.value}"
  source = "${var.test_images_path}/${each.value}"
  etag = filemd5("${var.test_images_path}/${each.value}")
}
