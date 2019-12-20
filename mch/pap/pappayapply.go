package pap

import (
	"strconv"

	"gopkg.in/skOak/wechat.v2/mch/core"
	"gopkg.in/skOak/wechat.v2/util"
)

// PapPayApply 申请扣款
func PayApply(clt *core.Client, req map[string]string) (resp map[string]string, err error) {
	return clt.PostXML("/pay/pappayapply", req)
}

// 文档：https://pay.weixin.qq.com/wiki/doc/api/pap.php?chapter=18_3&index=8
type PayApplyRequest struct {
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
	ContractId    string `xml:"contract_id"`    // 委托代扣签约成功后微信返回的委托代扣协议id

	// 可选参数
	NonceStr   string    `xml:"nonce_str"`  // 随机字符串，不长于32位。NOTE: 如果为空则系统会自动生成一个随机字符串。
	SignType   string    `xml:"sign_type"`  // 签名类型，默认为MD5，支持HMAC-SHA256和MD5。
	Detail     string    `xml:"detail"`     // 商品名称明细列表
	Attach     string    `xml:"attach"`     // 附加数据，在查询API和支付通知中原样返回，该字段主要用于商户携带订单的自定义数据
	FeeType    string    `xml:"fee_type"`   // 符合ISO 4217标准的三位字母代码，默认人民币：CNY，其他值列表详见货币类型
	GoodsTag   string    `xml:"goods_tag"`  // 商品标记，代金券或立减优惠功能的参数，说明详见代金券或立减优惠
	Receipt    string    `xml:"receipt"`    // 	Y，传入Y时，支付成功消息和支付详情页将出现开票入口。需要在微信支付商户平台或微信公众平台开通电子发票功能，传此字段才可生效
	SceneInfo  string    `xml:"scene_info"` // 该字段用于上报支付的场景信息,针对H5支付有以下三种场景,请根据对应场景上报,H5支付不建议在APP端使用，针对场景1，2请接入APP支付，不然可能会出现兼容性问题
}

func (req *PayApplyRequest) FieldsMap() map[string]string {
	m1 := make(map[string]string, 24)
	m1["body"] = req.Body
	m1["out_trade_no"] = req.OutTradeNo
	m1["total_fee"] = strconv.FormatInt(req.TotalFee, 10)
	m1["spbill_create_ip"] = req.SpbillCreateIP
	m1["notify_url"] = req.NotifyURL
	m1["trade_type"] = req.TradeType
	m1["contract_id"] = req.ContractId
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
	if req.GoodsTag != "" {
		m1["goods_tag"] = req.GoodsTag
	}
	if req.SceneInfo != "" {
		m1["scene_info"] = req.SceneInfo
	}
	if req.Receipt != "" {
		m1["receipt"] = req.Receipt
	}

	return m1
}

type PayApplyResponse struct {
	XMLName struct{} `xml:"xml" json:"-"`
	Body    string   `xml:"-" json:"-"` // 返回内容原文，主要是为了记录日志

	// 必选返回
	// 支付相关
	ReturnCode string `xml:"return_code"` // SUCCESS/FAIL	此字段是通信标识，非交易标识，交易是否成功需要查看result_code来判断
	ReturnMsg  string `xml:"return_msg"`  // 返回信息，如非空，为错误原因	签名失败	参数格式校验错误

	// 下面字段都是可选返回的(详细见微信支付文档), 为空值表示没有返回, 程序逻辑里需要判断
	ResultCode string `xml:"result_code"`  // 业务结果 SUCCESS/FAIL
	ErrCode    string `xml:"err_code"`     // 错误代码
	ErrCodeDes string `xml:"err_code_des"` // 错误代码描述
}

// PayApply2 统一下单.
func PayApply2(clt *core.Client, req *PayApplyRequest) (resp *PayApplyResponse, err error) {
	m2, err := PayApply(clt, req.FieldsMap())
	if err != nil {
		return nil, err
	}

	resp = &PayApplyResponse{
		ReturnCode:         m2["return_code"],
		ReturnMsg:          m2["return_msg"],
		ResultCode:         m2["result_code"],
		ErrCode:            m2["err_code"],
		ErrCodeDes:         m2["err_code_des"],
		Body:               m2[""], // 返回原文默认用空字符串指向
	}
	return resp, nil
}
