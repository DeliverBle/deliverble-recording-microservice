package preprocess

import (
	"context"
	"deliverble-recording-msa/data"
	postpb "deliverble-recording-msa/protos/v1/post"
	recordingpb "deliverble-recording-msa/protos/v1/recording"
	userpb "deliverble-recording-msa/protos/v1/user"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/labstack/echo"
	_ "github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net/http"
	"os"
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

func (s *S3Server) UploadRecording(ctx context.Context, req *recordingpb.UploadRecordingRequest) (*recordingpb.UploadRecordingResponse, error) {
	log.Println("UploadRecording ::::::::::::::::::: ", req.Recording)
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: "default",
		Config: aws.Config{
			Region: aws.String("us-west-2"),
		},
	})

	if err != nil {
		fmt.Printf("Failed to initialize new session: %v", err)
		return nil, err
	}

	bucketName := "deliverble-recording-bucket"
	uploader := s3manager.NewUploader(sess)
	filename := string(rune(rand.Intn(100000)))

	rec, err := os.Create("/tmp/" + filename + ".mp3")
	if err != nil {
		log.Println("Error creating file: ", err)
		return nil, err
	}
	defer rec.Close()

	fmt.Fprintf(rec, string(req.Recording))

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(rec.Name()),
		Body:   rec,
	})

	return &recordingpb.UploadRecordingResponse{}, nil
}

func UploadRecordingHandler(c echo.Context) error {
	file, err := c.FormFile("file") // file : "file" parsing
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return err
	}

	src, err := file.Open() // file api open
	if err != nil {
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

	client := recordingpb.NewRecordingTaskClient(conn)
	r, err := client.UploadRecording(ctx, &recordingpb.UploadRecordingRequest{Recording: buffer})

	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusCreated, r.Url)
}
