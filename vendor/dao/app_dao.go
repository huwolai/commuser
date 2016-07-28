package dao

import (
	"time"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"fmt"
	"gitlab.qiyunxin.com/tangtao/utils/util"
)

type APP struct  {
	Id uint64
	//应用ID
	AppId string
	//应用KEY
	AppKey string
	//应用名称
	AppName string
	//应用描述
	AppDesc string
	//应用状态 0.待审核 1.已审核
	Status int
	//openID
	OpenId string
	//创建时间
	CreateTime time.Time
	//修改时间
	UpdateTime time.Time
}

func NewAPP() *APP  {

	return &APP{}
}

func (self *APP)  Insert() bool{

	self.CreateTime=time.Now()
	self.UpdateTime=time.Now()

	_,err := db.NewSession().InsertInto("app").Columns("app_id","app_key","app_name","app_desc","status","create_time","update_time").Record(self).Exec()
	if err!=nil{
		fmt.Println(err)
		return false
	}

	return true
}

//查询可用的APP
func (self *APP) QueryCanUseApp(appId string) *APP {

	var app *APP
	_,err := db.NewSession().Select("id","app_id","app_key","app_name","app_desc","status").From("app").Where("app_id=? and status=?",appId,"1").LoadStructs(&app)
	util.CheckErr(err)

	return app


}
