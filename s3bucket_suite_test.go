package s3bucket_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestS3bucket(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "S3bucket Suite")
}
