package main

import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/pprof"
	"os"
	"strings"
	"time"

	"github.com/golang/glog"
)

func main() {

	glog.MaxSize = 1024 * 1024 * 20 // 20M自动分割
	loglevel := "3"
	logpath := "log"

	flag.Set("log_dir", logpath)
	flag.Set("alsologtostderr", "true")
	flag.Set("v", loglevel)
	flag.Parse()
	defer glog.Flush()

	pathFlag, pathErr := PathExists(logpath)
	if !pathFlag {
		glog.Error(pathErr)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/healthz", healthz)
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

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

	//ex10-func-01:为 HTTPServer 添加 0-2 秒的随机延时
	delay := randInt(0, 2000)
	time.Sleep(time.Millisecond * time.Duration(delay))
	glog.V(3).Infof("Respond in %d ms", delay)

	//ex12: add open tracing
	io.WriteString(w, "===================Details of the http request header:============\n")
	req, err := http.NewRequest("GET", "http://service1", nil)
	if err != nil {
		fmt.Printf("%s", err)
	}
	lowerCaseHeader := make(http.Header)
	for key, value := range r.Header {
		lowerCaseHeader[strings.ToLower(key)] = value
	}
	glog.Info("headers:", lowerCaseHeader)
	req.Header = lowerCaseHeader
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		glog.Info("HTTP get failed with error: ", "error", err)
	} else {
		glog.Info("HTTP get succeeded")
	}
	if resp != nil {
		resp.Write(w)
	}
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
