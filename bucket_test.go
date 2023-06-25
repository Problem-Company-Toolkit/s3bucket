package s3bucket_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

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
		FILE_TEST_NAME       = "file-test.txt"
		MOVED_FILE_TEST_NAME = "moved-file.txt"
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
		It("Downloads the file from S3", func() {
			sendFile, err := os.Open(FILE_TEST_PATH)

			if err != nil {
				Fail(err.Error())
				return
			}
			defer sendFile.Close()

			_, err = svc.PutObject(&s3.PutObjectInput{
				Body:   sendFile,
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

			testFile, err := os.Open(FILE_TEST_PATH)

			if err != nil {
				Fail(err.Error())
				return
			}
			defer testFile.Close()

			fileBuf := &bytes.Buffer{}
			if _, err := io.Copy(fileBuf, testFile); err != nil {
				Fail(err.Error())
				return
			}

			Expect(readerBuf.String()).To(Equal(fileBuf.String()))
		})
	})

	Describe("MoveFile", func() {
		It("Moves a file in s3", func() {
			sendFile, err := os.Open(FILE_TEST_PATH)

			if err != nil {
				Fail(err.Error())
				return
			}
			defer sendFile.Close()

			_, err = svc.PutObject(&s3.PutObjectInput{
				Body:   sendFile,
				Key:    aws.String(FILE_TEST_NAME),
				Bucket: aws.String(S3_BUCKET),
			})

			if err != nil {
				Fail(err.Error())
				return
			}

			err = bucket.MoveFile(FILE_TEST_NAME, MOVED_FILE_TEST_NAME)

			Expect(err).ShouldNot(HaveOccurred())

			_, err = svc.GetObject(&s3.GetObjectInput{
				Bucket: aws.String(S3_BUCKET),
				Key:    aws.String(FILE_TEST_NAME),
			})

			Expect(err).Should(MatchError(ContainSubstring(s3.ErrCodeNoSuchKey)))

			obj, err := svc.GetObject(&s3.GetObjectInput{
				Bucket: aws.String(S3_BUCKET),
				Key:    aws.String(MOVED_FILE_TEST_NAME),
			})

			Expect(err).ShouldNot(HaveOccurred())
			defer obj.Body.Close()

			testFile, err := os.Open(FILE_TEST_PATH)

			if err != nil {
				Fail(err.Error())
				return
			}
			defer testFile.Close()

			fileBuf := &bytes.Buffer{}
			if _, err := io.Copy(fileBuf, testFile); err != nil {
				Fail(err.Error())
				return
			}

			objBuf := &bytes.Buffer{}
			if _, err := io.Copy(objBuf, obj.Body); err != nil {
				Fail(err.Error())
				return
			}

			Expect(fileBuf.String()).To(Equal(objBuf.String()))
		})
	})

	Describe("DeleteFile", func() {
		It("Deletes a file from s3", func() {
			sendFile, err := os.Open(FILE_TEST_PATH)

			if err != nil {
				Fail(err.Error())
				return
			}
			defer sendFile.Close()

			_, err = svc.PutObject(&s3.PutObjectInput{
				Body:   sendFile,
				Key:    aws.String(FILE_TEST_NAME),
				Bucket: aws.String(S3_BUCKET),
			})

			if err != nil {
				Fail(err.Error())
				return
			}

			err = bucket.DeleteFile(FILE_TEST_NAME)

			Expect(err).ShouldNot(HaveOccurred())

			_, err = svc.GetObject(&s3.GetObjectInput{
				Bucket: aws.String(S3_BUCKET),
				Key:    aws.String(FILE_TEST_NAME),
			})

			Expect(err).Should(MatchError(ContainSubstring(s3.ErrCodeNoSuchKey)))
		})
	})

	Describe("UploadFile", func() {
		It("Uploads a file to s3", func() {
			sendFile, err := os.Open(FILE_TEST_PATH)

			if err != nil {
				Fail(err.Error())
				return
			}
			defer sendFile.Close()

			err = bucket.UploadFile(sendFile, FILE_TEST_NAME)

			Expect(err).ShouldNot(HaveOccurred())

			obj, err := svc.GetObject(&s3.GetObjectInput{
				Bucket: aws.String(S3_BUCKET),
				Key:    aws.String(FILE_TEST_NAME),
			})

			Expect(err).ShouldNot(HaveOccurred())
			defer obj.Body.Close()

			testFile, err := os.Open(FILE_TEST_PATH)

			if err != nil {
				Fail(err.Error())
				return
			}
			defer testFile.Close()

			fileBuf := &bytes.Buffer{}
			if _, err := io.Copy(fileBuf, testFile); err != nil {
				Fail(err.Error())
				return
			}

			objBuf := &bytes.Buffer{}
			if _, err := io.Copy(objBuf, obj.Body); err != nil {
				Fail(err.Error())
				return
			}

			Expect(fileBuf.String()).To(Equal(objBuf.String()))
		})
	})

	Describe("GetSignedUrl", func() {
		It("Generates a signed URL for an object", func() {
			signed, err := bucket.GetSignedUrl(FILE_TEST_NAME, time.Minute*5)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(signed).ShouldNot(BeEmpty())

			fmt.Printf("signed url: %s\n", signed)
		})
	})
})
