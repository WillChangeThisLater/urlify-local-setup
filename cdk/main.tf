variable "bucket_name" {}

provider "aws" {
  region = "us-east-2"
}

resource "aws_s3_bucket" "urlify_bucket" {
  bucket = var.bucket_name
}

resource "aws_s3_bucket_lifecycle_configuration" "urlify_lifecycle" {
  bucket = aws_s3_bucket.urlify_bucket.id

  rule {
    id      = "Lifecycle"
    status  = "Enabled"

    expiration {
      days = 1
    }
  }
}
