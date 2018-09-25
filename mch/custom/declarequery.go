package custom

import (
	"fmt"
	"strconv"

	"gopkg.in/skOak/wechat.v2/mch/core"
)

// DeclareQuery 报关查询.
// https://pay.weixin.qq.com/wiki/doc/api/external/declarecustom.php?chapter=18_2
func DeclareQuery(clt *core.Client, req map[string]string) (resp map[string]string, err error) {
	return clt.PostXML("/cgi-bin/mch/customs/customdeclarequery", req)
}

type DeclareQueryRequest struct {
	XMLName struct{} `xml:"xml" json:"-"`

	// 必选参数
	// 四选一，同时存在时优先级如下：sub_order_id> sub_order_no> transaction_id> out_trade_no
	OutTradeNo    string `xml:"out_trade_no"`   // 商户系统内部的订单号,32个字符内、可包含字母, 其他说明见商户订单号
	TransactionId string `xml:"transaction_id"` // 微信支付返回的订单号
	SubOrderNo    string `xml:"sub_order_no"`   // 商户子订单号，如有拆单则必传
	SubOrderId    string `xml:"sub_order_id"`   // 微信子订单号
	Customs       string `xml:"customs"`        // 海关名
}

func (req *DeclareQueryRequest) FieldsMap() map[string]string {
	m1 := make(map[string]string, 24)

	m1["out_trade_no"] = req.OutTradeNo
	m1["transaction_id"] = req.TransactionId
	m1["customs"] = req.Customs

	if req.SubOrderNo != "" {
		m1["sub_order_no"] = req.SubOrderNo
		m1["sub_order_id"] = req.SubOrderId
	}
	return m1
}

type DeclareQueryResponse struct {
	XMLName struct{} `xml:"xml" json:"-"`
	Body    string   `xml:"-" json:"-"` // 返回内容原文，主要是为了记录日志

	// 必选返回
	TransactionId string      `xml:"transaction_id"` // 微信支付返回的订单号
	Count         int         `xml:"count"`          // 笔数
	QueryList     []QueryItem `xml:"query_list"`     // 查询结果数据
}

type QueryItem struct {
	// 必选返回
	Customs         string `xml:"customs"`           // 海关
	State           string `xml:"state"`             // 状态码 UNDECLARED未申报 SUBMITTED申报已提交（订单已经送海关，商户重新申报，并且海关还有修改接口，那么记录的状态会是这个） PROCESSING申报中 SUCCESS申报成功 FAIL申报失败 EXCEPT海关接口异常
	ModifyTime      string `xml:"modify_time"`       // 最后更新时间，格式为yyyyMMddhhmmss，北京时间
	CertCheckResult string `xml:"cert_check_result"` // 订购人和支付人身份信息校验结果 UNCHECKED商户未上传订购人身份信息 SAME商户上传的订购人身份信息与支付人身份信息一致 DIFFERENT商户上传的订购人身份信息与支付人身份信息不一致

	// 可选返回
	SubOrderNo   string `xml:"sub_order_no"`   // 商户子订单号
	SubOrderId   string `xml:"sub_order_id"`   // 微信子订单号
	MchCustomsNo string `xml:"mch_customs_no"` // 商户在海关登记的备案号
	FeeType      string `xml:"fee_type"`       // 币种
	OrderFee     int64  `xml:"order_fee"`      // 子单应付金额，分
	Duty         int64  `xml:"duty"`           // 关税，分
	TransportFee int64  `xml:"transport_fee"`  // 物流费用，分
	ProductFee   int64  `xml:"product_fee"`    // 商品费用，分
	Explanation  string `xml:"explanation"`    // 申报结果说明，如果状态是失败或异常，显示失败原因
}

// DeclareOrder2 统一下单.
func DeclareQuery2(clt *core.Client, req *DeclareQueryRequest) (resp *DeclareQueryResponse, err error) {
	m2, err := DeclareQuery(clt, req.FieldsMap())
	if err != nil {
		return nil, err
	}

	resp = &DeclareQueryResponse{
		TransactionId: m2["transaction_id"],
		Body:          m2[""], // 返回原文默认用空字符串指向
	}
	if str := m2["count"]; str != "" {
		if n, err := strconv.ParseInt(str, 10, 64); err != nil {
			err = fmt.Errorf("parse count:%q to int64 failed: %s", str, err.Error())
			return nil, err
		} else {
			resp.Count = int(n)
		}
	}
	resp.QueryList = make([]QueryItem, resp.Count)
	for i := 0; i < resp.Count; i++ {
		resp.QueryList[i].Customs = m2["customs_"+strconv.Itoa(i)]
		resp.QueryList[i].State = m2["state_"+strconv.Itoa(i)]
		resp.QueryList[i].ModifyTime = m2["modify_time_"+strconv.Itoa(i)]
		resp.QueryList[i].CertCheckResult = m2["cert_check_result_"+strconv.Itoa(i)]
		resp.QueryList[i].SubOrderNo = m2["sub_order_no_"+strconv.Itoa(i)]
		resp.QueryList[i].SubOrderId = m2["sub_order_id_"+strconv.Itoa(i)]
		resp.QueryList[i].MchCustomsNo = m2["mch_customs_no_"+strconv.Itoa(i)]
		resp.QueryList[i].FeeType = m2["fee_type_"+strconv.Itoa(i)]
		resp.QueryList[i].Explanation = m2["explanation_"+strconv.Itoa(i)]

		if str := m2["order_fee_"+strconv.Itoa(i)]; str != "" {
			if n, err := strconv.ParseInt(str, 10, 64); err != nil {
				err = fmt.Errorf("parse order_fee_%d:%q to int64 failed: %s", i, str, err.Error())
				return nil, err
			} else {
				resp.QueryList[i].OrderFee = n
			}
		}

		if str := m2["duty_"+strconv.Itoa(i)]; str != "" {
			if n, err := strconv.ParseInt(str, 10, 64); err != nil {
				err = fmt.Errorf("parse duty_%d:%q to int64 failed: %s", i, str, err.Error())
				return nil, err
			} else {
				resp.QueryList[i].Duty = n
			}
		}

		if str := m2["transport_fee_"+strconv.Itoa(i)]; str != "" {
			if n, err := strconv.ParseInt(str, 10, 64); err != nil {
				err = fmt.Errorf("parse transport_fee_%d:%q to int64 failed: %s", i, str, err.Error())
				return nil, err
			} else {
				resp.QueryList[i].TransportFee = n
			}
		}

		if str := m2["product_fee_"+strconv.Itoa(i)]; str != "" {
			if n, err := strconv.ParseInt(str, 10, 64); err != nil {
				err = fmt.Errorf("parse product_fee_%d:%q to int64 failed: %s", i, str, err.Error())
				return nil, err
			} else {
				resp.QueryList[i].ProductFee = n
			}
		}
	}
	return resp, nil
}
