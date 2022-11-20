package test

import (
	"context"
	"deliverble-recording-msa/preprocess"
	_ "deliverble-recording-msa/preprocess"
	postpb "deliverble-recording-msa/protos/v1/post"
	userpb "deliverble-recording-msa/protos/v1/user"
	"github.com/stretchr/testify/require"
	_ "github.com/stretchr/testify/suite"
	"go.nhat.io/grpcmock"
	"google.golang.org/grpc"
	"testing"
)

func mockItemUserServiceServer(m ...grpcmock.ServerOption) grpcmock.ServerMockerWithContextDialer {
	opts := []grpcmock.ServerOption{grpcmock.RegisterService(userpb.RegisterUserServer)}
	opts = append(opts, m...)

	return grpcmock.MockServerWithBufConn(opts...)
}

// type UserClient interface {
//	GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*GetUserResponse, error)
//	ListUsers(ctx context.Context, in *ListUsersRequest, opts ...grpc.CallOption) (*ListUsersResponse, error)
//}
//
//type userClient struct {
//	cc grpc.ClientConnInterface
//}
//
//func NewUserClient(cc grpc.ClientConnInterface) UserClient {
//	return &userClient{cc}
//}

// https://jadekler.github.io/2020/10/08/stubbing-grpc.html
type fakeUserClient struct {
	GetUserFunc   func(ctx context.Context, in *userpb.GetUserRequest, opts ...grpc.CallOption) (*userpb.GetUserResponse, error)
	ListUsersFunc func(ctx context.Context, in *userpb.ListUsersRequest, opts ...grpc.CallOption) (*userpb.ListUsersResponse, error)
}

func (f fakeUserClient) GetUser(ctx context.Context, in *userpb.GetUserRequest, opts ...grpc.CallOption) (*userpb.GetUserResponse, error) {
	return f.GetUserFunc(ctx, in, opts...)
}

func (f fakeUserClient) ListUsers(ctx context.Context, in *userpb.ListUsersRequest, opts ...grpc.CallOption) (*userpb.ListUsersResponse, error) {
	return f.ListUsersFunc(ctx, in, opts...)
}

// https://github.com/nhatthm/grpcmock/blob/master/mock_example_test.go
func TestListPostsByUserId(t *testing.T) {
	// given: mocking user server
	ums := mockItemUserServiceServer(func(s *grpcmock.Server) {
		s.ExpectUnary("/v1.user.User/GetUser").WithPayload(&userpb.GetUserRequest{
			UserId: "1",
		}).Return(&userpb.GetUserResponse{
			UserMessage: &userpb.UserMessage{
				UserId:      "1",
				Name:        "홍주은",
				PhoneNumber: "010-1234-5678",
				Age:         23,
			},
		})
	})
	usernameWineSilver99 := "홍주은"

	testCases := []struct {
		username       string
		userId         string
		mockUserServer grpcmock.ServerMockerWithContextDialer
	}{
		{
			username:       usernameWineSilver99,
			userId:         "1",
			mockUserServer: ums,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.username, func(t *testing.T) {
			t.Parallel()

			_, dialerMockUser := tc.mockUserServer(t)

			// when
			out := &userpb.GetUserResponse{}
			err := grpcmock.InvokeUnary(context.Background(), "/v1.user.User/GetUser", &userpb.GetUserRequest{
				UserId: tc.userId,
			}, out, grpcmock.WithInsecure(), grpcmock.WithContextDialer(dialerMockUser))

			// then
			require.NoError(t, err)
			require.Equal(t, tc.username, out.UserMessage.Name)

			// given : create fake UserClient
			// https://jadekler.github.io/2020/10/08/stubbing-grpc.html
			uc := &fakeUserClient{}
			uc.GetUserFunc = func(ctx context.Context, in *userpb.GetUserRequest, opts ...grpc.CallOption) (*userpb.GetUserResponse, error) {
				return &userpb.GetUserResponse{
					UserMessage: &userpb.UserMessage{
						UserId:      "1",
						Name:        "홍주은",
						PhoneNumber: "010-1234-5678",
						Age:         23,
					},
				}, nil
			}

			ps := &preprocess.PostServer{}
			ps.UserClient = uc

			// when
			posts, err := ps.ListPostsByUserId(context.Background(), &postpb.ListPostsByUserIdRequest{
				UserId: tc.userId,
			})

			// then
			require.NoError(t, err)
			for _, post := range posts.PostMessages {
				require.Equal(t, tc.username, post.Author)
			}
		})
	}
}
