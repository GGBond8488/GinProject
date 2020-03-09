package main

import (
	"My-gin-Project/pkg/setting"
	"My-gin-Project/routers"
	"fmt"
	"net/http"
)

func main()  {
	r := routers.InitRouter()
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        r,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}


	//
	s.ListenAndServe()
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
}
