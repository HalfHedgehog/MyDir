package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"httpServer/src/global"
	"httpServer/src/initialize"
	"httpServer/src/middle"
	"httpServer/src/request"
	"httpServer/src/response"
	"httpServer/src/rpc/userRpc"
	"httpServer/src/specification"
	"strconv"
)

func main() {

	//读取全局配置文件
	initialize.Init()
	//启动gin
	router := gin.Default()
	//登陆
	router.POST("/login", login)
	//注册用户
	router.POST("/signIn", signIn)
	g1 := router.Group("/Api").Use(middle.Auth())
	g1.GET("/userInfo", getUserInfo)
	g1.POST("/updateUser", updateUser)

	_ = router.Run(":" + global.Config.Server.Port)
}

//获取用户信息
func getUserInfo(c *gin.Context) {
	userIdStr := c.GetHeader("userId")
	r := specification.Response{}
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		fmt.Println(err)
		c.JSON(500, r.Error("用户ID解析异常", nil))
		return
	}
	res, err := global.UserRpcClient.GetUserInfo(c, &userRpc.UserDto{UserID: userId})
	if err != nil {
		c.JSON(500, r.Error("获取信息异常", nil))
		return
	}
	if !res.Flag {
		c.JSON(500, r.Exception(res.Message, nil))
	}
	userInfo := response.UserInfoResp{
		UserId:         res.UserDto.UserID,
		NickName:       res.UserDto.NickName,
		ProfilePicture: res.UserDto.ProfilePicture,
	}
	c.JSON(200, r.Success("获取成功", userInfo))
}

//登陆
func login(c *gin.Context) {
	var user request.UserReq
	r := specification.Response{}
	err := c.ShouldBindJSON(&user)
	if err != nil {
		//global.Logger.WithFields(logrus.Fields{"req": c.Request.Body}).Info()
		global.Logger.Errorln("登陆参数校验错误")
		c.JSON(500, r.Error("参数错误", nil))
		return
	}
	userId, _ := strconv.ParseInt(user.UserId, 10, 64)
	rpcReq := userRpc.UserDto{
		Password: user.Password,
		UserID:   userId,
	}
	//调用rpc进行登陆操作
	res, err := global.UserRpcClient.LogIn(context.Background(), &rpcReq)
	if err != nil {
		global.Logger.WithFields(logrus.Fields{"ErrorMsq": err}).Errorln("rpc请求异常！")
		c.JSON(500, r.Error("内部异常", nil))
		return
	}
	if !res.Flag {
		c.JSON(500, r.Exception(res.Message, nil))
		return
	}
	c.JSON(200, r.Success("登陆成功！", userRpc.LogInRes{Token: res.Token}))
}

//注册用户
func signIn(c *gin.Context) {
	var user request.UserReq
	r := specification.Response{}
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(500, r.Error("参数错误", nil))
	}
	rpcReq := userRpc.UserDto{
		Password:       user.Password,
		NickName:       user.NickName,
		ProfilePicture: user.ProfilePicture,
	}
	res, erro := global.UserRpcClient.SignIn(context.Background(), &rpcReq)
	if erro != nil {
		global.Logger.WithFields(logrus.Fields{"ErrorMsq": err}).Errorln("rpc请求异常！")
		c.JSON(500, r.Error("注册失败，请稍后再试", nil))
	}
	if !res.Flag {
		c.JSON(500, r.Exception(res.Message, nil))
		return
	}
	c.JSON(200, r.Success("注册成功", response.SignInResp{UserId: strconv.FormatInt(res.UserId, 10)}))
}

//注册用户
func updateUser(c *gin.Context) {
	var user request.UserReq
	r := specification.Response{}
	err := c.ShouldBindJSON(&user)

	//UserId更新成token中的防止用户篡改
	user.UserId = c.GetHeader("userId")
	userId, _ := strconv.ParseInt(user.UserId, 10, 64)

	if err != nil {
		c.JSON(500, r.Error("参数错误", nil))
	}

	rpcReq := userRpc.UserDto{
		UserID:         userId,
		Password:       user.Password,
		NickName:       user.NickName,
		ProfilePicture: user.ProfilePicture,
	}
	res, err := global.UserRpcClient.UpdateUserInfo(context.Background(), &rpcReq)
	if err != nil {
		c.JSON(500, r.DefaultError())
		return
	}
	if !res.Flag {
		c.JSON(500, r.Exception(res.Message, nil))
		return
	}
	c.JSON(200, r.Success("更新成功！", nil))
}
