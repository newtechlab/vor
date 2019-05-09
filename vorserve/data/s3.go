package data

import (
	"bytes"
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/s3manager"
	"github.com/juju/errgo"
)

func newS3Storage(bucket string) (Storage, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, errgo.Mask(err)
	}
	s3s := s3Storage{
		upl:    s3manager.NewUploader(cfg),
		bucket: bucket,
	}

	// Try to upload a dummy data file to make sure it works, if possible
	// also try to delete it.
	data := bytes.NewBuffer([]byte("test"))
	_, err = s3s.upl.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3s.bucket),
		Key:    aws.String("__testobj"),
		Body:   data,
	})
	if err != nil {
		return nil, errgo.Mask(err)
	}

	req := s3s.upl.S3.DeleteObjectRequest(&s3.DeleteObjectInput{
		Bucket: aws.String(s3s.bucket),
		Key:    aws.String("__testobj"),
	})
	_, err = req.Send(context.Background())
	if err != nil {
		return nil, errgo.Mask(err)
	}

	return &s3s, nil
}

type s3Storage struct {
	upl    *s3manager.Uploader
	bucket string
}

func (s3s *s3Storage) Store(name string, data io.Reader) error {
	_, err := s3s.upl.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3s.bucket),
		Key:    aws.String(name),
		Body:   data,
	})
	return errgo.Mask(err)
}
