package s3Tools

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// BucketListing is an interface for the AWS API Call List Buckets
type BucketListing interface {
	ListBuckets(ctx context.Context, input *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)
}

// BucketLocation is an interface for the AWS API Call GetBucketLocation
type BucketLocation interface {
	GetBucketLocation(ctx context.Context, input *s3.GetBucketLocationInput, optFns ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error)
}

// BucketVersioning is an interface for the AWS API Call GetBucketVersioning
type BucketVersioning interface {
	GetBucketVersioning(ctx context.Context, input *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error)
}

// BucketEncryption is an interface for the AWS API Call GetBucketEncryption
type BucketEncryption interface {
	GetBucketEncryption(ctx context.Context, input *s3.GetBucketEncryptionInput, optFns ...func(*s3.Options)) (*s3.GetBucketEncryptionOutput, error)
}

// BucketLogging is an interface for the AWS API Call GetBucketLogging
type BucketLogging interface {
	GetBucketLogging(ctx context.Context, input *s3.GetBucketLoggingInput, optFns ...func(*s3.Options)) (*s3.GetBucketLoggingOutput, error)
}

// BucketVisibility is an interface for the AWS API Call GetBucketPolicy
type BucketVisibility interface {
	GetBucketPolicyStatus(ctx context.Context, input *s3.GetBucketPolicyStatusInput, optFns ...func(*s3.Options)) (*s3.GetBucketPolicyStatusOutput, error)
}

// ListBuckets is a function to list all S3 Buckets
func ListBuckets(ctx context.Context, client BucketListing) (*s3.ListBucketsOutput, error) {
	input := &s3.ListBucketsInput{}
	return client.ListBuckets(ctx, input)
}

// GetBucketLocation is a function to return the region in which a S3 bucket exists.
func GetBucketLocation(ctx context.Context, client BucketLocation, bucketName string) (*s3.GetBucketLocationOutput, error) {
	input := &s3.GetBucketLocationInput{
		Bucket: aws.String(bucketName),
	}
	return client.GetBucketLocation(ctx, input)
}

// GetBucketVersioning is a function in which gathers the version of a S3 Bucket
func GetBucketVersioning(ctx context.Context, client BucketVersioning, bucketName string) string {
	input := &s3.GetBucketVersioningInput{
		Bucket: aws.String(bucketName),
	}
	ver, err := client.GetBucketVersioning(ctx, input)
	if err != nil {
		log.Fatalf("failed to get bucket versioning status, %v", err)
	}
	switch ver.Status {
	case "Enabled":
		return "Enabled"
	case "Suspended":
		return "Suspended"
	default:
		return "Not Enabled"
	}
}

// GetBucketEncryption is a function in which gathers the encryption and encryption type of a S3 Bucket
func GetBucketEncryption(ctx context.Context, client BucketEncryption, bucketName string) (string, string) {
	input := &s3.GetBucketEncryptionInput{
		Bucket: aws.String(bucketName),
	}
	enc, err := client.GetBucketEncryption(ctx, input)
	if err != nil {
		return "Not Enabled", "None"
	}
	switch enc.ServerSideEncryptionConfiguration.Rules[0].ApplyServerSideEncryptionByDefault.SSEAlgorithm {
	case "AES256":
		return "Enabled", "SSE"
	case "aws:kms":
		return "Enabled", "KMS"
	default:
		return "Not Enabled", "None"
	}

}

// GetBucketLogging is a function in which gathers the logging of a S3 Bucket
func GetBucketLogging(ctx context.Context, client BucketLogging, bucketName string) (string, string) {
	input := &s3.GetBucketLoggingInput{
		Bucket: aws.String(bucketName),
	}
	logr, err := client.GetBucketLogging(ctx, input)
	if err != nil {
		log.Fatalf("failed to get bucket logging status, %v", err)
	}
	if logr.LoggingEnabled != nil {
		return "Enabled", *logr.LoggingEnabled.TargetBucket
	}
	return "Not Enabled", "None"
}

// GetBucketPolicyStatus is a function that determines whether a S3 Bucket is public or not.
func GetBucketPolicyStatus(ctx context.Context, client BucketVisibility, bucketName string) bool {
	input := &s3.GetBucketPolicyStatusInput{
		Bucket: aws.String(bucketName),
	}
	pol, err := client.GetBucketPolicyStatus(ctx, input)
	if err != nil {
		return false
	}
	return pol.PolicyStatus.IsPublic
}
