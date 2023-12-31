package main

import (
	"fmt"
	"github.com/Eden/go-gin-example/models"
	"github.com/Eden/go-gin-example/pkg/gredis"
	"github.com/Eden/go-gin-example/pkg/logging"
	"github.com/Eden/go-gin-example/pkg/setting"
	"github.com/Eden/go-gin-example/routers"
	"github.com/fvbock/endless"
	"log"
	"syscall"
)

func main() {
	setting.Setup()
	models.Setup()
	logging.Setup()
	gredis.Setup()
	endless.DefaultReadTimeOut = setting.ServerSetting.ReadTimeout
	endless.DefaultWriteTimeOut = setting.ServerSetting.WriteTimeout
	endless.DefaultMaxHeaderBytes = 1 << 20
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)

	server := endless.NewServer(endPoint, routers.InitRouter())
	server.BeforeBegin = func(add string) {
		log.Printf("Actual pid is %d", syscall.Getpid())
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Server err: %v", err)
	}
}

//func main() {
//	s := &http.Server{
//		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
//		Handler:        routers.InitRouter(),
//		ReadTimeout:    setting.ReadTimeout,
//		WriteTimeout:   setting.WriteTimeout,
//		MaxHeaderBytes: 1 << 20,
//	}
//	s.ListenAndServe()
//}
