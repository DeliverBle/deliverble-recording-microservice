package preprocess

import (
	"context"
	"deliverble-recording-msa/data"
	postpb "deliverble-recording-msa/protos/v1/post"
	recordingpb "deliverble-recording-msa/protos/v1/recording"
	userpb "deliverble-recording-msa/protos/v1/user"
	"deliverble-recording-msa/server/s3_server/client"
	ffw "deliverble-recording-msa/server/s3_server/client"
	"fmt"
	_ "github.com/aws/aws-sdk-go-v2/aws"
	_ "github.com/aws/aws-sdk-go-v2/config"
	_ "github.com/aws/aws-sdk-go-v2/credentials"
	_ "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	_ "github.com/aws/aws-sdk-go-v2/service/s3"
	_ "github.com/aws/aws-sdk-go-v2/service/s3/types"
	_ "github.com/aws/aws-sdk-go/aws"
	_ "github.com/aws/aws-sdk-go/aws/session"
	_ "github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	_ "github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type UserServer struct {
	userpb.UserServer
}

type PostServer struct {
	postpb.PostServer
	UserClient userpb.UserClient
}

type S3Server struct {
	recordingpb.RecordingTaskServer
}

/*
GetUser returns user message by user_id
*/
func (s *UserServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	userID := req.UserId

	var userMessage *userpb.UserMessage
	for _, ud := range data.MockUserData {
		if ud.UserId != userID {
			continue
		}
		userMessage = ud
		break
	}

	return &userpb.GetUserResponse{
		UserMessage: userMessage,
	}, nil
}

func (s *UserServer) ListUsers(context.Context, *userpb.ListUsersRequest) (*userpb.ListUsersResponse, error) {
	userMessages := make([]*userpb.UserMessage, len(data.MockUserData))
	for i, u := range data.MockUserData {
		userMessages[i] = u
	}

	return &userpb.ListUsersResponse{
		UserMessages: userMessages,
	}, nil
}

/*
ListPostsByUserId returns post messages by user_id
*/
func (s *PostServer) ListPostsByUserId(ctx context.Context, req *postpb.ListPostsByUserIdRequest) (*postpb.ListPostsByUserIdResponse, error) {
	userID := req.UserId

	response, err := s.UserClient.GetUser(ctx, &userpb.GetUserRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, err
	}

	var postMessagesByUserId []*postpb.PostMessage
	for _, up := range data.MockUserPosts {
		// if not wanted userid, then pass
		if up.UserId != userID {
			continue
		}

		// inject author name
		for _, post := range up.Posts {
			post.Author = response.UserMessage.Name
		}

		// add all posts then break
		postMessagesByUserId = append(postMessagesByUserId, up.Posts...)
		break
	}

	return &postpb.ListPostsByUserIdResponse{
		PostMessages: postMessagesByUserId,
	}, nil
}

func (s *PostServer) ListAllPosts(ctx context.Context, req *postpb.ListAllPostsRequest) (*postpb.ListAllPostsResponse, error) {
	var postMessages []*postpb.PostMessage
	for _, up := range data.MockUserPosts {
		resp, err := s.UserClient.GetUser(ctx, &userpb.GetUserRequest{
			UserId: up.UserId,
		})
		if err != nil {
			return nil, err
		}

		for _, post := range up.Posts {
			post.Author = resp.UserMessage.Name
		}

		postMessages = append(postMessages, up.Posts...)
	}

	return &postpb.ListAllPostsResponse{
		PostMessages: postMessages,
	}, nil
}

func (s *S3Server) UploadRecording(_ context.Context, req *recordingpb.UploadRecordingRequest) (*recordingpb.UploadRecordingResponse, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	info := client.S3Info{
		AwsS3Region:  os.Getenv("AP_NORTHEAST_2"),
		BucketName:   os.Getenv("DELIVERBLE_BUCKET_NAME"),
		AwsAccessKey: os.Getenv("DELIVERBLE_ACCESS_KEY"),
		AwsSecretKey: os.Getenv("DELIVERBLE_SECRET_KEY"),
	}

	client, err := info.InitS3DefaultConfig()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("S3 client init success ::::::::::::: ", client)

	filename := fmt.Sprintf("%v.mp3", time.Now().Unix())
	filepath := "/tmp/" + filename + ".mp3"

	err = ioutil.WriteFile(filepath, req.Recording, 0644)
	if err != nil {
		log.Println("Error writing to file: ", err)
		return nil, err
	}

	var response *recordingpb.UploadRecordingResponse
	recording, err := info.UploadRecording(filename, filepath)
	if err != nil {
		response = &recordingpb.UploadRecordingResponse{
			Result: false,
			Url:    "",
			Key:    "",
		}
		return response, err
	} else {
		response = &recordingpb.UploadRecordingResponse{
			Result: true,
			Url:    recording.Location,
			Key:    *recording.Key,
		}
		return response, nil
	}
}

func (s *S3Server) UploadRecordingV2(_ context.Context, req *recordingpb.UploadRecordingRequest) (*recordingpb.UploadRecordingResponse, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	info := client.S3Info{
		AwsS3Region:  os.Getenv("AP_NORTHEAST_2"),
		BucketName:   os.Getenv("DELIVERBLE_BUCKET_NAME"),
		AwsAccessKey: os.Getenv("DELIVERBLE_ACCESS_KEY"),
		AwsSecretKey: os.Getenv("DELIVERBLE_SECRET_KEY"),
	}

	client, err := info.InitS3DefaultConfig()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("S3 client init success ::::::::::::: ", client)

	// 0. work with given uploaded file
	filename := fmt.Sprintf("%v.mp3", time.Now().Unix())
	filepath := "/tmp/" + filename + ".mp3"

	err = ioutil.WriteFile(filepath, req.Recording, 0644)
	if err != nil {
		log.Println("Error writing to file: ", err)
		return nil, err
	}

	// 1. change `uploaded` mp3 file to webm file
	fileNameWebm, errChange := ffw.ChangeFileNameMp3ToWebm(filepath)
	if errChange != nil {
		log.Println("UploadRecordingV2 Error ::::::: ", errChange)
	}
	filename = *fileNameWebm

	// 2. convert converted `.webm` file to `.mp3` file
	// filename to remove `.mp3` but not to add	`.webm`
	errConvert := ffw.ConvertWebmBlobToMp3File(strings.Replace(filename, ".mp3", "", -1))
	if errConvert != nil {
		log.Println("UploadRecordingV2 Error ::::::: ", errConvert)
	}

	var response *recordingpb.UploadRecordingResponse
	recording, err := info.UploadRecording(filename, filepath)
	if err != nil {
		response = &recordingpb.UploadRecordingResponse{
			Result: false,
			Url:    "",
			Key:    "",
		}
		return response, err
	} else {
		response = &recordingpb.UploadRecordingResponse{
			Result: true,
			Url:    recording.Location,
			Key:    *recording.Key,
		}
		return response, nil
	}
}

func UploadRecordingHandlerV2(c echo.Context) error {
	file, err := c.FormFile("file") // file : "file" parsing
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return err
	}

	src, err := file.Open() // file api open
	if err != nil {
		log.Println("Error opening file: ", err)
		err := c.JSON(http.StatusInternalServerError, err)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return err
		}
		return err
	}

	defer src.Close()
	buffer := make([]byte, file.Size) // file size buf define
	_, err = src.Read(buffer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return err
	} // file read

	ctx := context.Background()

	conn, err := grpc.Dial("localhost:8020", grpc.WithInsecure())
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	taskClient := recordingpb.NewRecordingTaskClient(conn)
	r, err := taskClient.UploadRecordingV2(ctx, &recordingpb.UploadRecordingRequest{Recording: buffer})

	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	successResponse := &client.UploadRecordingHandlerResponse{
		Code: http.StatusCreated,
		Url:  r.Url,
		Key:  r.Key,
	}
	return c.JSON(http.StatusCreated, successResponse)
}

func UploadRecordingHandler(c echo.Context) error {
	file, err := c.FormFile("file") // file : "file" parsing
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return err
	}

	src, err := file.Open() // file api open
	if err != nil {
		log.Println("Error opening file: ", err)
		err := c.JSON(http.StatusInternalServerError, err)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return err
		}
		return err
	}

	defer src.Close()
	buffer := make([]byte, file.Size) // file size buf define
	_, err = src.Read(buffer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return err
	} // file read

	ctx := context.Background()

	conn, err := grpc.Dial("localhost:8020", grpc.WithInsecure())
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	taskClient := recordingpb.NewRecordingTaskClient(conn)
	r, err := taskClient.UploadRecording(ctx, &recordingpb.UploadRecordingRequest{Recording: buffer})

	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	successResponse := &client.UploadRecordingHandlerResponse{
		Code: http.StatusCreated,
		Url:  r.Url,
		Key:  r.Key,
	}
	return c.JSON(http.StatusCreated, successResponse)
}
