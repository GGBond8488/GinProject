package main

import (
	"My-gin-Project/models"
	"My-gin-Project/pkg/logging"
	"My-gin-Project/pkg/setting"
	"My-gin-Project/routers"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func init()  {
	setting.Setup()
	models.Setup()
	logging.Setup()
}

func main() {
	router := routers.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.ServerSetting.HttpPort),
		Handler:        router,
		ReadTimeout:    setting.ServerSetting.ReadTimeout,
		WriteTimeout:   setting.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Printf("Listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<- quit

	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exiting")
}
	/*
	func (engine *Engine) Run(addr ...string) (err error) {
		defer func() { debugPrintError(err) }()

		address := resolveAddress(addr)
		debugPrint("Listening and serving HTTP on %s\n", address)
		err = http.ListenAndServe(address, engine)
		return
	}
	 */
	//test := gin.Default()
	//test.Run()

