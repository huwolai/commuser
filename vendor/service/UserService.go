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
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"gitlab.qiyunxin.com/tangtao/utils/config"
	"net/http"
)

const (

	//用户状态 正常
	USER_STATUS_NORMAL= 1

)

type LoginResult struct  {
	OpenId string `json:"open_id"`
	Rid string `json:"r_id"`
	Token string `json:"token"`

}


//无密码登录
func LoginForNoPwd(mobile string,appId string) (*LoginResult,error) {
	if mobile=="" {
		return nil,errors.New("用户名不能为空!")
	}

	if appId=="" {
		return nil,errors.New("appId不能为空!")
	}

	user :=  dao.NewUser()
	user,err := user.QueryUserByUsername(mobile,appId)
	if err!=nil {
		log.Error(err)
		return nil,errors.New("查询用户信息失败!")
	}

	if user==nil {
		user = dao.NewUser()
		user.Mobile = mobile
		user.AppId=appId
		user.Status = USER_STATUS_NORMAL //
		user,err =user.Insert()
		if err!=nil{
			log.Error(err)
			return nil,errors.New("用户保存失败!")
		}
	}

	userdata,err := GetUserInfoFromUCR(strconv.FormatInt(user.Id,20))
	if err!=nil{
		log.Error(err)
		return nil,errors.New("获取UCR数据失败!")
	}
	log.Debug("获取到UCR的用户信息:",string(userdata))
	var loginResult *LoginResult
	err =util.ReadJsonByByte(userdata,&loginResult)
	if err!=nil{
		return nil,errors.New("解析数据错误!")
	}

	if user.OpenId=="" {
		err =user.UpdateUserOpenId(loginResult.OpenId,user.Id,appId)
		if err!=nil{
			log.Error(err)
			return nil,errors.New("更新用户中心ID失败!")
		}

	}

	return loginResult,err;
}



//登录
func Login(username string,password string,appId string) (*LoginResult,error)  {

	if username=="" {
		return nil,errors.New("用户名不能为空!")
	}

	if appId=="" {
		return nil,errors.New("appId不能为空!")
	}

   	user :=  dao.NewUser()
	user,err := user.QueryUserByUsername(username,appId)
	if err!=nil {
		log.Error(err)
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
		log.Error(err)
		return nil,errors.New("获取UCR数据失败!")
	}
	log.Debug("获取到UCR的用户信息:",string(userdata))
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


	respose,err :=network.Post(config.GetValue("ucr_url").ToString()+"/users/auth",paramsBytes,GetAuthHeader(params))
	if respose.StatusCode==http.StatusOK {

		return []byte(respose.Body),nil
	}else if respose.StatusCode==http.StatusBadRequest {
		var resultMap map[string]interface{}
		err =util.ReadJsonByByte([]byte(respose.Body),&resultMap)
		if err!=nil{

			return nil,err
		}

		return nil,errors.New(resultMap["err_msg"].(string))
	}

	return []byte(respose.Body),nil

}

func GetAuthHeader(params map[string]interface{}) map[string]string   {

	apiId :=  config.GetValue("ucr_appid").ToString();
	apikey := config.GetValue("ucr_appkey").ToString()
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
