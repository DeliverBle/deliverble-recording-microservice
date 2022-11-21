package client

import "github.com/aws/aws-sdk-go-v2/service/s3"

type S3Info struct {
	AwsS3Region    string
	AwsAccessKey   string
	AwsSecretKey   string
	AwsProfileName string
	BucketName     string
	S3Client       *s3.Client
}
