// Package model 数据模型
package model

import (
	"log"
	"testing"
)

// TestCurl 测试函数
func TestCurl(t *testing.T) {
	// ../curl.txt
	c, err := ParseTheFile("../curl/post.curl.txt")
	log.Println(c, err)

	if err != nil {
		return
	}
	log.Printf("curl:%s \n", c.String())
	log.Printf("url:%s \n", c.GetURL())
	log.Printf("method:%s \n", c.GetMethod())
	log.Printf("body:%v \n", c.GetBody())
	log.Printf("body string:%v \n", c.GetBody())
	log.Printf("headers:%s \n", c.GetHeadersStr())
}
