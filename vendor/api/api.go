package api

import (
	"net/http"
	"fmt"
	"gitlab.qiyunxin.com/tangtao/utils/network"
	"crypto/md5"
	"time"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"encoding/json"
)

//登录
func Login(w http.ResponseWriter, r *http.Request)  {

	fmt.Println("---登录---")

	apikey :="262573d2e673477f95a8f6097c6751e0"
	noncestr :="12345"
	timestamp :=fmt.Sprintf("%d",time.Now().Unix())

	signStr := apikey +noncestr+timestamp
	bytes  := md5.Sum([]byte(signStr))
	basesign :=fmt.Sprintf("%X",bytes)

	params :=map[string]interface{}{

		"r_id":"1",
	}
	sign := util.SignWithBaseSign(params,apikey,basesign,nil)

	sign =fmt.Sprintf("%s.%s",basesign,sign)

	paramsBytes,err := json.Marshal(params)
	util.CheckErr(err)

	data,err :=network.Post("http://usercenterapi.test.svc.cluster.local:8080/users/auth",paramsBytes,map[string]string{
		"app_id":"commuser",
		"sign":sign,
		"timestamp":timestamp,
		"noncestr":noncestr,
	})

	util.CheckErr(err)

	fmt.Println(string(data))
}
