package routes

import (
	"net/http"
	"web_app/controllers"
	"web_app/logger"
	"web_app/middlewares"

	"github.com/gin-gonic/gin"
)

//Setup 配置路由
func Setup(mode string) *gin.Engine {
	//一共三种模式，其余2种模式都当做调试模式
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) //gin设置为发布模式,就不会输出Gin-debug信息了
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	//前台API
	v1 := r.Group("api/v1/")
	//用户登录
	v1.POST("/login", controllers.UserLoginHandler)
	//退出登录
	v1.POST("/logout", controllers.LogoutHandler)
	v1.Use(middlewares.JWTAuthMiddleware())
	{
		//获取用户信息
		v1.GET("/user/info", controllers.GetUserInfoHandler)
		//添加公司
		v1.POST("/company", controllers.AddCompanyHandler)
		//获取公司列表
		v1.GET("/company", controllers.GetCompanyListHandler)
		//删除公司
		v1.DELETE("/company", controllers.DeleteCompanyHandler)
		//修改公司
		v1.PUT("/company", controllers.UpdateCompanyHandler)
		//添加任务
		v1.POST("/task", controllers.AddTaskHandler)
		//获取任务列表
		v1.GET("/task", controllers.GetTaskListHandler)
		//删除任务
		v1.DELETE("/task", controllers.DeleteTaskHandler)
		//修改任务
		v1.PUT("/task", controllers.UpdateTaskHandler)
		//开始任务
		v1.PUT("/task/:id/start", controllers.StartTaskHandler)
		//停止任务
		v1.PUT("/task/:id/stop", controllers.StopTaskHandler)
		//获取指定公司的主机列表
		v1.GET("/host/list", controllers.GetHostListHandler)
		//获取主机详情
		v1.GET("/host/detail", controllers.GetHostDetailHandler)
		//对指定url进行目录扫描
		v1.GET("/url/dirscan/start", controllers.StartURLDirScanHandler)
		//停止对指定url的目录扫描
		v1.GET("/url/dirscan/stop", controllers.StopURLDirScanHandler)
		//删除指定的url的子目录
		v1.DELETE("/url/subdir", controllers.DeleteURLSubDirHandler)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}
