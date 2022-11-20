package data

import postpb "deliverble-recording-msa/protos/v1/post"

type WrappedPostByUserId struct {
	UserId string
	Posts  []*postpb.PostMessage
}

var MockUserPosts = []*WrappedPostByUserId{
	{
		UserId: "1",
		Posts: []*postpb.PostMessage{
			{
				PostId: "1",
				// Post 서비스는 User 서비스에 user id에 해당하는 유저의 정보를 받고 그 정보중에서 유저의 이름만 빼와서 Author 이란 메세지 필드값에 유저의 이름 값을 채워주는 것이다.
				Author: "", // Post 서비스는 자체적으로 유저의 이름은 알지 못한다
				Title:  "첫번째 포스트",
				Body:   "첫번째 포스트 내용",
				Tags:   []string{"딜리버블", "테스트", "gRPC"},
			},
			{
				PostId: "2",
				Author: "",
				Title:  "두번째 포스트",
				Body:   "두번째 포스트 내용",
				Tags:   []string{"짱리버블", "나는요", "두번째"},
			},
			{
				PostId: "3",
				Author: "",
				Title:  "세번째 포스트",
				Body:   "세번째 포스트 내용",
				Tags:   []string{"SOPT", "홍주은", "포스트"},
			},
		},
	},
	{
		UserId: "3",
		Posts: []*postpb.PostMessage{
			{
				PostId: "1",
				Author: "",
				Title:  "첫번째 포스트",
				Body:   "첫번째 포스트 내용",
				Tags:   []string{"딜리버블", "테스트", "gRPC"},
			},
			{
				PostId: "2",
				Author: "",
				Title:  "두번째 포스트",
				Body:   "두번째 포스트 내용",
				Tags:   []string{"짱리버블", "나는요", "두번째"},
			},
			{
				PostId: "3",
				Author: "",
				Title:  "세번째 포스트",
				Body:   "세번째 포스트 내용",
				Tags:   []string{"SOPT", "홍주은", "포스트"},
			},
		},
	},
}
