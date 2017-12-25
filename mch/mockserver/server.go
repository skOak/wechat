/*
	可以提供MockServer，供测试用，比如在账号支付权限还没有下来的时候
*/
package main

import (
	"flag"
	"fmt"
	"net/http"
)

var (
	addr   string
	ApiKey string
    NotifyDelaySeconds int
)

func main() {
	flag.StringVar(&addr, "addr", ":18080", "address the mock server running on")
	flag.StringVar(&ApiKey, "apikey", "", "the apikey to use")
    flag.IntVar(&NotifyDelaySeconds, "nds", 10, "the senconds to delay before notification")
	flag.Parse()

	if ApiKey == "" {
		panic("ApiKey is empty")
	}

	http.HandleFunc("/pay/unifiedorder", UnifiedOrder)

	fmt.Println("mockserver is running at", addr)
	fmt.Println(http.ListenAndServe(addr, nil))
}

type RequestCommon struct {
	AppId string `xml:"appid,omitempty"`
	MchId string `xml:"mch_id,omitempty"`
	Sign  string `xml:"sign,omitempty"`
}

func (req *RequestCommon) FieldsMap() map[string]string {
	m1 := make(map[string]string, 3)
	if req.AppId != "" {
		m1["appid"] = req.AppId
	}
	if req.MchId != "" {
		m1["mch_id"] = req.MchId
	}
	if req.Sign != "" {
		m1["sign"] = req.Sign
	}
	return m1
}
