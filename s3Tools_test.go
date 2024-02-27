package s3Tools

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"testing"
)

type mockBucketVersioning func(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error)

func (m mockBucketVersioning) GetBucketVersioning(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
	return m(ctx, params, optFns...)
}

type MockBucketEncryption struct{}

func (m *MockBucketEncryption) GetBucketEncryption(ctx context.Context, params *s3.GetBucketEncryptionInput, optFns ...func(*s3.Options)) (*s3.GetBucketEncryptionOutput, error) {
	// Mock response
	return &s3.GetBucketEncryptionOutput{
		ServerSideEncryptionConfiguration: &types.ServerSideEncryptionConfiguration{
			Rules: []types.ServerSideEncryptionRule{
				{
					ApplyServerSideEncryptionByDefault: &types.ServerSideEncryptionByDefault{
						SSEAlgorithm: "AES256",
					},
				},
			},
		},
	}, nil
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

func TestGetBucketEncryption(t *testing.T) {
	mockClient := &MockBucketEncryption{}
	resp1, resp2 := GetBucketEncryption(context.Background(), mockClient, "test-bucket")
	if resp1 != "Enabled" || resp2 != "SSE" {
		t.Errorf("unexpected status: got %s %s, want %s %s", resp1, resp2, "Enabled", "SSE")
	}
}
