package client

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// InitS3DefaultConfig /*
// https://velog.io/@tae2089/Go%EC%97%90%EC%84%9C-S3-%EC%82%AC%EC%9A%A9%ED%95%B4%EB%B3%B4%EA%B8%B0
func (s *S3Info) InitS3DefaultConfig() (*s3.Client, error) {
	creds := credentials.NewStaticCredentialsProvider(s.AwsAccessKey, s.AwsSecretKey, "")
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(creds),
		config.WithRegion(s.AwsS3Region),
	)
	if err != nil {
		log.Printf("error: %v", err)
		panic(err)
		return nil, err
	}
	s.S3Client = s3.NewFromConfig(cfg)
	return s.S3Client, nil
}

/*
UploadRecordingV2 : uploading process to deliverble s3 bucket when served by the main restful server
*/
func (s *S3Info) UploadRecordingV2(filename string, filepath string) (*manager.UploadOutput, error) {
	uploader := manager.NewUploader(s.S3Client)

	// change filename to remove `/tmp` path location at the forehead of the name
	filename = strings.TrimPrefix(filename, "/tmp/")

	// Add the ".mp3" extension to the file name
	filename = strings.Join([]string{filename, ".mp3"}, "")

	// open file by filepath
	file, err := os.Open(filepath)
	if err != nil {
		log.Println("UploadRecording File Open Error ::::::: ", err)
		return nil, err
	}

	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(filename),
		Body:   file,
	})
	if err != nil {
		log.Fatal("UploadRecording Error ::::::: ", err)
		return nil, err
	}
	return result, nil
}

/*
UploadRecordingV2 : uploading process to deliverble s3 bucket when served by the main restful server (with using ffmpeg	to convert the file)
*/
func (s *S3Info) UploadRecording(filename string, filepath string) (*manager.UploadOutput, error) {
	uploader := manager.NewUploader(s.S3Client)

	// open file by filepath
	file, err := os.Open(filepath)
	if err != nil {
		log.Println("UploadRecording File Open Error ::::::: ", err)
		return nil, err
	}

	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(filename),
		Body:   file,
	})
	if err != nil {
		log.Fatal("UploadRecording Error ::::::: ", err)
		return nil, err
	}
	return result, nil
}

/*
DownloadRecording : downloading process to deliverble s3 bucket when served by the main restful server
*/
func (s *S3Info) DownloadRecording(targetDirectory, key string) (*os.File, error) {
	// 1. create the directory in the path
	splitKeyArr := strings.Split(key, "/")
	file := filepath.Join(targetDirectory, splitKeyArr[len(splitKeyArr)-1])
	if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
		log.Fatal("DownloadRecording Error ::::::: ", err)
		return nil, err
	}

	// 2. setting up the local file
	fd, err := os.Create(file)
	if err != nil {
		log.Fatal("DownloadRecording Error ::::::: ", err)
		return nil, err
	}

	downloader := manager.NewDownloader(s.S3Client)
	_, err = downloader.Download(context.TODO(), fd, &s3.GetObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
	})

	// https://youtrack.jetbrains.com/issue/GO-13454/Unresolved-reference-Close-for-os.File
	defer func(fd *os.File) {
		err := fd.Close()
		if err != nil {
			log.Fatal("DownloadRecording Error ::::::: ", err)
		}
	}(fd)

	return fd, err
}

/*
GetItems : get all items in the bucket
*/
func (s *S3Info) GetItems(prefix string) []types.Object {
	var responses []types.Object
	paginator := s3.NewListObjectsV2Paginator(s.S3Client, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.BucketName),
		Prefix: aws.String(prefix),
	})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			log.Fatal("GetItems Error ::::::: ", err)
		}
		responses = append(responses, page.Contents...)
	}
	return responses
}
