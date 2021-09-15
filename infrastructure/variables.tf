variable "aws_region" {
  type        = string
  description = "AWS region"
}

variable "aws_credentials_path" {
  type        = string
  description = "AWS credentials path"
}

variable "bucket_name" {
  type        = string
  description = "Name of photos bucket"
}

variable "function_name" {
  type        = string
  description = "Name of lambda function"
}

variable "source_path" {
  type        = string
  description = "Source path in bucket"
}

variable "target_path" {
  type        = string
  description = "Target path in bucket"
}

variable "image_width" {
  type        = number
  description = "Default scaled images width"
}

variable "image_height" {
  type        = number
  description = "Default scaled images height"
}

variable "test_images_path" {
  type        = string
  description = "Test images path"
}