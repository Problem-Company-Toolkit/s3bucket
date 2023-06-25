package s3bucket

import (
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

//go:generate mockgen -destination=mocks/mock_bucket.go -package=mocks github.com/problem-company-toolkit/s3bucket Bucket

// Abstraction over interacting with a S3 bucket.
type Bucket interface {
	// Downloads the file from the specified bucket path to the specified host path.
	DownloadFile(bucketKey string) (io.ReadCloser, error)

	// Move file inside the bucket from source to target destination.
	MoveFile(sourceDest, targetDest string) error

	// Deletes the target file in the bucket.
	DeleteFile(targetFile string) error

	// Uplodates a file to the configured S3 bucket.
	UploadFile(content io.ReadSeeker, targetDest string) error

	// Gets a signed URL that allows public access to a resource. Expires after the provided duration.
	GetSignedUrl(bucketKey string, duration time.Duration) (string, error)
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
	_, err := b.svc.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(b.bucketName),
		Key:        aws.String(targetDest),
		CopySource: aws.String(fmt.Sprintf("%s/%s", b.bucketName, sourceDest)),
	})

	if err != nil {
		return err
	}

	_, err = b.svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(b.bucketName),
		Key:    aws.String(sourceDest),
	})

	if err != nil {
		return err
	}

	return nil
}

func (b bucket) DeleteFile(targetFile string) error {
	_, err := b.svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(b.bucketName),
		Key:    aws.String(targetFile),
	})

	return err
}

func (b bucket) UploadFile(content io.ReadSeeker, targetDest string) error {
	_, err := b.svc.PutObject(&s3.PutObjectInput{
		Body:   content,
		Bucket: aws.String(b.bucketName),
		Key:    aws.String(targetDest),
	})

	return err
}

func (b bucket) GetSignedUrl(bucketKey string, duration time.Duration) (out string, err error) {
	req, _ := b.svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(b.bucketName),
		Key:    aws.String(bucketKey),
	})

	out, err = req.Presign(duration)
	return out, err
}
