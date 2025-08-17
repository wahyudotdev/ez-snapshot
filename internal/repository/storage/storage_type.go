package storage

type StorageType int

const (
	AwsS3 StorageType = 0
	Gcs   StorageType = 1
)
