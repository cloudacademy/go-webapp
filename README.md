# Go Web App

[![Go](https://github.com/cloudacademy/go-webapp/actions/workflows/go.yml/badge.svg)](https://github.com/cloudacademy/go-webapp/actions/workflows/go.yml)

## Introduction
The [web application](https://github.com/cloudacademy/go-webapp/releases) contained within this repo is used within the [Optimize a Deployed AWS Web Application](https://cloudacademy.com/lab/aws-cloud-optimization/) hands-on lab. The web application is designed to read from a configured directory the **latest** file written into it and return back it's contents (SHA256 string). The directory from which it reads can either be a filesystem directory or an S3 bucket.

## Design
The web application is developed using Go and leverages the [Echo](https://echo.labstack.com/) framework. 

## Startup
The web application executable when started will listen on port ```8080```, and must be passed **either** the `S3Bucket` or `FilePath` environment variable, which instructs the application where to read the **latest** modified file from.

In the following example the web application is started up with the `FilePath` environment variable set to `/cloudacademy/files`

```
FilePath=/cloudacademy/files ./webapp
environemnt variable S3Bucket: 
environment variable FilePath: /cloudacademy/hashfiles

   ____    __
  / __/___/ /  ___
 / _// __/ _ \/ _ \
/___/\__/_//_/\___/ v4.6.3
High performance, minimalist Go web framework
https://echo.labstack.com
____________________________________O/_______
                                    O\
⇨ http server started on [::]:8080
```

## Application Endpoints
The ```main()``` function within the web application sets up and configures the following routes:

```
func main() {
	e := echo.New()

	e.GET("/", getSHA)
	e.GET("/sha256", getSHA)
	e.GET("/thrash", thrash)
	e.GET("/version", version)
	e.GET("/ok", ok)

	e.Logger.Fatal(e.Start(":8080"))
}
```

* `/` - root path, serves up the latest SHA256 string
* `/sha256` - serves up the latest SHA256 string (same as the root path)
* `/thrash` - used to load up the CPU, useful for causing ASG scale out events
* `/version` - metadata, returns versioning information
* `/ok` - used for health checks when configured behind a load balancer, e.g. ALB

## Build
The source code can compiled using the following command:
```
go build -o webapp 
```

The following script is used to build OS specific versions of the executable:
```
go get -v -t -d ./...
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o webapp
tar -czf release-${{ env.RELEASE_VERSION }}.linux-amd64.tar.gz webapp
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o webapp
tar -czf release-${{ env.RELEASE_VERSION }}.linux-arm64.tar.gz webapp
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o webapp
tar -czf release-${{ env.RELEASE_VERSION }}.darwin-amd64.tar.gz webapp
```
