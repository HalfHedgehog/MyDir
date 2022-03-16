package initialize

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"httpServer/src/config"
	"httpServer/src/global"
	"httpServer/src/rpc/userRpc"
	"path"
	"time"
)

func Init() {
	//读取全局配置文件
	initConfig()
	//初始化RPC
	global.UserRpcClient = createUserRpc()
	//初始化Redis
	global.RedisHelper = createRedis()

	global.Logger = createLog()
}

// InitConfig 初始化配置
func initConfig() {
	viper.AddConfigPath("./src")
	viper.SetConfigName("apps")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("read config file failed, %v", err)
	}
	c := config.Cfg{}
	if err := viper.Unmarshal(&c); err != nil {
		fmt.Printf("unmarshal config file failed, %v", err)
	}
	global.Config = c
}

// CreateUserRpc 创建UserRpc
func createUserRpc() userRpc.SearchServiceClient {
	l, err := grpc.Dial(global.Config.UserRpc.Address, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return userRpc.NewSearchServiceClient(l)
}

// CreateRedis 初始化Redis
func createRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     global.Config.Redis.Address,
		Password: global.Config.Redis.Password,
		DB:       global.Config.Redis.DB,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	return client
}

func createLog() *logrus.Logger {
	//file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

	//if err != nil {
	//	fmt.Println("文件创建或打开失败")
	//}

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	//global.Logger.Out = file

	//设置info日志路径
	fileNameInfo := path.Join("./log/info", "sys.log")
	logWriterInfo, _ := rotatelogs.New(
		fileNameInfo+"-%Y%m%d.log",
		rotatelogs.WithLinkName(fileNameInfo),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(1*time.Hour), //
	)

	fileNameWarn := path.Join("./log/warn", "sys.log")
	logWriterWarn, _ := rotatelogs.New(
		fileNameWarn+"-%Y%m%d.log",
		rotatelogs.WithLinkName(fileNameInfo),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(1*time.Hour), //
	)

	fileNameErorr := path.Join("./log/Error", "sys.log")
	logWriteErorr, _ := rotatelogs.New(
		fileNameWarn+"-%Y%m%d.log",
		rotatelogs.WithLinkName(fileNameErorr),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(1*time.Hour), //
	)

	writerMap := lfshook.WriterMap{
		logrus.InfoLevel: logWriterInfo,
		//logrus.FatalLevel: logWriter,
		//logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriterWarn,
		logrus.ErrorLevel: logWriteErorr,
		//logrus.PanicLevel: logWriter,
	}

	logger.AddHook(lfshook.NewHook(writerMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-02-02 17:23:21",
	}))
	return logger
}
