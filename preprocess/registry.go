package preprocess

import (
	"context"
	"deliverble-recording-msa/data"
	postpb "deliverble-recording-msa/protos/v1/post"
	userpb "deliverble-recording-msa/protos/v1/user"
)

type UserServer struct {
	userpb.UserServer
}

type PostServer struct {
	postpb.PostServer
	UserClient userpb.UserClient
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
