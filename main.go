// Package main go 实现的压测工具
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go-stress-testing/model"
	"go-stress-testing/server"
	"log"
	"runtime"
	"strings"
	"sync"
)

// array 自定义数组参数
type array []string

// String string
func (a *array) String() string {
	return fmt.Sprint(*a)
}

// Set set
func (a *array) Set(s string) error {
	*a = append(*a, s)

	return nil
}

var (
	concurrency uint64  = 1       // 并发数
	totalNumber uint64  = 1       // 请求数(单个并发/协程)
	debugStr            = "false" // 是否是debug
	requestURL          = ""      // 压测的url 目前支持，http/https ws/wss
	path                = ""      // curl文件路径 http接口压测，自定义参数设置
	verify              = ""      // verify 验证方法 在server/verify中 http 支持:statusCode、json webSocket支持:json
	headers     array             // 自定义头信息传递给服务器
	body        = ""              // HTTP POST方式传送数据
	maxCon      = 1               // 单个连接最大请求数
	maxProcess  = 0               //最大CPU并发
	code        = 200             //成功状态码
	http2       = false           // 是否开http2.0
	keepalive   = false           // 是否开启长连接
	sign        = "false"
	partnerId   = ""
	secretStr   = ""
	signGen     = "false"
)

func init() {
	flag.Uint64Var(&concurrency, "c", concurrency, "并发数")
	flag.Uint64Var(&totalNumber, "n", totalNumber, "请求数(单个并发/协程)")
	flag.StringVar(&debugStr, "d", debugStr, "调试模式")
	flag.StringVar(&requestURL, "u", requestURL, "压测地址")
	flag.StringVar(&path, "p", path, "curl文件路径")
	flag.StringVar(&verify, "v", verify, "验证方法 http 支持:statusCode、json webSocket支持:json")
	flag.Var(&headers, "H", "自定义头信息传递给服务器 示例:-H 'Content-Type: application/json'")
	flag.StringVar(&body, "data", body, "HTTP POST方式传送数据")
	flag.IntVar(&maxCon, "m", maxCon, "单个host最大连接数")
	flag.IntVar(&maxProcess, "mp", maxProcess, "最大CPU并发")
	flag.IntVar(&code, "code", code, "请求成功的状态码")
	flag.BoolVar(&http2, "http2", http2, "是否开http2.0")
	flag.BoolVar(&keepalive, "k", keepalive, "是否开启长连接")
	flag.StringVar(&partnerId, "partner", partnerId, "合作方编号")
	flag.StringVar(&secretStr, "secret", secretStr, "秘钥")
	flag.StringVar(&sign, "sign", sign, "是否为验签请求")
	flag.StringVar(&signGen, "sg", signGen, "是否为仅生成验签请求")
	// 解析参数
	flag.Parse()
}

// main go 实现的压测工具
// 编译可执行文件
//
//go:generate go build main.go
func main() {
	if maxProcess != 0 {
		runtime.GOMAXPROCS(maxProcess)
	}
	log.Printf("max process cpu:%v", runtime.NumCPU())
	if concurrency == 0 || totalNumber == 0 || (requestURL == "" && path == "") {
		log.Printf("示例: go-stress-testing.exe -p ./curl/test.post.curl.txt -sign 'true' -partner 'CRM' -secret '0cb4521f6723215fec1s4f9f9a756y31' -c 10 -n 20\n")
		log.Printf("压测地址或curl路径必填 \n")
		log.Printf("当前请求参数: -c %d -n %d -d %v -u %s -sign %v -partner %v -secret %v \n", concurrency, totalNumber, debugStr, requestURL, sign, partnerId, secretStr)
		flag.Usage()
		return
	}
	signGenFlag := strings.ToLower(signGen) == "true"
	debug := strings.ToLower(debugStr) == "true"
	signFlag := strings.ToLower(sign) == "true"
	if signGenFlag {
		var wg sync.WaitGroup
		var data map[string]interface{}
		if strings.TrimSpace(path) == "" {
			err := json.Unmarshal([]byte(body), &data)
			if err != nil {
				log.Println("解析JSON时出错：", err)
				return
			}
		}
		request, err := model.NewSignRequest(totalNumber, path, requestURL, verify, code, 0, debug, headers, partnerId, secretStr, data, maxCon, http2, keepalive)
		if err != nil {
			log.Printf("参数不合法 %v \n", err)
			return
		}
		wg.Add(len(request))
		for _, req := range request {
			req := req
			go func() {
				log.Printf("request:%v\n", req)
				wg.Done()
			}()
		}
		wg.Wait()
		return
	}
	if signFlag {
		log.Printf("当前请求参数: -c %d -n %d -d %v -u %s -sign %v -partner %v -secret %v -data %v \n", concurrency, totalNumber, debugStr, requestURL, sign, partnerId, secretStr, body)
		var data map[string]interface{}
		if strings.TrimSpace(path) == "" {
			err := json.Unmarshal([]byte(body), &data)
			if err != nil {
				log.Println("解析JSON时出错：", err)
				return
			}
		}
		request, err := model.NewSignRequest(totalNumber, path, requestURL, verify, code, 0, debug, headers, partnerId, secretStr, data, maxCon, http2, keepalive)
		if err != nil {
			log.Printf("参数不合法 %v \n", err)
			return
		}
		log.Printf("\n 开始启动  并发数:%d 请求数:%d\n", concurrency, totalNumber)
		// 开始处理
		server.Dispose(concurrency, totalNumber, nil, request)
	} else {
		request, err := model.NewRequest(requestURL, verify, code, 0, debug, path, headers, body, maxCon, http2, keepalive)
		if err != nil {
			log.Printf("参数不合法 %v \n", err)
			return
		}
		log.Printf("\n 开始启动  并发数:%d 请求数:%d 请求参数: \n", concurrency, totalNumber)
		request.Print()
		// 开始处理
		server.Dispose(concurrency, totalNumber, request, nil)
	}
	return
}
