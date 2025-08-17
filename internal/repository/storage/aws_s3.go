package storage

import (
	"context"
	"io"
)

type AwsS3Impl struct {
}

func NewAwsS3Impl() *AwsS3Impl {
	return &AwsS3Impl{}
}

func (a AwsS3Impl) Upload(ctx context.Context, key string, r io.Reader) (string, error) {
	return key, nil
}

func (a AwsS3Impl) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	//TODO implement me
	panic("implement me")
}

func (a AwsS3Impl) Delete(ctx context.Context, key string) error {
	//TODO implement me
	panic("implement me")
}
