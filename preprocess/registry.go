package preprocess

import (
	"context"
	"deliverble-recording-msa/data"
	postpb "deliverble-recording-msa/protos/v1/post"
	recordingpb "deliverble-recording-msa/protos/v1/recording"
	userpb "deliverble-recording-msa/protos/v1/user"
	"fmt"
	_ "github.com/labstack/echo/v4"
	"io"
	"log"
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
	return &recordingpb.UploadRecordingResponse{}, nil
}

func uploadsHandler(w http.ResponseWriter, r *http.Request) {
	uploadRecording, header, err := r.FormFile("upload_recording")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(w, err)
		return
	}

	// process with the file
	dirname := "./uploads"
	err = os.MkdirAll(dirname, 0777)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(w, err)
		return
	}

	filepath := fmt.Sprintf("%s/%s", dirname, header.Filename)
	file, err := os.Create(filepath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(w, err)
		return
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}(file)

	_, err = io.Copy(file, uploadRecording)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = fmt.Fprint(w, filepath)
	if err != nil {
		log.Println(w, err)
		return
	}

	// return file path
	_, err = w.Write([]byte(filepath))
	if err != nil {
		log.Println(w, err)
		return
	}
}
