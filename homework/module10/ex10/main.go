package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"runtime"
	"strings"
	"time"

	"httpserver/metrics"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	glog.MaxSize = 1024 * 1024 * 20 // 20M自动分割
	loglevel := "3"
	logpath := "log"

	//读取配置文件，设置glog的level
	config := InitConfig("config/app.properties")

	if config != nil {
		loglevel = config["loglevel"]
		logpath = config["logpath"]
	}

	flag.Set("log_dir", logpath)
	flag.Set("alsologtostderr", "true")
	flag.Set("v", loglevel)
	flag.Parse()
	defer glog.Flush()

	pathFlag, pathErr := PathExists(logpath)
	if !pathFlag {
		glog.Error(pathErr)
	}

	//ex10-func-02:增加metrics的register，注册prometheus的指标采集器（直方图采集器，采集指标是execution_latency_seconds，详见metrics.go代码
	metrics.Register()

	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/healthz", healthz)
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	//ex10-func-02:/metrics路径注册为prometheus的handler
	mux.Handle("/metrics", promhttp.Handler())

	//根据环境变量获取访问端口，默认为80
	httpport := os.Getenv("HTTP_PORT")
	//fmt.Printf("GET ENV: HTTP_PORT=%s", httpport)
	glog.V(5).Infof("GET ENV: HTTP_PORT=%s", httpport)
	if httpport == "" {
		httpport = "80"
	}
	err := http.ListenAndServe(":"+httpport, mux)
	if err != nil {
		//log.Fatal(err)
		glog.Fatal(err)
	}
}

//4、当访问 localhost/healthz 时，应返回 200
func healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<html h1>system is working... httpcode: %d </html>", 200)
	glog.V(3).Infof("system is working... httpcode: %d \n", 200)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {

	//ex10-func-02:记录整个函数的执行时间
	timer := metrics.NewTimer()
	defer timer.ObserveTotal()

	//ex10-func-01:为 HTTPServer 添加 0-2 秒的随机延时
	delay := randInt(0, 2000)
	time.Sleep(time.Millisecond * time.Duration(delay))
	glog.V(3).Infof("Respond in %d ms", delay)

	//1、接收客户端 request，并将 request 中带的 header 写入 response header
	for k, v := range r.Header {
		for _, vv := range v {
			//fmt.Printf("Header key: %s, value: %s \n", k, vv)
			glog.V(5).Infof("Header key: %s, value: %s \n", k, vv)
			w.Header().Set(k, vv)
		}
	}

	//2、读取当前系统的环境变量中的 VERSION 配置，并写入 response header
	version := runtime.Version()
	//version := os.Getenv("GOVERSION") //TODO:为什么不行？
	//path := os.Getenv("GOPATH")
	//fmt.Printf("PATH: %s \n", path)	//OK
	//fmt.Printf("VERSION: %s \n", version)
	glog.V(5).Infof("VERSION: %s \n", version)
	w.Header().Set("version", version)

	//3、Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
	clientIp := RemoteIp(r)
	retCode := 200
	//log.Printf("Clientip: %s \n", clientIp)
	glog.V(5).Infof("Clientip: %s \n", clientIp)
	//log.Printf("Http Return Code: %d", retCode)
	glog.V(5).Infof("Http Return Code: %d", retCode)

	fmt.Fprint(w, "<html h1>Welcome to cncamp...</html>")
	glog.V(3).Info("Welcome to cncamp...")
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

//读取key=value类型的配置文件
func InitConfig(path string) map[string]string {
	config := make(map[string]string)

	f, err := os.Open(path)
	if f != nil {
		defer f.Close()
	} else {
		return nil
	}
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		s := strings.TrimSpace(string(b))
		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}
		key := strings.TrimSpace(s[:index])
		if len(key) == 0 {
			continue
		}
		value := strings.TrimSpace(s[index+1:])
		if len(value) == 0 {
			continue
		}
		config[key] = value
	}
	return config
}

//PathExists 判断文件夹是否存在,不存在就创建
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		// 创建文件夹
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir failed![%v]\n", err)
		} else {
			return true, nil
		}
	}
	return false, err
}

//ex10
func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}
