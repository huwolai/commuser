package main

import (
	"gitlab.qiyunxin.com/tangtao/utils/startup"
	"gitlab.qiyunxin.com/tangtao/utils/config"
        "github.com/gin-gonic/gin"
	"api"
	"os"
)



func main() {
	//os.Setenv("APPID","commuser")
	//os.Setenv("CONFIG_URL","http://configtest.qiyunxin.com")
	if !startup.IsInstall() {
		startup.InitDBData()
	}

	config.Init()

	env := os.Getenv("GO_ENV")
	if env=="tests" {
		gin.SetMode(gin.TestMode)
	}else if env== "prod" {
		gin.SetMode(gin.ReleaseMode)
	}else if env == "pre" {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.Default()

	router.POST("/login",api.Login)

	router.Run(":8080")
}
