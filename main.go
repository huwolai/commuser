package main

import (
	"gitlab.qiyunxin.com/tangtao/utils/startup"
	"gitlab.qiyunxin.com/tangtao/utils/config"
)

func main() {

	if startup.IsInstall() {
		startup.InitDBData()
	}

	config.Init()
}
