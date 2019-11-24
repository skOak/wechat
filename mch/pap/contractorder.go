package pap

import (
	"fmt"
	"strconv"
	"time"

	"gopkg.in/skOak/wechat.v2/mch/core"
	"gopkg.in/skOak/wechat.v2/util"
)

// ContractOrder支付中签约
func ContractOrder(clt *core.Client, req map[string]string) (resp map[string]string, err error) {
	return clt.PostXML("/pay/contractorder", req)
}

// 文档：https://pay.weixin.qq.com/wiki/doc/api/pap.php?chapter=18_13&index=5
type ContractOrderRequest struct {
	XMLName struct{} `xml:"xml" json:"-"`

	// 必选参数
	// 支付相关
	Body           string `xml:"body"`             // 商品或支付单简要描述
	OutTradeNo     string `xml:"out_trade_no"`     // 商户系统内部的订单号,32个字符内、可包含字母, 其他说明见商户订单号
	TotalFee       int64  `xml:"total_fee"`        // 订单总金额，单位为分，详见支付金额
	SpbillCreateIP string `xml:"spbill_create_ip"` // APP和网页支付提交用户端ip，Native支付填调用微信支付API的机器IP。
	NotifyURL      string `xml:"notify_url"`       // 接收微信支付异步通知回调地址，通知url必须为直接可访问的url，不能携带参数。
	TradeType      string `xml:"trade_type"`       // 取值如下：JSAPI，NATIVE，APP，详细说明见参数规定

	// 签约相关
	ContractMchid          string `xml:"contract_mchid"`           // 签约商户号，必须与mch_id一致
	ContractAppid          string `xml:"contract_appid"`           // 签约公众号，必须与appid一致
	PlanId                 int64  `xml:"plan_id"`                  // 协议模板id
	ContractCode           string `xml:"contract_code"`            // 签约协议号
	RequestSerial          int64  `xml:"request_serial"`           // 商户请求签约时的序列号，要求唯一性。序列号主要用于排序，不作为查询条件。
	ContractDisplayAccount string `xml:"contract_display_account"` // 签约用户的名称，用于页面展示，参数不支持UTF8非3字节编码的字符。
	ContractNotifyUrl      string `xml:"contract_notify_url"`      // 签约信息回调通知的url

	// 可选参数
	DeviceInfo string    `xml:"device_info"` // 终端设备号(门店号或收银设备ID)，注意：PC网页或公众号内支付请传"WEB"
	NonceStr   string    `xml:"nonce_str"`   // 随机字符串，不长于32位。NOTE: 如果为空则系统会自动生成一个随机字符串。
	SignType   string    `xml:"sign_type"`   // 签名类型，默认为MD5，支持HMAC-SHA256和MD5。
	Detail     string    `xml:"detail"`      // 商品名称明细列表
	Attach     string    `xml:"attach"`      // 附加数据，在查询API和支付通知中原样返回，该字段主要用于商户携带订单的自定义数据
	FeeType    string    `xml:"fee_type"`    // 符合ISO 4217标准的三位字母代码，默认人民币：CNY，其他值列表详见货币类型
	TimeStart  time.Time `xml:"time_start"`  // 订单生成时间，格式为yyyyMMddHHmmss，如2009年12月25日9点10分10秒表示为20091225091010。其他详见时间规则
	TimeExpire time.Time `xml:"time_expire"` // 订单失效时间，格式为yyyyMMddHHmmss，如2009年12月27日9点10分10秒表示为20091227091010。其他详见时间规则
	GoodsTag   string    `xml:"goods_tag"`   // 商品标记，代金券或立减优惠功能的参数，说明详见代金券或立减优惠
	ProductId  string    `xml:"product_id"`  // trade_type=NATIVE，此参数必传。此id为二维码中包含的商品ID，商户自行定义。
	LimitPay   string    `xml:"limit_pay"`   // no_credit--指定不能使用信用卡支付
	OpenId     string    `xml:"openid"`      // rade_type=JSAPI，此参数必传，用户在商户appid下的唯一标识。
	SubOpenId  string    `xml:"sub_openid"`  // trade_type=JSAPI，此参数必传，用户在子商户appid下的唯一标识。openid和sub_openid可以选传其中之一，如果选择传sub_openid,则必须传sub_appid。
	SceneInfo  string    `xml:"scene_info"`  // 该字段用于上报支付的场景信息,针对H5支付有以下三种场景,请根据对应场景上报,H5支付不建议在APP端使用，针对场景1，2请接入APP支付，不然可能会出现兼容性问题
}

