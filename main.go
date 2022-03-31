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
    flag.StringVar(&port, "port", "9999", "api service port")
    flag.StringVar(&configPath, "config", "./config/docker_config.json", "service config filepath")
    flag.Parse()
    err := service.InitDockerRunner(configPath)
    if err != nil {
        log.Fatalf("initilize service failure: %s", err.Error())
    }
    http.HandleFunc("/languages", ctl.LanguagesController)
    http.HandleFunc("/run", ctl.RunController)
    addr := net.JoinHostPort("0.0.0.0", port)
    log.Println("running to " + addr)
    log.Fatalln(http.ListenAndServe(addr, nil))
}
