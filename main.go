package main

import (
	"gitlab.qiyunxin.com/tangtao/utils/startup"
	"gitlab.qiyunxin.com/tangtao/utils/config"
	. "gitlab.qiyunxin.com/tangtao/utils/route"
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"log"
	"api"
	"fmt"
)

type LogHandler struct  {

}

type AllowOrigin struct {

}

func (self AllowOrigin)  Handler(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("跨域设置..")
		w.Header().Set("Access-Control-Allow-Origin","*");
		w.Header().Set("Access-Control-Allow-Headers","*");
		w.Header().Set("Access-Control-Allow-Method","*")

		inner.ServeHTTP(w, r)
	})
}

func (self LogHandler)  Handler(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}


func GetRouters()  *mux.Router{

	//日志处理
	var logHander IHandler  =LogHandler{}
	//跨域处理
	var allowOrigin IHandler = AllowOrigin{}

	return  NewRouterWithHandle([]Route{

		Route{  //应用申请
			"Login",
			"POST",
			"/login",
			api.Login,
		},


	},[]IHandler{logHander,allowOrigin})
}


func main() {
	//os.Setenv("APPID","commuser")
	//os.Setenv("CONFIG_URL","http://configtest.qiyunxin.com")
	if !startup.IsInstall() {
		startup.InitDBData()
	}

	config.Init()

	log.Fatal(http.ListenAndServe(":8080", GetRouters()))
}
