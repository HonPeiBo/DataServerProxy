package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	defaultServerPath = "http://127.0.0.1:9200"
)

func fillBadData(w http.ResponseWriter, err error) {
	w.Write([]byte(fmt.Sprintf("err: %v", err)))
}

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	path := r.URL.Path
	fmt.Printf("path: %v, method: %v\n", path, method)
	if method == "OPTIONS" {
		fmt.Printf("this is options from client %v\n", r.Header)
	}
	request, err := http.NewRequest(method, defaultServerPath+path, r.Body)
	if err != nil {
		fillBadData(w, err)
		return
	}
	request.Header["Content-Type"] = []string{"application/json"}
	rsp, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		fillBadData(w, err)
		return
	}
	defer rsp.Body.Close()

	bodyData, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		fillBadData(w, err)
		return
	}
	fmt.Printf("从服务端返回来的信息是: %v|\n", string(bodyData))
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,HEAD,DELETE,PUT")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Write(bodyData)
	if method == "OPTIONS" {
		for key, vals := range rsp.Header {
			fmt.Printf("res header: %v %v\n", key, vals)
			for _, val := range vals {
				w.Header().Set(key, val)
				fmt.Printf("   add: %v %v\n", key, val)

			}
		}
	}
	//w.Write(data2)
}

//Server HTTP方法，绑定TestHandler
func Server() {
	fmt.Printf("开始进入服务状态\n")
	http.HandleFunc("/", DefaultHandler) //根路由
	http.ListenAndServe("0.0.0.0:8009", nil)
}

func main() {
	//Get请求
	go Server()
	for {
		time.Sleep(100 * time.Second)
	}
}
