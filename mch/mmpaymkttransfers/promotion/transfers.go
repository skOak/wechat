package promotion

import (
	"errors"
	"strconv"
	"time"

	"gopkg.in/skOak/wechat.v2/mch/core"
	wechatutil "gopkg.in/skOak/wechat.v2/util"
)

// 企业付款.
//  NOTE: 请求需要双向证书
func Transfer(clt *core.Client, req map[string]string) (resp map[string]string, err error) {
	return clt.PostXML("/mmpaymkttransfers/promotion/transfers", req)
}

type CheckNameOption string

const (
	OptionNoCheck    CheckNameOption = "NO_CHECK"
	OptionForceCheck CheckNameOption = "FORCE_CHECK"
)

type TransferRequest struct {
	XMLName struct{} `xml:"xml" json:"-"`

	// 必选参数
	PartnerTradeNo string          `xml:"partner_trade_no"` // 商户订单号，需保持唯一性（只能是字母或者数字，不能包含有其他字符）
	Openid         string          `xml:"openid"`           // 商品appid下，某用户的openid
	CheckName      CheckNameOption `xml:"check_name"`       // 校验用户姓名选项
	Amount         int64           `xml:"amount"`           // 企业付款金额，单位为分
	Desc           string          `xml:"desc"`             // MAX100;企业付款备注；备注中的敏感词会被转成字符
	SpbillCreateIp string          `xml:"spbill_create_ip"` // MAX32;该IP同在商户平台设置的IP白名单中的IP没有关联，该IP可传用户端或者服务端的IP。

	// 可选参数
	SignType   string `xml:"sign_type"`    // 签名类型，目前支持HMAC-SHA256和MD5，默认为MD5
	DeviceInfo string `xml:"device_info"`  // MAX32;微信支付分配的终端设备号
	NonceStr   string `xml:"nonce_str"`    // 随机字符串，不长于32位。NOTE: 如果为空则系统会自动生成一个随机字符串。
	ReUserName string `xml:"re_user_name"` // MAX64;收款用户真实姓名，如果check_name为FORCE_CHECK，则必填用户真实姓名
}

type TransferResponse struct {
	XMLName struct{} `xml:"xml" json:"-"`
	Body    string   `xml:"-" json:"-"` // 返回内容原文，主要是为了记录日志

	// 必选返回
	PartnerTradeNo string    `xml:"partner_trade_no"` // 商户订单号，需保持历史全局唯一性(只能是字母或者数字，不能包含有其他字符)
	PaymentNo      string    `xml:"payment_no"`       // 企业付款成功，返回的微信付款单号
	PaymentTime    time.Time `xml:"payment_time"`     // 企业付款成功时间 2015-05-19 15:26:59
}

// 企业付款2.
//  NOTE: 请求需要双向证书
func Transfer2(clt *core.Client, req *TransferRequest) (resp *TransferResponse, err error) {
	m1 := make(map[string]string, 8)
	m1["mch_appid"] = clt.AppId()
	m1["mch_id"] = clt.MchId()
	if req.DeviceInfo != "" {
		m1["device_info"] = req.DeviceInfo
	}
	if req.NonceStr != "" {
		m1["nonce_str"] = req.NonceStr
	} else {
		m1["nonce_str"] = wechatutil.NonceStr()
	}
	if req.SignType != "" {
		m1["sign_type"] = req.SignType
	}

	if req.PartnerTradeNo == "" {
		return nil, errors.New("empty partner_trade_no")
	}
	m1["partner_trade_no"] = req.PartnerTradeNo
	if req.Openid == "" {
		return nil, errors.New("empty openid")
	}
	m1["openid"] = req.Openid
	if string(req.CheckName) == "" {
		return nil, errors.New("empty check_name")
	}
	m1["check_name"] = string(req.CheckName)
	if req.ReUserName != "" {
		m1["re_user_name"] = req.ReUserName
	}
	m1["amount"] = strconv.FormatInt(req.Amount, 10)
	if req.Desc != "" {
		m1["desc"] = req.Desc
	}
	if req.SpbillCreateIp != "" {
		m1["spbill_create_ip"] = req.SpbillCreateIp
	}
	m2, err := Transfer(clt, m1)
	if err != nil {
		return nil, err
	}

	resp = &TransferResponse{
		PartnerTradeNo: m2["partner_trade_no"],
		PaymentNo:      m2["payment_no"],
	}
	if m2["payment_time"] != "" {
		resp.PaymentTime, err = time.ParseInLocation("2006-01-02 15:04:05", m2["payment_time"], wechatutil.BeijingLocation)
		if err != nil {
			return nil, err
		}
	}

	// 返回原文默认用空字符串指向
	resp.Body = m2[""]
	return
}
