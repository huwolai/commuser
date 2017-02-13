package main

import (
	"gitlab.qiyunxin.com/tangtao/utils/startup"
	"gitlab.qiyunxin.com/tangtao/utils/config"
        "github.com/gin-gonic/gin"
	"api"
	"os"
	"gitlab.qiyunxin.com/tangtao/utils/util"
)


func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, app_id, open_id")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	//os.Setenv("APPID","commuser")
	//os.Setenv("CONFIG_URL","http://configtest.qiyunxin.com")

	err := config.Init(false)
	util.CheckErr(err)
	startup.InitDBData()
	env := os.Getenv("GO_ENV")
	if env=="tests" {
		gin.SetMode(gin.DebugMode)
	}else if env== "production" {
		gin.SetMode(gin.ReleaseMode)
	}else if env == "preproduction" {
		gin.SetMode(gin.TestMode)
	}

	router := gin.Default()

	router.Use(CORSMiddleware())

	v1 :=router.Group("/v1")
	{
		v1.POST("/login",api.Login)
		v1.POST("/sms/:mobile/code",api.SendCodeSMS)
		v1.POST("/loginSMS",api.LoginForSMS)
		v1.POST("/register",api.Register)
		
		v1.POST("/chagepwd",api.ChagePassword)
		
		v1.GET("/test",api.Test)
	}

	router.Static("/swagger","./config/swagger")
	router.Run(":8080")
}
