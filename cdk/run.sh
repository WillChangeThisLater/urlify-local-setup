#!/bin/bash

# since S3 bucket names are globally unique, 
# this will need to be changed everytime
BUCKETNAME="urlify"

# make sure the user has some sort of AWS creds installed
if ! aws sts get-caller-identity >/dev/null 2>&1; then
	echo "AWS Credentials not found"
	exit 1
fi

terraform apply -var="bucket_name=$BUCKETNAME"
