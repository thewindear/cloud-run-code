package main

import (
	"cloud-run-code/ctl"
	"cloud-run-code/service"
	"flag"
	"log"
	"net"
	"net/http"
)

func main() {
	var port string
	var configPath string
	flag.StringVar(&port, "port", "9999", "接口服务端口号，默认为:9999")
	flag.StringVar(&configPath, "config", "./config/docker_config.json", "服务配置文件,默认当前目录下config.json")
	flag.Parse()
	log.SetPrefix("[cloud-run-code]")
	err := service.InitDockerRunner(configPath)
	if err != nil {
		log.Fatalf("path config error: %s", err.Error())
	}
	http.HandleFunc("/run", ctl.RunController)
	addr := net.JoinHostPort("0.0.0.0", port)
	log.Println("running to " + addr)
	log.Fatalln(http.ListenAndServe(addr, nil))
}
