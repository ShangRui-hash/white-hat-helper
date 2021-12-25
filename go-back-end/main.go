package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/logger"
	"web_app/pkg/snowflake"
	"web_app/pkg/translator"
	"web_app/routes"
	"web_app/settings"

	"go.uber.org/zap"
)

func main() {
	//检查权限
	if os.Geteuid() != 0 {
		log.Fatal("请使用root用户运行")
	}

	//0.接收命令行参数
	var configFilePath string
	flag.StringVar(&configFilePath, "config", "./config.json", "指定configPath")
	flag.Parse()

	//1.加载配置
	if err := settings.Init(configFilePath); err != nil {
		fmt.Printf("init setting failed,err:%v\n", err)
		return
	}

	//2.初始化日志
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Printf("init logger failed,err:%v\n", err)
		return
	}

	defer zap.L().Sync() //同步到日志

	//3.初始化MySQL连接
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Printf("init mysql failed,err:%v\n", err)
		return
	}
	defer mysql.Close()

	//4.初始化Redis连接
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed,err:%v\n", err)
		return
	}
	defer redis.Close()
	//5.初始化分布式ID生成器
	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		fmt.Printf("init snowflake failed,err:%v\n", err)
		return
	}
	//6.初始化gin框架 binding validate 使用的翻译器
	if err := translator.InitTrans("zh"); err != nil {
		fmt.Printf("init translator failed,err:%v\n", err)
		return
	}
	//7.注册路由
	router := routes.Setup(settings.Conf.Mode)

	//8.启动服务（优雅关机）
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Conf.Port),
		Handler: router,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")

	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}
	//修改所有任务状态为停止
	// if err := redis.StopAllRunningTask(); err != nil {
	// 	zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	// }
	zap.L().Info("Server exiting")
}
