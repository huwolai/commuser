package dao

import "gitlab.qiyunxin.com/tangtao/utils/db"

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
	//密码
	Password string
	//昵称
	Nickname string
	//是否开通支付API
	IsPayapi int
	//状态
	Status int

}

func NewUser() *User  {

	return &User{}
}

func (self *User) Insert() (*User,error)  {

	result,err :=db.NewSession().InsertInto("user").Columns("app_id","open_id","email","mobile","password","nickname","status","is_payapi").Record(self).Exec()
	if err!=nil{

		return nil,err
	}

	lastId,_ :=result.LastInsertId()
	self.Id = lastId

	return self,nil
}

func (self *User) UpdateUserOpenId(openId string,rid int64,appId string) error {

	_,err :=db.NewSession().Update("user").Set("open_id",openId).Where("r_id=?",rid).Where("app_id=?",appId).Exec()

	return err
}

//查询用户信息通过用户名
func (self *User) QueryUserByUsername(username string,appId string) (*User,error)  {

	var user *User
	_,err :=db.NewSession().Select("*").From("user").Where("(email=? or mobile=? or username=?) and app_id=?",username,username,username,appId).LoadStructs(&user)

	return user,err;
}