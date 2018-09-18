package custom

import (
	"strconv"

	"gopkg.in/skOak/wechat.v2/mch/core"
)

// DeclareOrder 订单报关.
// https://pay.weixin.qq.com/wiki/doc/api/external/declarecustom.php?chapter=18_1
func DeclareOrder(clt *core.Client, req map[string]string) (resp map[string]string, err error) {
	return clt.PostXML("/cgi-bin/mch/customs/customdeclareorder", req)
}

type DeclareOrderRequest struct {
	XMLName struct{} `xml:"xml" json:"-"`

	// 必选参数
	OutTradeNo    string `xml:"out_trade_no"`   // 商户系统内部的订单号,32个字符内、可包含字母, 其他说明见商户订单号
	TransactionId string `xml:"transaction_id"` // 微信支付返回的订单号
	Customs       string `xml:"customs"`        // 海关名
	MchCustomsNo  string `xml:"mch_customs_no"` // 商户在海关登记的备案号
	Duty          int64  `xml:"duty"`           // 关税，以分为单位，少数海关特殊要求上传该字段时需要

	// 以下字段在拆单或重新报关时必传
	SubOrderNo   string `xml:"sub_order_no"`  // 商户子订单号，如有拆单则必传
	FeeType      string `xml:"fee_type"`      // 符合ISO 4217标准的三位字母代码，默认人民币：CNY，其他值列表详见货币类型
	OrderFee     int64  `xml:"order_fee"`     // 子订单金额，以分为单位，不能超过原订单金额，order_fee=transport_fee+product_fee
	TransportFee int64  `xml:"transport_fee"` // 物流费用，以分为单位
	ProductFee   int64  `xml:"product_fee"`   // 商品费用，以分为单位

	// 可选参数
	ActionType string `xml:"action_type"` // 不传，默认是新增 ADD新增报关申请 MODIFY修改报关信息

	// 用户实名信息将以微信侧的为准，推送给海关。以下字段上传后，如与微信侧的信息不一致，会反馈给商户，便于商户收集正确的信息用于订单推送，不影响报关结果。如用户是未实名微信用户，请联系用户完成实名后再报关。
	CertType string `xml:"cert_type"` // 请传固定值IDCARD,暂只支持大陆身份证。
	CertId   string `xml:"cert_id"`   // 用户大陆身份证号，尾号为字母X的身份证号，请大写字母X。
	Name     string `xml:"name"`      // 用户姓名
}

func (req *DeclareOrderRequest) FieldsMap() map[string]string {
	m1 := make(map[string]string, 24)

	m1["out_trade_no"] = req.OutTradeNo
	m1["transaction_id"] = req.TransactionId
	m1["customs"] = req.Customs
	m1["mch_customs_no"] = req.MchCustomsNo
	m1["duty"] = strconv.FormatInt(req.Duty, 10)

	if req.SubOrderNo != "" {
		m1["sub_order_no"] = req.SubOrderNo
		m1["fee_type"] = req.FeeType
		m1["order_fee"] = strconv.FormatInt(req.OrderFee, 10)
		m1["transport_fee"] = strconv.FormatInt(req.TransportFee, 10)
		m1["product_fee"] = strconv.FormatInt(req.ProductFee, 10)
	}
	if req.ActionType != "" {
		m1["action_type"] = req.ActionType
	} else {
		m1["action_type"] = "ADD"
	}
	m1["cert_type"] = req.CertType
	m1["cert_id"] = req.CertId
	m1["name"] = req.Name
	return m1
}

type DeclareOrderResponse struct {
	XMLName struct{} `xml:"xml" json:"-"`
	Body    string   `xml:"-" json:"-"` // 返回内容原文，主要是为了记录日志

	// 必选返回
	State           string `xml:"state"`             // 状态码 UNDECLARED未申报 SUBMITTED申报已提交（订单已经送海关，商户重新申报，并且海关还有修改接口，那么记录的状态会是这个） PROCESSING申报中 SUCCESS申报成功 FAIL申报失败 EXCEPT海关接口异常
	TransactionId   string `xml:"transaction_id"`    // 微信支付返回的订单号
	OutTradeNo      string `xml:"out_trade_no"`      // 商户系统内部订单号
	ModifyTime      string `xml:"modify_time"`       // 最后更新时间，格式为yyyyMMddhhmmss，北京时间
	CertCheckResult string `xml:"cert_check_result"` // 订购人和支付人身份信息校验结果 UNCHECKED商户未上传订购人身份信息 SAME商户上传的订购人身份信息与支付人身份信息一致 DIFFERENT商户上传的订购人身份信息与支付人身份信息不一致

	// 下面字段都是拆单时必回的
	SubOrderNo string `xml:"sub_order_no"` // 商户子订单号
	SubOrderId string `xml:"sub_order_id"` // 微信子订单号
}

// DeclareOrder2 统一下单.
func DeclareOrder2(clt *core.Client, req *DeclareOrderRequest) (resp *DeclareOrderResponse, err error) {
	m2, err := DeclareOrder(clt, req.FieldsMap())
	if err != nil {
		return nil, err
	}

	resp = &DeclareOrderResponse{
		State:           m2["state"],
		TransactionId:   m2["transaction_id"],
		OutTradeNo:      m2["out_trade_no"],
		ModifyTime:      m2["modify_time"],
		CertCheckResult: m2["cert_check_result"],
		SubOrderNo:      m2["sub_order_no"],
		SubOrderId:      m2["sub_order_id"],
		Body:            m2[""], // 返回原文默认用空字符串指向
	}
	return resp, nil
}
