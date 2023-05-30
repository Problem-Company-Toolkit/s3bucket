package s3bucket_test

import (
	"bytes"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/problem-company-toolkit/s3bucket"
)

var _ = Describe("Bucket", func() {
	const (
		FILE_TEST_NAME = "file-test.txt"
	)

	var (
		FILE_TEST_PATH        = os.Getenv("FILE_TEST_PATH")
		AWS_ENDPOINT          = os.Getenv("AWS_ENDPOINT")
		AWS_REGION            = os.Getenv("AWS_DEFAULT_REGION")
		AWS_SECRET_ACCESS_KEY = os.Getenv("AWS_SECRET_ACCESS_KEY")
		AWS_ACCESS_KEY_ID     = os.Getenv("AWS_ACCESS_KEY_ID")
		S3_BUCKET             = os.Getenv("S3_BUCKET")
	)

	var (
		awsSession *session.Session
		svc        *s3.S3
		bucket     s3bucket.Bucket
	)

	BeforeEach(func() {
		var err error
		awsSession, err = session.NewSession(&aws.Config{
			Endpoint:         aws.String(AWS_ENDPOINT),
			Region:           aws.String(AWS_REGION),
			Credentials:      credentials.NewStaticCredentials(AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, ""),
			S3ForcePathStyle: aws.Bool(true),
		})

		if err != nil {
			Fail(err.Error())
			return
		}

		svc = s3.New(awsSession)

		bucket = s3bucket.NewS3(s3bucket.AWSConfig{
			Session: awsSession,
			Bucket:  S3_BUCKET,
		})
	})

	AfterEach(func() {
		_, err := svc.DeleteObject(&s3.DeleteObjectInput{
			Key:    aws.String(FILE_TEST_NAME),
			Bucket: aws.String(S3_BUCKET),
		})

		if err != nil {
			Fail(err.Error())
			return
		}
	})

	Describe("DownloadFile", func() {
		It("Should to download the file from S3", func() {
			testFile, err := os.Open(FILE_TEST_PATH)

			if err != nil {
				Fail(err.Error())
				return
			}
			defer testFile.Close()

			_, err = svc.PutObject(&s3.PutObjectInput{
				Body:   testFile,
				Key:    aws.String(FILE_TEST_NAME),
				Bucket: aws.String(S3_BUCKET),
			})

			if err != nil {
				Fail(err.Error())
				return
			}

			reader, err := bucket.DownloadFile(FILE_TEST_NAME)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(reader).ShouldNot(BeNil())
			defer reader.Close()

			readerBuf := &bytes.Buffer{}
			if _, err := io.Copy(readerBuf, reader); err != nil {
				Fail(err.Error())
				return
			}

			fileBuf := &bytes.Buffer{}
			if _, err := io.Copy(fileBuf, testFile); err != nil {
				Fail(err.Error())
				return
			}

			Expect(readerBuf.String()).Should(Equal(fileBuf.String()))
		})
	})
})
