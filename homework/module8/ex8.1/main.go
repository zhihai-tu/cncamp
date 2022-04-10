package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"

	"net/http/pprof"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/healthz", healthz)
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	//根据环境变量获取访问端口，默认为80
	httpport := os.Getenv("HTTP_PORT")
	fmt.Printf("GET ENV: HTTP_PORT=%s", httpport)
	if httpport == "" {
		httpport = "80"
	}
	err := http.ListenAndServe(":"+httpport, mux)
	if err != nil {
		log.Fatal(err)
	}
}

//4、当访问 localhost/healthz 时，应返回 200
func healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "system is working... httpcode: %d \n", 200)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {

	//1、接收客户端 request，并将 request 中带的 header 写入 response header
	for k, v := range r.Header {
		for _, vv := range v {
			fmt.Printf("Header key: %s, value: %s \n", k, vv)
			w.Header().Set(k, vv)
		}
	}

	//2、读取当前系统的环境变量中的 VERSION 配置，并写入 response header
	version := runtime.Version()
	//version := os.Getenv("GOVERSION") //TODO:为什么不行？
	//path := os.Getenv("GOPATH")
	//fmt.Printf("PATH: %s \n", path)	//OK
	fmt.Printf("VERSION: %s \n", version)
	w.Header().Set("version", version)

	//3、Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
	clientIp := RemoteIp(r)
	retCode := 200
	log.Printf("Clientip: %s \n", clientIp)
	log.Printf("Http Return Code: %d", retCode)
}

func RemoteIp(req *http.Request) string {
	remoteAddr := req.RemoteAddr
	if ip := req.Header.Get("X-Real-IP"); ip != "" {
		remoteAddr = ip
	} else if ip = req.Header.Get("X-Forwarded-For"); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}

	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}

	return remoteAddr
}
