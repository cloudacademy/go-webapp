package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/labstack/echo/v4"
)

func getSHA(c echo.Context) error {
	var data string

	if s3Bucket != "" {
		bucket := aws.String(s3Bucket)
		// Get the list of items
		resp, err := s3Svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: bucket})
		if err != nil {
			log.Fatal(err)
		}

		/*
			if (len(resp.Contents) == 0) {
				fmt.Println("[ERROR] No files in bucket. ")
				os.Exit(0)
			}
		*/

		mostRecentObj := *resp.Contents[0]
		for _, item := range resp.Contents {
			if item.LastModified.After(*mostRecentObj.LastModified) {
				mostRecentObj = *item
			}
		}

		fmt.Printf("LATEST: %s\n", *mostRecentObj.Key)

		rawObject, _ := s3Svc.GetObject(
			&s3.GetObjectInput{
				Bucket: bucket,
				Key:    aws.String(*mostRecentObj.Key),
			})

		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(rawObject.Body)
		if err != nil {
			log.Fatal(err)
		}

		data = buf.String()

	} else if filePath != "" {
		//read from filepath
		files, _ := ioutil.ReadDir(filePath)
		var newestFile string
		var newestTime int64 = 0
		for _, f := range files {
			fi, err := os.Stat(fmt.Sprintf("%s/%s", filePath, f.Name()))
			if err != nil {
				fmt.Println(err)
			}
			currTime := fi.ModTime().Unix()
			if currTime > newestTime {
				newestTime = currTime
				newestFile = f.Name()
			}
		}

		content, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", filePath, newestFile))
		if err != nil {
			fmt.Println(err)
		}

		data = string(content)
	} else {
		return errors.New("location of sha256 files unknown")
	}

	c.Response().Header().Add("Content-Type", "text/plain")
	return c.String(http.StatusOK, data)
}

func ok(c echo.Context) error {
	c.Response().Header().Add("Content-Type", "text/plain")
	return c.String(http.StatusOK, "ok!")
}

var s3Bucket string
var filePath string

var s3Svc *s3.S3

func init() {
	s3Bucket = os.Getenv("S3Bucket")
	filePath = os.Getenv("FilePath")

	if s3Bucket != "" {
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String("us-west-2"),
		})
		if err != nil {
			fmt.Println(err.Error())
		}

		s3Svc = s3.New(sess)
	}

	fmt.Printf("environemnt variable S3Bucket: %s\n", s3Bucket)
	fmt.Printf("environment variable FilePath: %s\n", filePath)
}

func main() {
	e := echo.New()

	e.GET("/ok", ok)
	e.GET("/sha256", getSHA)

	e.Logger.Fatal(e.Start(":8080"))
}
