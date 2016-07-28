package service

import (
	"dao"
	"errors"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"encoding/json"
	"gitlab.qiyunxin.com/tangtao/utils/network"
	"fmt"
	"crypto/md5"

	"time"
	"strconv"
)

const (

	//用户中心URL
	UCR_URL = "http://usercenterapi.test.svc.cluster.local:8080/"
	//用户中心APPID
	UCR_APP_ID ="commuser"
	//用户中心APPKEY
	UCR_APP_KEY = "262573d2e673477f95a8f6097c6751e0"
)

type LoginResult struct  {
	OpenId string `json:"open_id"`
	Rid string `json:"r_id"`
	Token string `json:"token"`

}

//登录
func Login(username string,password string) (*LoginResult,error)  {

   	user :=  dao.NewUser()
	user,err := user.QueryUserByUsername(username)
	if err!=nil {

		return nil,errors.New("查询用户信息失败!")
	}

	if user==nil {
		return nil,errors.New("用户未找到!")
	}

	if md5Data(user.Password)!=md5Data(password){
		return nil,errors.New("密码不正确!")
	}

	userdata,err := GetUserInfoFromUCR(strconv.FormatInt(user.Id,20))
	if err!=nil{
		return nil,errors.New("获取UCR数据失败!")
	}
	var loginResult *LoginResult
	err =util.ReadJsonByByte(userdata,&loginResult)
	if err!=nil{
		return nil,errors.New("解析数据错误!")
	}

	return loginResult,err;

}

func md5Data(data string) string {

	//m5data := md5.Sum([]byte(data))

	//return string(m5data)

	return data
}

//获取用户信息从用户中心
func GetUserInfoFromUCR(rid string) ([]byte,error) {
	params :=map[string]interface{}{
		"r_id":rid,
	}
	paramsBytes,err := json.Marshal(params)
	util.CheckErr(err)

	data,err :=network.Post(UCR_URL+"/users/auth",paramsBytes,GetAuthHeader(params))

	return data,err
}

func GetAuthHeader(params map[string]interface{}) map[string]string   {

	apiId :=  UCR_APP_ID;
	apikey := UCR_APP_KEY
	noncestr :="12345"
	timestamp :=fmt.Sprintf("%d",time.Now().Unix())

	signStr := apikey +noncestr+timestamp
	bytes  := md5.Sum([]byte(signStr))
	basesign :=fmt.Sprintf("%X",bytes)

	sign := util.SignWithBaseSign(params,apikey,basesign,nil)

	sign =fmt.Sprintf("%s.%s",basesign,sign)

	return map[string]string{
		"app_id":apiId,
		"sign":sign,
		"timestamp":timestamp,
		"noncestr":noncestr,
	}
}