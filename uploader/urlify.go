package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "S3 uploader",
		Version: "1.0.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "output",
				Value: "line",
				Usage: "format of the URLs output ('line', 'json', 'csv')",
			},
		},
		Action: func(c *cli.Context) error {

			// CLI variables
			region := "us-east-2"
			bucket := "urlify"
			prefix := "urlify"

			// make sure there are some arguments
			// TODO: print usage
			if c.Args().Len() == 0 {
				cli.ShowAppHelpAndExit(c, 0)
			}

			// Initialize a session that the SDK will use to load configuration,
			// credentials, and region from the shared config file. (~/.aws/config).
			sess, err := session.NewSession(&aws.Config{
				Region: aws.String(region)},
			)
			if err != nil {
				fmt.Printf("could not load session in %s: %v\n", region, err)
				os.Exit(1)
			}

			// Setup the S3 Upload Manager
			uploader := s3manager.NewUploader(sess)

			var urls []string

			for i := 0; i < c.Args().Len(); i++ {

				// This idiom (reading the file to a byte array just to re-create
				// a reader from it instead of just using 'file' directly)
				// is a little strange. We do this to allow the
				// script to work with non-seekable files. If we just provide
				// file directly to UploadInput a la
				//
				// ```go
				// _, err = uploader.Upload(&s3manager.UploadInput{
				// 	Bucket: aws.String(bucket),
				// 	Key:    aws.String(key),
				// 	Body:   file
				// })
				// ```
				//
				// The script won't work with non-seekable files. So
				// stuff like process substitution
				//
				// ```bash
				// go run main.go <(echo "input here")
				// ```
				//
				// won't work
				localFileName := c.Args().Get(i)
				file, err := os.Open(localFileName)
				if err != nil {
					fmt.Printf("failed to open file %q, %v\n", c.Args().Get(i), err)
					continue
				}
				defer file.Close()

				fileBytes, err := io.ReadAll(file)
				if err != nil {
					fmt.Printf("failed to open file %q, %v\n", localFileName, err)
					continue
				}

				// Upload the file to S3 Bucket
				remoteFileName := namesgenerator.GetRandomName(3)
				extension := filepath.Ext(localFileName)
				key := fmt.Sprintf("%s/%s%s", prefix, remoteFileName, extension)
				_, err = uploader.Upload(&s3manager.UploadInput{
					Bucket: aws.String(bucket),
					Key:    aws.String(key),
					Body:   bytes.NewReader(fileBytes),
				})

				if err != nil {
					fmt.Printf("failed to upload file, %v\n", err)
					continue
				}

				svc := s3.New(sess)
				req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
					Bucket: aws.String(bucket),
					Key:    aws.String(key),
				})
				urlStr, err := req.Presign(5 * time.Minute)

				if err != nil {
					fmt.Println("failed to sign request", err)
				}
				urls = append(urls, urlStr)
			}

			switch c.String("output") {
			case "line":
				for _, url := range urls {
					fmt.Println(url)
				}
			case "json":
				jsonOut, err := json.Marshal(urls)
				if err != nil {
					fmt.Println("failed to serialize urls to JSON", err)
				}
				fmt.Println(string(jsonOut))
			case "csv":
				fmt.Println(strings.Join(urls, ","))
			default:
				fmt.Println("unrecognized output format")
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
