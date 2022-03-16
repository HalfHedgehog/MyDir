package response

type UserInfoResp struct {
	UserId         int64  `json:"userId"`
	NickName       string `json:"nickName"`
	ProfilePicture string `json:"profilePicture"`
}
