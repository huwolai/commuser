package main

import (
	"gitlab.qiyunxin.com/tangtao/utils/startup"
	"gitlab.qiyunxin.com/tangtao/utils/config"
        "github.com/gin-gonic/gin"
	"api"
)



func main() {
	//os.Setenv("APPID","commuser")
	//os.Setenv("CONFIG_URL","http://configtest.qiyunxin.com")
	if !startup.IsInstall() {
		startup.InitDBData()
	}

	config.Init()

	router := gin.Default()

	router.POST("/login",api.Login)

	router.Run(":8080")
}
