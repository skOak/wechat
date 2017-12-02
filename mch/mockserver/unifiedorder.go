package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	"gopkg.in/skOak/wechat.v2/mch/core"
	"gopkg.in/skOak/wechat.v2/mch/pay"
	"gopkg.in/skOak/wechat.v2/util"
)

type UnifiedOrderReq struct {
	RequestCommon
	pay.UnifiedOrderRequest
}

// 字段解释参考：https://pay.weixin.qq.com/wiki/doc/api/app/app.php?chapter=9_7&index=3
type PayResultNotify struct {
	ReturnCode string `xml:"return_code,omitempty"`
	ReturnMsg  string `xml:"return_msg,omitempty"`

	AppId       string `xml:"appid,omitempty"`
	MchId       string `xml:"mch_id,omitempty"`
	DeviceInfo  string `xml:"device_info,omitempty"`
	NonceStr    string `xml:"nonce_str,omitempty"`
	Sign        string `xml:"sign,omitempty"`
	ResultCode  string `xml:"result_code,omitempty"`
	ErrCode     string `xml:"err_code,omitempty"`
	ErrCodeDes  string `xml:"err_code_des,omitempty"`
	OpenId      string `xml:"openid,omitempty"`
	IsSubscribe string `xml:"is_subscribe,omitempty"`
	TradeType   string `xml:"trade_type,omitempty"`
	BankType    string `xml:"bank_type,omitempty"`
	TotalFee    int    `xml:"total_fee,omitempty"`
	FeeType     string `xml:"fee_type,omitempty"`
	CashFee     int    `xml:"cash_fee,omitempty"`
	CashFeeType string `xml:"cash_fee_type,omitempty"`
	CouponFee   int    `xml:"coupon_fee,omitempty"`
	//CouponCount int    `xml:"coupon_count,omitempty"`
	//CouponIdN    string `xml:"coupon_id_n,omitempty"`
	//CouponFeeN   int    `xml:"coupon_fee_n,omitempty"`
	TransactionId string `xml:"trasaction_id,omitempty"`
	OutTradeNo    string `xml:"out_trade_no,omitempty"`
	Attach        string `xml:"attach,omitempty"`
	TimeEnd       string `xml:"time_end,omitempty"`
}

func UnifiedOrder(w http.ResponseWriter, r *http.Request) {
	req := &UnifiedOrderReq{}
	resp := make(map[string]string)
	defer func() {
		// omit marshal error
		rb, _ := xml.Marshal(resp)
		w.Write(rb)
	}()
	err := xml.NewDecoder(r.Body).Decode(req)
	if err != nil {
		// Decode Error
		resp["return_code"] = "FAIL"
		resp["return_msg"] = fmt.Sprintf("Decode Error:%v", err.Error())
		return
	}

	fieldsMap := req.UnifiedOrderRequest.FieldsMap()
	for k, v := range req.RequestCommon.FieldsMap() {
		fieldsMap[k] = v
	}
	reqSign := req.Sign
	signWant := ""
	signType := req.SignType
	switch signType {
	case core.SignType_HMAC_SHA256:
		signWant = core.Sign2(fieldsMap, ApiKey, hmac.New(sha256.New, []byte(ApiKey)))
	default:
		signWant = core.Sign2(fieldsMap, ApiKey, md5.New())
	}
	if reqSign != signWant {
		// Sign Check Error
		resp["return_code"] = "FAIL"
		resp["return_msg"] = "Signature Incorrect"
		return
	}

	resp["return_code"] = "SUCCESS"
	resp["appid"] = req.AppId
	resp["mch_id"] = req.MchId
	resp["device_info"] = req.DeviceInfo
	resp["nonce_str"] = util.NonceStr()
	resp["result_code"] = "SUCCESS"
	resp["trade_type"] = req.TradeType
	resp["prepay_id"] = "wx" + core.FormatTime(time.Now()) + util.NonceStr()
	switch signType {
	case core.SignType_HMAC_SHA256:
		resp["sign"] = core.Sign2(resp, ApiKey, hmac.New(sha256.New, []byte(ApiKey)))
	default:
		resp["sign"] = core.Sign2(resp, ApiKey, md5.New())
	}
	go func() {
		// 发送交易结果通知
		<-time.After(time.Second * 10)
		// Notify client in 10 seconds
		notify := PayResultNotify{
			ReturnCode:    "SUCCESS",
			AppId:         req.AppId,
			MchId:         req.MchId,
			DeviceInfo:    req.DeviceInfo,
			NonceStr:      util.NonceStr(),
			ResultCode:    "SUCCESS",
			OpenId:        "wxopenid" + util.NonceStr(),
			IsSubscribe:   "Y",
			TradeType:     req.TradeType,
			BankType:      "CMB_CREDIT",
			TotalFee:      int(req.TotalFee),
			FeeType:       req.FeeType,
			CashFee:       0,
			TransactionId: "wxtr_" + util.NonceStr(),
			OutTradeNo:    req.OutTradeNo,
			Attach:        req.Attach,
			TimeEnd:       core.FormatTime(time.Now()),
		}
		switch signType {
		case core.SignType_HMAC_SHA256:
			notify.Sign = core.Sign2(util.StructFieldsMap(notify, "xml"), ApiKey, hmac.New(sha256.New, []byte(ApiKey)))
		default:
			notify.Sign = core.Sign2(util.StructFieldsMap(notify, "xml"), ApiKey, md5.New())
		}
		nb, _ := xml.Marshal(notify)
		Notify(req.NotifyURL, "application/xml", bytes.NewBuffer(nb), 0)
	}()
	return
}
