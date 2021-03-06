package redis

import (
	"github.com/fzzy/radix/extra/sentinel"
	"os"
	"github.com/fzzy/radix/redis"
	"log"
	"github.com/fzzy/radix/extra/pool"
	"gitlab.qiyunxin.com/tangtao/utils/config"
)

var client *sentinel.Client

var MASTER_NAME string ="mymaster"

var pl *pool.Pool;


func Init()  {

	log.Println("init redis..")
	//var err error;
	//client,err = sentinel.NewClient("tcp","127.0.0.1:26378",10,MASTER_NAME)
	//
	//if err!=nil{
	//	log.Println("init redis error",err)
	//	os.Exit(0)
	//}

	var err error;
	pl,err = pool.NewPool("tcp",config.GetValue("redis_address").ToString(),10)

	if err!=nil{

		log.Println("redis is error=",err)
		os.Exit(0)
	}

	log.Println("init redis success")

}

//func GetConn()  (*redis.Client){
//
//	if client==nil{
//
//		Init()
//	}
//
//	conn,err  :=client.GetMaster(MASTER_NAME)
//
//	if err!=nil{
//
//		log.Fatal(err);
//		return nil;
//	}
//
//	return conn;
//}

func GetConn()  (*redis.Client){

	if pl==nil{

		Init()
	}

	conn,err  :=pl.Get()

	if err!=nil{

		log.Fatal(err);
		return nil;
	}

	return conn;
}

func PutConn(conn *redis.Client)  {

	//client.PutMaster(MASTER_NAME,conn);

	pl.Put(conn)
}

func Set(key string,value interface{})  {

	conn := GetConn();

	conn.Cmd("set",key,value)


	defer func() {
		PutConn(conn)
	}()
}

//expire 单位 秒
func SetAndExpire(key string,value interface{},expire float32)  {


	conn := GetConn();

	conn.Cmd("set",key,value)

	conn.Cmd("expire",key,expire);


}


func GetString(key string)  string{

	conn := GetConn();
	defer PutConn(conn)

	result,_:=conn.Cmd("get",key).Str()

	return result

}