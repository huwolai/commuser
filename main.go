package main

import (
	"gitlab.qiyunxin.com/tangtao/utils/startup"
	"gitlab.qiyunxin.com/tangtao/utils/config"
	"os"
)

func main() {
	os.Setenv("APPID","commuser")
	os.Setenv("CONFIG_URL","http://configtest.qiyunxin.com")
	if !startup.IsInstall() {
		startup.InitDBData()
	}

	config.Init()
}
