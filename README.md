# S3Bucket Package

The `s3bucket` package is an open-source Go package that provides an abstraction for interacting with an S3 bucket. It simplifies common operations such as downloading files, moving files within the bucket, deleting files, and uploading files.

## Installation

To use the `s3bucket` package in your Go project, you can install it using the `go get` command:

```shell
go get github.com/problem-company-toolkit/s3bucket
```

## Usage

To start using the `s3bucket` package, you need to import it in your Go code:

```go
import "github.com/problem-company-toolkit/s3bucket"
```

### Creating an S3 Bucket

To create a new `Bucket` object that represents an S3 bucket, you can use the `NewS3` function:

```go
config := s3bucket.AWSConfig{
    Session: session,
    Bucket:  "your-bucket-name",
}

bucket := s3bucket.NewS3(config)
```

The `AWSConfig` struct requires an AWS `session.Session` object and the name of the S3 bucket.

### Downloading a File

To download a file from the S3 bucket, you can use the `DownloadFile` method:

```go
reader, err := bucket.DownloadFile("path/to/file.txt")
if err != nil {
    // Handle error
}

defer reader.Close()

// Read the file content from the reader
```

The `DownloadFile` method returns an `io.ReadCloser` that provides access to the downloaded file content. Make sure to close the reader when you're done reading the file.

### Moving a File

To move a file within the S3 bucket, you can use the `MoveFile` method:

```go
err := bucket.MoveFile("source/file.txt", "destination/file.txt")
if err != nil {
    // Handle error
}
```

The `MoveFile` method moves the file from the source path to the target path within the S3 bucket.

### Deleting a File

To delete a file from the S3 bucket, you can use the `DeleteFile` method:

```go
err := bucket.DeleteFile("path/to/file.txt")
if err != nil {
    // Handle error
}
```

The `DeleteFile` method deletes the specified file from the S3 bucket.

### Uploading a File

To upload a file to the S3 bucket, you can use the `UploadFile` method:

```go
file, err := os.Open("path/to/local/file.txt")
if err != nil {
    // Handle error
}
defer file.Close()

err = bucket.UploadFile(file, "destination/file.txt")
if err != nil {
    // Handle error
}
```

The `UploadFile` method takes an `io.ReadSeeker` that represents the file content and uploads it to the specified destination within the S3 bucket.