package main

import (
	"azmesh-gateway/config"
	"azmesh-gateway/services"
	"flag"
	l4g "github.com/alecthomas/log4go"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	l4g.LoadConfiguration("./log4go.xml")
	// add go run command flag config
	var port, configPath string
	flag.StringVar(&configPath, "config", "./config/config.json", "config file path")
	flag.StringVar(&port, "port", "4180", "server port")
	flag.Parse()

	// start pprof
	go func() {
		err := http.ListenAndServe(":6065", nil)
		if err != nil {
			l4g.Error("pprof err:%v", err)
			return
		}
		l4g.Info("pprof listen:6065")
	}()

	// init config
	config.InitConfig(configPath)

	// start server
	l4g.Info("server start！")
	services.StartServer(port)

	//等待退出
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-c
}
