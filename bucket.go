package s3bucket

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

//go:generate <implement mockgen code here>

// Abstraction over interacting with a S3 bucket.
type Bucket interface {
	// Downloads the file from the specified bucket path to the specified host path.
	DownloadFile(bucketKey string) (io.ReadCloser, error)

	// Move file inside the bucket from source to target destination.
	MoveFile(sourceDest, targetDest string) error

	// Deletes the target file in the bucket.
	DeleteFile(targetFile string) error

	UploadFile(content io.ReadSeeker, targetDest string) error
}

type bucket struct {
	svc        *s3.S3
	bucketName string
}

type AWSConfig struct {
	Session *session.Session
	Bucket  string
}

func NewS3(
	config AWSConfig,
) Bucket {
	return &bucket{
		svc:        s3.New(config.Session),
		bucketName: config.Bucket,
	}
}

func (b bucket) DownloadFile(bucketKey string) (io.ReadCloser, error) {
	response, err := b.svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(b.bucketName),
		Key:    aws.String(bucketKey),
	})

	if err != nil {
		return nil, err
	}

	return response.Body, nil
}

func (b bucket) MoveFile(sourceDest, targetDest string) error {
	return nil
	// Implement the MoveFile method here using b.uploader and b.downloader.
}

func (b bucket) DeleteFile(targetFile string) error {
	return nil
	// Implement the DeleteFile method here using b.session.
}

func (b bucket) UploadFile(content io.Writer, targetDest string) error {
	return nil
}