func (req *ContractOrderRequest) FieldsMap() map[string]string {
	m1 := make(map[string]string, 24)
	m1["body"] = req.Body
	m1["out_trade_no"] = req.OutTradeNo
	m1["total_fee"] = strconv.FormatInt(req.TotalFee, 10)
	m1["spbill_create_ip"] = req.SpbillCreateIP
	m1["notify_url"] = req.NotifyURL
	m1["trade_type"] = req.TradeType
	m1["contract_mchid"] = req.ContractMchid
	m1["contract_appid"] = req.ContractAppid
	m1["plan_id"] = strconv.FormatInt(req.PlanId, 10)
	m1["contract_code"] = req.ContractCode
	m1["request_serial"] = strconv.FormatInt(req.RequestSerial, 10)
	m1["contract_display_account"] = req.ContractDisplayAccount
	m1["contract_notify_url"] = req.ContractNotifyUrl
	if req.DeviceInfo != "" {
		m1["device_info"] = req.DeviceInfo
	}
	if req.NonceStr != "" {
		m1["nonce_str"] = req.NonceStr
	} else {
		m1["nonce_str"] = util.NonceStr()
	}
	if req.SignType != "" {
		m1["sign_type"] = req.SignType
	}
	if req.Detail != "" {
		m1["detail"] = req.Detail
	}
	if req.Attach != "" {
		m1["attach"] = req.Attach
	}
	if req.FeeType != "" {
		m1["fee_type"] = req.FeeType
	}
	if !req.TimeStart.IsZero() {
		m1["time_start"] = core.FormatTime(req.TimeStart)
	}
	if !req.TimeExpire.IsZero() {
		m1["time_expire"] = core.FormatTime(req.TimeExpire)
	}
	if req.GoodsTag != "" {
		m1["goods_tag"] = req.GoodsTag
	}
	if req.ProductId != "" {
		m1["product_id"] = req.ProductId
	}
	if req.LimitPay != "" {
		m1["limit_pay"] = req.LimitPay
	}
	if req.OpenId != "" {
		m1["openid"] = req.OpenId
	}
	if req.SubOpenId != "" {
		m1["sub_openid"] = req.SubOpenId
	}
	if req.SceneInfo != "" {
		m1["scene_info"] = req.SceneInfo
	}

	return m1
}

type ContractOrderResponse struct {
	XMLName struct{} `xml:"xml" json:"-"`
	Body    string   `xml:"-" json:"-"` // 返回内容原文，主要是为了记录日志

	// 必选返回
	// 支付相关
	PrepayId  string `xml:"prepay_id"`  // 微信生成的预支付回话标识，用于后续接口调用中使用，该值有效期为2小时
	TradeType string `xml:"trade_type"` // 调用接口提交的交易类型，取值如下：JSAPI，NATIVE，APP，详细说明见参数规定

	// 签约相关
	ContractResultCode string `xml:"contract_result_code"`  // 预签约结果
	ContractErrCode    string `xml:"contract_err_code"`     // 预签约错误代码
	ContractErrCodeDes string `xml:"contract_err_code_des"` // 预签约错误描述

	// 下面字段都是可选返回的(详细见微信支付文档), 为空值表示没有返回, 程序逻辑里需要判断
	DeviceInfo string `xml:"device_info"` // 调用接口提交的终端设备号。
	CodeURL    string `xml:"code_url"`    // trade_type 为 NATIVE 时有返回，可将该参数值生成二维码展示出来进行扫码支付
	MWebURL    string `xml:"mweb_url"`    // trade_type 为 MWEB 时有返回
}

// ContractOrder2 统一下单.
func ContractOrder2(clt *core.Client, req *ContractOrderRequest) (resp *ContractOrderResponse, err error) {
	m2, err := ContractOrder(clt, req.FieldsMap())
	if err != nil {
		return nil, err
	}

	// 校验 trade_type
	respTradeType := m2["trade_type"]
	if respTradeType != req.TradeType {
		err = fmt.Errorf("trade_type mismatch, have: %s, want: %s", respTradeType, req.TradeType)
		return nil, err
	}

	resp = &ContractOrderResponse{
		PrepayId:           m2["prepay_id"],
		TradeType:          respTradeType,
		DeviceInfo:         m2["device_info"],
		CodeURL:            m2["code_url"],
		MWebURL:            m2["mweb_url"],
		ContractResultCode: m2["contract_result_code"],
		ContractErrCode:    m2["contract_err_code"],
		ContractErrCodeDes: m2["contract_err_code_des"],
		Body:               m2[""], // 返回原文默认用空字符串指向
	}
	return resp, nil
}
