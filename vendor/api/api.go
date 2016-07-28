package api

import (
	"net/http"
	"fmt"
	"crypto/md5"
	"time"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"strings"
	"strconv"
	"dao"
	"io/ioutil"
	"log"
	"errors"
	"service"
)



type AppDto struct  {
	AppId string `json:"app_id"`
	AppKey string `json:"app_key"`
	AppName string `json:"app_name"`
	AppDesc string `json:"app_desc"`
	Status int `json:"status"`
}

//登录
func Login(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	_,data,err := CheckRequest(r)
	if err!=nil{
		util.ResponseError(w,http.StatusUnauthorized,"认证失败!")
		return;
	}
	var paramMap map[string]string
	util.CheckErr(util.ReadJsonByByte(data,&paramMap))

	username :=paramMap["username"]
	password := paramMap["password"]

	loginResult,err := service.Login(username,password)
	if err!=nil {
		util.ResponseError400(w,err.Error())
		return
	}
	util.WriteJson(w,loginResult)
}

//检查请求合法性
func CheckRequest( r *http.Request) (string,[]byte,error)  {

	bodyBytes,err := ioutil.ReadAll(r.Body)
	appId,appKey,basesign,err:=AppIsOk(r);
	if err!=nil{
		return appId,bodyBytes,err;
	}
	sign := r.Header.Get("sign")
	signs :=strings.Split(sign,".")
	if len(signs)!=2 {
		return appId,bodyBytes,errors.New("非法请求!")
	}

	if err!=nil{
		return appId,bodyBytes,errors.New("参数有误!")
	}

	var paramMap map[string]interface{}
	util.CheckErr(util.ReadJsonByByte(bodyBytes,&paramMap))

	wantSign := util.SignWithBaseSign(paramMap,appKey,basesign,nil)
	gotSign :=signs[1];
	if wantSign!=gotSign {
		log.Println("wantSign: ",wantSign,"gotSign: ",gotSign)

		return appId,bodyBytes,errors.New("签名不匹配!")
	}
	return appId,bodyBytes,nil
}





//提交应用申请
func SubmitApp(w http.ResponseWriter, r *http.Request)  {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var appDto *AppDto
	util.CheckErr(util.ReadJson(r.Body,&appDto))

	app := dao.NewAPP()
	app.AppId = fmt.Sprintf("%d",util.GenerAppId())
	app.AppName = appDto.AppName
	app.AppDesc = appDto.AppDesc
	app.Status=0
	app.AppKey = util.GenerUUId()

	if !app.Insert() {
		util.ResponseError(w,http.StatusBadRequest,"添加APP失败!")
		return;
	}else{
		util.ResponseSuccess(w)
	}

}


func AppIsOk(r *http.Request) (appId string,appKey string,basesign string,er error) {
	app_id := r.Header.Get("app_id");
	if app_id=="" {

		return "","","",errors.New("app_id不能为空!");
	}

	app := dao.NewAPP()
	app = app.QueryCanUseApp(app_id)

	if app==nil {
		return app_id,"","",errors.New("系统中没有此应用信息!");
	}
	sign :=r.Header.Get("sign")
	if sign =="" {
		return app_id,app.AppKey,"",errors.New("签名信息(sign)不能为空!");
	}
	signs := strings.Split(sign,".")
	gotSign := signs[0]

	noncestr :=r.Header.Get("noncestr")
	timestamp :=r.Header.Get("timestamp")

	if noncestr=="" {
		return app_id,app.AppKey,"",errors.New("随机码不能为空!");
	}

	if timestamp=="" {
		return app_id,app.AppKey,"",errors.New("时间戳不能为空!");
	}


	timestam64,_ := strconv.ParseInt(timestamp,10,64)
	timeBtw := time.Now().Unix()-int64(timestam64)
	if timeBtw > 5*60 {
		return app_id,app.AppKey,"",errors.New("签名已失效!");
	}

	signStr:= fmt.Sprintf("%s%s%s",app.AppKey,noncestr,timestamp)
	wantSign :=fmt.Sprintf("%X",md5.Sum([]byte(signStr)))

	if gotSign!=wantSign {
		fmt.Println("wantSign: ",wantSign,"gotSign: ",gotSign)
		return app_id,app.AppKey,"",errors.New("请求不合法!");
	}

	if app==nil{
		return app_id,app.AppKey,"",errors.New("应用信息未找到!请检查APPID是否正确!");
	}

	return app_id,app.AppKey,gotSign,nil;
}
