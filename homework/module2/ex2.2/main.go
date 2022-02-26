package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"runtime"

	"net/http/pprof"

	"github.com/golang/glog"
	"github.com/thinkeridea/go-extend/exnet"
)

func main() {

	//flag.Set("v", "4")
	//flag.Set("log_dir", "log")
	flag.Parse()
	glog.V(2).Info("Starting http server...")

	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/healthz", healthz)
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	err := http.ListenAndServe(":80", mux)
	if err != nil {
		log.Fatal(err)
	}
}

//当访问 localhost/healthz 时，应返回 200
func healthz(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "200\n")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("entering root handler")
	user := r.URL.Query().Get("user")
	if user != "" {
		io.WriteString(w, fmt.Sprintf("hello [%s]\n", user))
	} else {
		io.WriteString(w, "hello [stranger]\n")
	}
	//接收客户端 request，并将 request 中带的 header 写入 response header
	io.WriteString(w, "===================Details of the http request header:============\n")
	for k, v := range r.Header {
		io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
	}
	//读取当前系统的环境变量中的 VERSION 配置，并写入 response header
	io.WriteString(w, "===================Details of the env of Go:============\n")
	io.WriteString(w, fmt.Sprintf("GOVERSION=%s\n", runtime.Version()))

	//Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
	glog.V(2).Info("The Client IP is : ", RemoteIp(r))
	w.WriteHeader(200)
	glog.V(2).Info("Http retcode is : ", 200)
}

func RemoteIp(req *http.Request) string {
	remoteAddr := req.RemoteAddr
	if ip := exnet.ClientPublicIP(req); ip != "" {
		remoteAddr = ip
	} else if ip := exnet.ClientIP(req); ip != "" {
		remoteAddr = ip
	} else if ip := req.Header.Get("X-Real-IP"); ip != "" {
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
