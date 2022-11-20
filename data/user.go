package data

import userpb "deliverble-recording-msa/protos/v1/user"

var MockUserData = []*userpb.UserMessage{
	{
		UserId:      "1",
		Name:        "홍주은",
		PhoneNumber: "010-1234-5678",
		Age:         23,
	},
	{
		UserId:      "2",
		Name:        "류성경",
		PhoneNumber: "010-3843-5678",
		Age:         23,
	},
	{
		UserId:      "3",
		Name:        "주송아",
		PhoneNumber: "010-3823-5678",
		Age:         23,
	},
}
