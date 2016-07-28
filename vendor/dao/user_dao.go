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
	//状态
	Status int

}

func NewUser() *User  {

	return &User{}
}

//查询用户信息通过用户名
func (self *User) QueryUserByUsername(username string) (*User,error)  {

	var user *User
	_,err :=db.NewSession().Select("*").From("user").Where("email=? or mobile=?",username,username).LoadStructs(&user)

	return user,err;
}