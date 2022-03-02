package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// GetItem Get item from specified S3 bucket and key, throw error otherwise
func GetItem(sess *session.Session, bucket *string, key *string) error {
	// Setup service
	svc := s3manager.NewDownloader(sess)

	// Create new file for downloaded item
	file, err := os.Create(*key)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	// Get item
	numBytes, err := svc.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(*bucket),
			Key:    aws.String(*key),
		})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
	return nil
}

// InfoItem Get info on bucket item from specified S3 bucket and key, throw error otherwise
func InfoItem(sess *session.Session, bucket *string, key *string) error {
	// Setup service
	svc := s3.New(sess)
	input := &s3.GetObjectInput{
		Bucket: aws.String(*bucket),
		Key:    aws.String(*key),
	}

	// Retrive object
	result, err := svc.GetObject(input)

	// Error report
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				fmt.Println(s3.ErrCodeNoSuchKey, aerr.Error())
			case s3.ErrCodeInvalidObjectState:
				fmt.Println(s3.ErrCodeInvalidObjectState, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}

	// Print info
	fmt.Println(result)
	return nil
}

// DeleteItem Delete item from specified S3 bucket and key, throw error otherwise
func DeleteItem(sess *session.Session, bucket *string, item *string) error {
	svc := s3.New(sess)

	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: bucket,
		Key:    item,
	})
	if err != nil {
		return err
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: bucket,
		Key:    item,
	})
	if err != nil {
		return err
	}

	return nil
}

// UploadItem Upload item to specified S3 bucket and key, throw error otherwise
func UploadItem(ctx context.Context, sess *session.Session, bucket *string, key *string) error {

	// create s3 svc
	svc := s3.New(sess)
	file, err := os.Open("mytest.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()

	buffer := make([]byte, size)

	// read file content to buffer
	file.Read(buffer)
	// converted to io.ReadSeeker type
	fileBytes := bytes.NewReader(buffer)
	// upload file to s3 bucket
	_, err = svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(*bucket),
		Key:    aws.String(*key),
		Body:   fileBytes,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == request.CanceledErrorCode {
			// If the SDK can determine the request or retry delay was canceled
			// by a context the CanceledErrorCode error code will be returned.
			fmt.Fprintf(os.Stderr, "upload canceled due to timeout, %v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "failed to upload object, %v\n", err)
		}
		os.Exit(1)
	}
	return nil
}

// Uploads a file to S3 given a bucket and object key. Also takes a duration
// value to terminate the update if it doesn't complete within that time.
//
// Usage:
//   # Upload myfile.txt to myBucket/myKey. Must complete within 10 minutes or will fail
//   go run withContext.go -b mybucket -k myKey -d 10m < myfile.txt
func main() {
	var timeout time.Duration
	bucket := flag.String("b", "", "Bucket name.")
	key := flag.String("k", "", "Object key name.")
	remove := flag.Bool("r", false, "Delete object")
	upload := flag.Bool("u", false, "Upload object")
	get := flag.Bool("g", false, "Get object")
	info := flag.Bool("i", false, "Get info on object")
	flag.DurationVar(&timeout, "d", 0, "Upload timeout.")
	flag.Parse()

	// create session
	sess := session.Must(session.NewSession())

	// create a context for timeout
	ctx := context.Background()
	var cancelFn func()
	if timeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, timeout)
	}
	// Ensure the context is canceled to prevent leaking.
	// See context package for more information, https://golang.org/pkg/context/
	if cancelFn != nil {
		defer cancelFn()
	}
	// Perform action based on flag
	if *remove == true {
		// -r delete item
		err := DeleteItem(sess, bucket, key)
		if err != nil {
			fmt.Println("Got an error deleting item:")
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("successfully deleted file on %s/%s\n", *bucket, *key)
	} else if *upload {
		// -u upload item
		err := UploadItem(ctx, sess, bucket, key)
		if err != nil {
			fmt.Println("Got an error uploading item:")
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("successfully uploaded file to %s/%s\n", *bucket, *key)
	} else if *get {
		// -i info item
		err := GetItem(sess, bucket, key)
		if err != nil {
			fmt.Println("Got an error getting item:")
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("successfully retrieved file from %s/%s\n", *bucket, *key)
	} else if *info {
		// -g get item
		err := InfoItem(sess, bucket, key)
		if err != nil {
			fmt.Println("Got an error getting item:")
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("successfully retrieved info on item from %s/%s\n", *bucket, *key)
	} else {
		// No flag specified
		fmt.Printf("No command flag detected.\nUse -u for Upload, -r for Remove, and -g for Get/Retrieve")
	}

}
