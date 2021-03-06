package dao

import (
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"errors"
)

type User struct  {

	Id int64
	//APPID
	AppId string
	//用户中心的OPENID
	OpenId string
	//邮箱
	Email string
	//手机号
	Mobile string
	Username string
	//密码
	Password string
	//昵称
	Nickname string
	//是否开通支付API
	IsPayapi int
	//状态
	Status int

}

type Authority struct  {
	Id 		 string	`json:"id"`
	OpenId   string	`json:"open_id"`
	Json     string	`json:"json"`
	Username string	`json:"username"`
	Nickname string	`json:"nickname"`
	Mobile   string	`json:"mobile"`
}

func NewUser() *User  {

	return &User{}
}

func (self *User) Insert() (*User,error)  {

	result,err :=db.NewSession().InsertInto("user").Columns("app_id","open_id","email","username","mobile","password","nickname","json","flag","status").Record(self).Exec()
	if err!=nil{

		return nil,err
	}

	lastId,_ :=result.LastInsertId()
	self.Id = lastId

	return self,nil
}

func (self *User) UpdateUserOpenId(openId string,rid int64,appId string) error {

	_,err :=db.NewSession().Update("user").Set("open_id",openId).Where("id=?",rid).Where("app_id=?",appId).Exec()

	return err
}

//查询用户信息通过用户名
func (self *User) QueryUserByUsername(username string,appId string) (*User,error)  {

	var user *User
	_,err :=db.NewSession().Select("*").From("user").Where("(email=? or mobile=? or username=?) and app_id=?",username,username,username,appId).LoadStructs(&user)

	return user,err;
}

func (self *User) ChagePassword(openId string,password string,newPassword string,appId string) error {

	res,err :=db.NewSession().Update("user").Set("password",newPassword).Where("open_id=?",openId).Where("app_id=?",appId).Where("password=?",password).Exec()
	
	if err==nil {
		count, _ := res.RowsAffected()
		if count < 1 {
			return errors.New("修改失败")
		}
	}
	
	return err
}
//下级
func (self *User)Lower(openId string,appId string,pageIndex uint64,pageSize uint64) ([]*Authority,error) {
	var authority []*Authority
	_,err :=db.NewSession().Select("*").From("user").Where("super_id=?",openId).Where("app_id=?",appId).Limit(pageSize).Offset((pageIndex-1)*pageSize).LoadStructs(&authority)

	return authority,err
}
func (self *User)LowerCount(openId string,appId string) (int64,error) {
	return db.NewSession().Select("count(id)").From("user").Where("super_id=?",openId).Where("app_id=?",appId).ReturnInt64()
}
//修改权限
func (self *User)Authority(appId string,openId string,json string) error {
	_,err :=db.NewSession().Update("user").Set("json",json).Where("open_id=?",openId).Where("app_id=?",appId).Limit(1).Exec()

	return err
}
//下级
func (self *User)AuthorityByOpenId(openId string,appId string) (string,error) {
	var authority *Authority
	_,err :=db.NewSession().Select("*").From("user").Where("open_id=?",openId).Where("app_id=?",appId).LoadStructs(&authority)

	return authority.Json,err
}


















