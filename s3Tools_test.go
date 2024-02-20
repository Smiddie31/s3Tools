package s3Tools

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"testing"
)

type mockBucketVersioning func(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error)

func (m mockBucketVersioning) GetBucketVersioning(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
	return m(ctx, params, optFns...)
}

func TestGetBucketVersioning(t *testing.T) {
	mockResponse := &s3.GetBucketVersioningOutput{
		Status: "Enabled",
	}

	mockClient := mockBucketVersioning(func(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
		return mockResponse, nil
	})

	resp := GetBucketVersioning(context.Background(), mockClient, "test-bucket")

	if resp != "Enabled" {
		t.Errorf("unexpected status: got %s, want %s", resp, "Enabled")
	}
}
