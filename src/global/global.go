package global

import (
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"httpServer/src/config"
	"httpServer/src/rpc/userRpc"
)

var (
	UserRpcClient userRpc.SearchServiceClient
	Config        config.Cfg
	RedisHelper   *redis.Client
	Logger        *logrus.Logger
)
