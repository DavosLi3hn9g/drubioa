package routers

import (
	"VGO/pi/internal/auth"
	"VGO/pi/internal/config"
	"VGO/pi/internal/routers/api/socket"
	"VGO/pi/internal/routers/api/v1"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
)

// List 路由列表设定
var configENV = config.ENV

func List(isDebug bool) *gin.Engine {
	r := gin.New()
	ginConf := cors.DefaultConfig()
	ginConf.AllowCredentials = true
	ginConf.AddAllowHeaders("Authorization,x-requested-with,withcredentials")
	if isDebug {
		ginConf.AllowOrigins = strings.Split(configENV["allow_origin"], ",")
		gin.SetMode(gin.TestMode)
	} else {
		ginConf.AllowAllOrigins = true
		gin.SetMode(gin.ReleaseMode)
	}
	r.Use(gin.Logger(), gin.Recovery(), cors.New(ginConf))
	r.NoRoute(Index)
	r.GET("/", Index)
	r.GET("/index.html", Index)
	var s = socket.Client{}
	r.GET("/ws", s.WsPrint)
	r.Static("/data", "./data")
	StaticList(r, "./templates/"+configENV["template"]+"/dist/")
	var vAuth = new(v1.Auth)
	var vLogsCall = new(v1.LogsCall)
	var vLogsSms = new(v1.LogsSms)
	var vIntention = new(v1.Intention)
	var vPolicy = new(v1.Policy)
	var vSetting = new(v1.Setting)
	var vUser = new(v1.User)
	var vLogsPi = new(v1.LogsPi)
	var vCloud = new(v1.Cloud)
	var vDevice = new(v1.Device)
	var vTest = new(v1.Test)
	var vVersion = new(v1.Version)

	api := r.Group("/api/v1")
	{
		anonymous := api
		anonymous.GET("/first", vAuth.First)
		anonymous.POST("/login", vAuth.Login)
		anonymous.POST("/register", vAuth.Register)
		anonymous.GET("/logout", vAuth.Logout)
		anonymous.GET("/version", vVersion.GetVersion)

		api.Use(auth.CheckAuth())
		au := api.Group("/auth")
		au.GET("/info", vAuth.Info)
		au.POST("/update", vAuth.Update)

		api.GET("/version_update", vVersion.DoUpdate)
		api.GET("/version/download_progress", vVersion.Progress)
		api.GET("/reload", vVersion.Reload)

		api.GET("/logs_at/all", vLogsPi.All)
		api.POST("/logs_at/add", vLogsPi.Add)
		api.POST("/logs_at/del", vLogsPi.Clear)
		api.GET("/logs_gpio/list", vLogsPi.GPIO)
		api.GET("/logs_sys/list", vLogsPi.SysCache)
		api.GET("/logs_call/list", vLogsCall.List)
		api.POST("/logs_call/del", vLogsCall.Del)
		api.GET("/logs_call/cost", vLogsCall.Cost)
		api.GET("/logs_sms/list", vLogsSms.List)
		api.GET("/logs_sms/update", vLogsSms.AddOrUpdate)
		api.POST("/logs_sms/del", vLogsSms.Del)

		api.GET("/setting/list", vSetting.List)
		api.POST("/setting/set", vSetting.Set)

		api.GET("/policy/list", vPolicy.List)
		api.POST("/policy/add_update", vPolicy.AddOrUpdate)
		api.POST("/policy/del", vPolicy.Del)

		api.GET("/intention/list", vIntention.List)
		api.POST("/intention/add_update", vIntention.AddOrUpdateSid)
		api.POST("/intention/del", vIntention.Del)

		api.POST("/query/add_update", vIntention.AddOrUpdateQuery)
		api.POST("/query/del", vIntention.DelQuery)

		api.GET("/user/list", vUser.List)
		api.POST("/user/add_update", vUser.AddOrUpdate)
		api.POST("/user/add_call", vUser.AddCall)
		api.POST("/user/del", vUser.Del)
		api.POST("/user/del_call", vUser.DelCall)

		api.POST("/cloud/list_oss", vCloud.ListOSS)
		api.GET("/cloud/list_tts", vCloud.ListTTS)

		api.GET("/device/list_tty", vDevice.ListTTY)

		api.POST("/test/file_asr", vTest.FileASR)
		api.GET("/test/list_pcm", vTest.ListPCM)
		api.POST("/test/setting", vTest.Setting)

	}

	return r
}
func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "iQiar智能反骚扰-控制面板",
	})
}
func HTML(c *gin.Context) {
	if c.Request.URL.Path != "/" {
		http.Error(c.Writer, "Not found", http.StatusNotFound)
		return
	}
	if c.Request.Method != "GET" {
		http.Error(c.Writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func StaticList(r *gin.Engine, path string) *gin.Engine {
	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".html") {
			r.LoadHTMLFiles(path + f.Name())
		} else {
			r.StaticFile(f.Name(), path+f.Name())
		}
	}
	r.Static("/css", "./templates/"+configENV["template"]+"/dist/css")
	r.Static("/img", "./templates/"+configENV["template"]+"/dist/img")
	r.Static("/js", "./templates/"+configENV["template"]+"/dist/js")
	return r
}
