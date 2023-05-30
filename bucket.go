package s3bucket

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

//go:generate <implement mockgen code here>

// Abstraction over interacting with a S3 bucket.
type Bucket interface {
	// Downloads the file from the specified bucket path to the specified host path.
	DownloadFile(bucketFilepath string) (io.ReadCloser, error)

	// Move file inside the bucket from source to target destination.
	MoveFile(sourceDest, targetDest string) error

	// Deletes the target file in the bucket.
	DeleteFile(targetFile string) error

	UploadFile(content io.ReadSeeker, targetDest string) error
}

type bucket struct {
	ctx        context.Context
	session    *session.Session
	downloader *s3manager.Downloader
	uploader   *s3manager.Uploader
	bucketName string
}

type AWSConfig struct {
	Session *session.Session
	Bucket  string
}

func NewS3(
	config AWSConfig,
) Bucket {

	ctx := context.TODO()
	downloader := s3manager.NewDownloader(config.Session)
	uploader := s3manager.NewUploader(config.Session)

	return &bucket{
		ctx:        ctx,
		session:    config.Session,
		downloader: downloader,
		uploader:   uploader,
		bucketName: config.Bucket,
	}
}

func (b bucket) DownloadFile(bucketFilepath string) (io.ReadCloser, error) {
	return nil, nil
	// Implement the DownloadFile method here using b.downloader.
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
