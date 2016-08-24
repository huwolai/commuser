package api

import (
	"net/http"
	"fmt"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"dao"
	"io/ioutil"
	"errors"
	"service"
	"github.com/gin-gonic/gin"
	"gitlab.qiyunxin.com/tangtao/utils/log"

	"strconv"
	"time"
	"math/rand"
	"redis"
)
const CODE_PREFIX  = "CODE_"

//code失效时间(单位秒)
const CODE_EXPIRE  = 60*5


type AppDto struct  {
	AppId string `json:"app_id"`
	AppKey string `json:"app_key"`
	AppName string `json:"app_name"`
	AppDesc string `json:"app_desc"`
	Status int `json:"status"`
}

type LoginParam struct {
	Username     string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type LoginForSMSParam struct  {
	//手机号
	Mobile string `json:"mobile"`
	//验证码
	Code string `json:"code"`
}

func LoginForSMS(c *gin.Context)  {
	var loginParam LoginForSMSParam
	err := c.BindJSON(&loginParam)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"数据解析错误!")
		return
	}
	appId := getAppId(c)
	if appId=="" {
		util.ResponseError400(c.Writer,"app_id不能为空!")
		return
	}

	if loginParam.Mobile=="" {
		util.ResponseError400(c.Writer,"手机号不能为空!")
		return
	}

	if loginParam.Code=="" {
		util.ResponseError400(c.Writer,"验证码不能为空!")
		return
	}

	//从缓存中读取验证码
	code :=redis.GetString(CODE_PREFIX+loginParam.Mobile)

	if code=="" {
		util.ResponseError400(c.Writer,"请先获取验证码!")
		return
	}

	if loginParam.Code!=code {
		util.ResponseError400(c.Writer,"验证码不正确!")
		return
	}

	loginResult,err :=service.LoginForNoPwd(loginParam.Mobile,appId)
	if err!=nil {
		util.ResponseError400(c.Writer,err.Error())
		return
	}

	c.JSON(http.StatusOK,loginResult)
}

//登录
func Login(c *gin.Context)  {

	var loginParam LoginParam
	err := c.BindJSON(&loginParam)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"数据解析错误!")
		return
	}

	appId := getAppId(c)

	if appId=="" {
		util.ResponseError400(c.Writer,"app_id不能为空!")
		return
	}

	loginResult,err := service.Login(loginParam.Username,loginParam.Password,appId)
	if err!=nil {
		util.ResponseError400(c.Writer,err.Error())
		return
	}
	util.WriteJson(c.Writer,loginResult)
}



func SendCodeSMS(c *gin.Context) {

	mobile := c.Param("mobile")

	if mobile=="" {

		util.ResponseError400(c.Writer,"请输入手机号!")
		return
	}

	if len(mobile)!=11 {
		util.ResponseError400(c.Writer,"手机号输入有误!")
		return
	}

	code :=redis.GetString(CODE_PREFIX+mobile)
	if code== ""{
		//code =GetRandCode()
		code="1111"
	}
	redis.SetAndExpire(CODE_PREFIX+mobile,code,CODE_EXPIRE)

	//err :=service.SendCodeSMS(mobile,code)
	//if err!=nil{
	//	log.Error(err)
	//	util.ResponseError400(c.Writer,"短信发送失败!")
	//	return
	//}

	util.ResponseSuccess(c.Writer)
}

func GetRandCode() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var code string
	for i:=0; i<4; i++ {
		code+=strconv.Itoa(r.Intn(9))
	}

	return code
}

func getAppId(c *gin.Context) string {

	appId := c.Request.Header.Get("app_id")
	if appId=="" {

		appId = c.Query("app_id")
	}

	return appId
}

//检查请求合法性
func CheckRequest( r *http.Request) (string,[]byte,error)  {

	bodyBytes,err := ioutil.ReadAll(r.Body)
	appId,_,err:=AppIsOk(r);
	if err!=nil{
		return appId,bodyBytes,err;
	}
	//sign := r.Header.Get("sign")
	//signs :=strings.Split(sign,".")
	//if len(signs)!=2 {
	//	return appId,bodyBytes,errors.New("非法请求!")
	//}
	//
	//if err!=nil{
	//	return appId,bodyBytes,errors.New("参数有误!")
	//}
	//
	//var paramMap map[string]interface{}
	//util.CheckErr(util.ReadJsonByByte(bodyBytes,&paramMap))
	//
	//wantSign := util.SignWithBaseSign(paramMap,appKey,basesign,nil)
	//gotSign :=signs[1];
	//if wantSign!=gotSign {
	//	log.Println("wantSign: ",wantSign,"gotSign: ",gotSign)
	//
	//	return appId,bodyBytes,errors.New("签名不匹配!")
	//}
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


func AppIsOk(r *http.Request) (appId string,appKey string,er error) {
	app_id := r.Header.Get("app_id");
	app_key := r.Header.Get("app_key");
	if app_id=="" {

		return "","",errors.New("app_id不能为空!");
	}

	if app_key=="" {
		return "","",errors.New("app_key不能为空!");
	}

	app := dao.NewAPP()
	app = app.QueryCanUseApp(app_id)

	if app==nil {
		return app_id,"",errors.New("系统中没有此应用信息!");
	}

	if app.AppKey!=app_key{
		return app_id,"",errors.New("app_key不正确!");
	}
	//sign :=r.Header.Get("sign")
	//if sign =="" {
	//	return app_id,app.AppKey,"",errors.New("签名信息(sign)不能为空!");
	//}
	//signs := strings.Split(sign,".")
	//gotSign := signs[0]
	//
	//noncestr :=r.Header.Get("noncestr")
	//timestamp :=r.Header.Get("timestamp")
	//
	//if noncestr=="" {
	//	return app_id,app.AppKey,"",errors.New("随机码不能为空!");
	//}
	//
	//if timestamp=="" {
	//	return app_id,app.AppKey,"",errors.New("时间戳不能为空!");
	//}


	//timestam64,_ := strconv.ParseInt(timestamp,10,64)
	//timeBtw := time.Now().Unix()-int64(timestam64)
	//if timeBtw > 5*60 {
	//	return app_id,app.AppKey,"",errors.New("签名已失效!");
	//}
	//
	//signStr:= fmt.Sprintf("%s%s%s",app.AppKey,noncestr,timestamp)
	//wantSign :=fmt.Sprintf("%X",md5.Sum([]byte(signStr)))
	//
	//if gotSign!=wantSign {
	//	fmt.Println("wantSign: ",wantSign,"gotSign: ",gotSign)
	//	return app_id,app.AppKey,"",errors.New("请求不合法!");
	//}
	//
	//if app==nil{
	//	return app_id,app.AppKey,"",errors.New("应用信息未找到!请检查APPID是否正确!");
	//}

	return app_id,app.AppKey,nil;
}
