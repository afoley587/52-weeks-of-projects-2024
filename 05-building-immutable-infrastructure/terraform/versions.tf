// versions.tf
terraform {
  required_version = "= 1.6.5"

  # backend "s3" {
  #   key            = "..."
  #   region         = "..."
  #   bucket         = "..."
  #   dynamodb_table = "..."
  #   encrypt        = true
  # }


  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "= 5.33.0"
    }
  }
}