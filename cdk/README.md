# Deploying the S3 Bucket for the Urlify Tool

## Requirements
- You have installed [Terraform](https://learn.hashicorp.com/tutorials/terraform/install-cli) and [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html).
- You have the necessary [AWS credentials setup](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-note1.html).

## Setup
1. Clone this repository
2. In the directory, adjust the region in the `provider` block in `main.tf` if you wish to create the bucket in a different region.
3. Update the `bucket` parameter in the `aws_s3_bucket` resource in `main.tf` to a unique name of your choosing.

## Deployment
To deploy the S3 bucket, navigate to the directory containing this README and `main.tf` and run the following commands:
```bash
terraform init
terraform apply -var='bucket_name=BUCKET_NAME_HERE'

