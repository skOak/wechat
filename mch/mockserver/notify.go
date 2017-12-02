package main

import (
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"time"
)

type NotifyResp struct {
	ReturnCode string `xml:"return_code"` // SUCCESS/FAIL	SUCCESS表示商户接收通知成功并校验成功
	ReturnMsg  string `xml:"return_msg"`  // 返回信息，如非空，为错误原因：	签名失败	参数格式校验错误
}

func Notify(notifyUrl string, contentType string, body io.Reader, waitMS int64) error {
	if notifyUrl == "" {
		return errors.New("invalid notify url")
	}
	if waitMS > 0 {
		<-time.After(time.Millisecond * time.Duration(waitMS))
	}
	resp, err := http.Post(notifyUrl, contentType, body)
	if err != nil {
		return err
	}
	res := &NotifyResp{}
	if err := xml.NewDecoder(resp.Body).Decode(res); err != nil {
		return err
	}
	if res.ReturnCode != "SUCCESS" {
		return errors.New(res.ReturnMsg)
	}
	return nil
}
