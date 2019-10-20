package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"log"
	"net/http"
	"time"
)

func (s *session) download() {
	if s == nil {
		s = &session{err: fmt.Errorf("nil session")}
	} else if s.err != nil {
	} else if len(s.files) == 0 {
		s.err = fmt.Errorf("missing files")
	}
	if s.err != nil {
		return
	}
	s.err = s.downloadInternal()
}

func (s session) downloadInternal() error {
	s3Session, err := awsSession.NewSession()
	if err != nil {
		return fmt.Errorf("create aws session: %w", err)
	}
	s3Client := s3.New(s3Session)
	for _, f := range s.files {
		exists, err := f.exists(s3Client)
		if err != nil {
			return fmt.Errorf("check exists: %w", err)
		}
		if exists {
			log.Printf("%s exists, skipping", f.s3Key())
			continue
		}
		log.Printf("%s does not exist, downloading", f.s3Key())
		err = f.download(s3Session)
		if err != nil {
			return fmt.Errorf("download file: %w", err)
		}
	}
	return nil
}

func (f file) s3Key() string {
	t := time.Unix(f.LastModified/1000, 0)
	name := t.Format(time.RFC3339)
	return "us/malganis/" + name
}

func (f file) exists(s3Client *s3.S3) (bool, error) {
	req := s3.HeadObjectInput{
		Bucket: aws.String("wowah"),
		Key:    aws.String(f.s3Key()),
	}
	_, err := s3Client.HeadObject(&req)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			switch awsErr.Code() {
			case s3.ErrCodeNoSuchKey:
			case "NotFound":
				return false, nil
			}
		}
		return false, fmt.Errorf("head s3 object %s: %w", f.s3Key(), err)
	}
	return true, nil
}

func (f file) download(s3Session *awsSession.Session) error {
	getResp, err := http.Get(f.URL)
	if err != nil {
		return fmt.Errorf("get file %s: %w", f.URL, err)
	}
	defer getResp.Body.Close()
	if getResp.StatusCode != http.StatusOK {
		return fmt.Errorf("get file %s: %s", f.URL, getResp.Status)
	}

	uploader := s3manager.NewUploader(s3Session)
	uploadReq := s3manager.UploadInput{
		Bucket: aws.String("wowah"),
		Key:    aws.String(f.s3Key()),
		Body:   getResp.Body,
	}
	_, err = uploader.Upload(&uploadReq)
	if err != nil {
		return fmt.Errorf("s3 upload failed %s: %w", f.s3Key(), err)
	}
	return nil
}
