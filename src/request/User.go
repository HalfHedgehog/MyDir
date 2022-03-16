package request

type UserReq struct {
	UserId         string `json:"userId"`
	Password       string `json:"password"`
	NickName       string `json:"nickName"`
	ProfilePicture string `json:"profilePicture"`
}
