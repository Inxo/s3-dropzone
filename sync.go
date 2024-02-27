package main

import (
	"fyne.io/fyne/v2/widget"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"os"
	"strings"
	"time"
)

type Sync struct {
	Progress   *widget.ProgressBarInfinite
	SVC        *s3.S3
	BucketName string
}

func (s *Sync) Init() error {
	bucketName := os.Getenv("BUCKET_NAME")

	// Validate environment variables
	if bucketName == "" {
		log.Fatal(".env file not loaded or missing required variables")
		return nil
	}

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), ""),
		Endpoint:         aws.String(os.Getenv("AWS_ENDPOINT")),
		Region:           aws.String(os.Getenv("AWS_REGION")),
		S3ForcePathStyle: aws.Bool(true),
	}

	// Create a new AWS session
	sess, err := session.NewSession(s3Config)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new S3 client
	s.SVC = s3.New(sess)
	s.BucketName = bucketName

	return nil
}

func (s *Sync) UploadToS3(filename string, expire time.Duration) (string, error) {
	// Create S3 service client
	svc := s.SVC

	// Open the file
	filename, _ = strings.CutPrefix(filename, "file://")
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
	}(file)

	// Upload the file to S3
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(filename),
		Body:   file,
	})
	if err != nil {
		return "", err
	}
	url, err := s.generateDownloadURL(filename, expire)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (s *Sync) generateDownloadURL(filename string, expiry time.Duration) (string, error) {
	// Create S3 service client
	svc := s.SVC

	// Generate a pre-signed URL for downloading the file
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(filename),
	})
	url, err := req.Presign(expiry)
	if err != nil {
		return "", err
	}

	return url, nil
}
