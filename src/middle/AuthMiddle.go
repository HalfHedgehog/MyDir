package middle

import (
	"context"
	"github.com/gin-gonic/gin"
	"httpServer/src/Util"
	"httpServer/src/enum"
	"httpServer/src/global"
	"httpServer/src/specification"
)

//Auth 鉴权
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//业务执行前
		r := specification.Response{}
		token := c.GetHeader("Token")
		if token == "" {
			c.JSON(500, r.Exception("非法访问！", nil))
			c.Abort()
			return
		}
		userId, err := Util.DecryToken(token)
		if err != nil {
			c.JSON(500, r.Exception("请重新登陆！", nil))
			c.Abort()
			return
		}
		val, err := global.RedisHelper.Get(context.Background(), string(enum.LogInPre)+userId).Result()
		if err != nil {
			c.JSON(500, r.DefaultError())
			c.Abort()
			return
		}
		if val == "" {
			c.JSON(500, r.Exception("请重新登陆！", nil))
			c.Abort()
			return
		}
		//通过了权限就吧userID放上下文里，供后面调用
		//c.Header("userId", val)
		c.Request.Header.Add("userId", val)
		//执行具体业务方法
		c.Next()
	}
}
