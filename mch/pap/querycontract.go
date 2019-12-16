package pap

import (
	"gopkg.in/skOak/wechat.v2/mch/core"
	"strconv"
)

// QueryContract查询签约关系
func QueryContract(clt *core.Client, req map[string]string) (resp map[string]string, err error) {
	return clt.PostXML("/papay/querycontract", req)
}

// 文档：https://pay.weixin.qq.com/wiki/doc/api/pap.php?chapter=18_2&index=7
type QueryContractRequest struct {
	XMLName struct{} `xml:"xml" json:"-"`

	// 必选参数(contract_id和plan_id+contract_code二选一)
	ContractId   string `xml:"contract_id"`   // 委托代扣签约成功后微信返回的委托代扣协议id
	PlanId       int64  `xml:"plan_id"`       // 协议模板id
	ContractCode string `xml:"contract_code"` // 签约协议号
	// Version      string `xml:"version"`       // 固定值1.0
}

func (req *QueryContractRequest) FieldsMap() map[string]string {
	m1 := make(map[string]string, 4)
	if req.ContractId != "" {
		m1["contract_id"] = req.ContractId
	} else {
		m1["plan_id"] = strconv.FormatInt(req.PlanId, 10)
		m1["contract_code"] = req.ContractCode
	}
	m1["version"] = "1.0"

	return m1
}

type QueryContractResponse struct {
	XMLName struct{} `xml:"xml" json:"-"`
	Body    string   `xml:"-" json:"-"` // 返回内容原文，主要是为了记录日志

	// 必选返回
	ContractId             string `xml:"contract_id"`              // 委托代扣签约成功后微信返回的委托代扣协议id
	PlanId                 int64  `xml:"plan_id"`                  // 协议模板id
	ContractCode           string `xml:"contract_code"`            // 签约协议号
	RequestSerial          int64  `xml:"request_serial"`           // 商户请求签约时的序列号，要求唯一性。序列号主要用于排序，不作为查询条件。
	ContractState          int    `xml:"contract_state"`           // 协议状态,0已签约;1未签约;9签约进行中
	ContractSignedTime     string `xml:"contract_signed_time"`     // 协议签署时间,2015-07-01 10:00:00
	ContractExpiredTime    string `xml:"contract_expired_time"`    // 协议到期时间,2015-07-01 10:00:00
	ContractTerminatedTime string `xml:"contract_terminated_time"` // 协议解约时间,2015-07-01 10:00:00
	// 协议解约方式，仅state为1时有效;0未解约;1有效期过自动解约;2用户主动解约;3商户API解约;4商户平台解约;5注销
	ContractTerminationMode   int    `xml:"contract_termination_mode"`
	ContractTerminationRemark string `xml:"contract_termination_remark"` // 解约备注；state为1时有效
	Openid                    string `xml:"openid"`                      // 商户appid下用户的唯一标识
	ContractDisplayAccount    string `xml:"contract_display_account"`    // 签约用户的名称，用于页面展示
}

// QueryContract2 统一下单.
func QueryContract2(clt *core.Client, req *QueryContractRequest) (resp *QueryContractResponse, err error) {
	m2, err := QueryContract(clt, req.FieldsMap())
	if err != nil {
		return nil, err
	}

	resp = &QueryContractResponse{
		ContractId:                m2["contract_id"],
		ContractCode:              m2["contract_code"],
		ContractSignedTime:        m2["contract_signed_time"],
		ContractExpiredTime:       m2["contract_expired_time"],
		ContractTerminatedTime:    m2["contract_terminated_time"],
		ContractTerminationRemark: m2["contract_termination_remark"],
		Openid:                 m2["openid"],
		ContractDisplayAccount: m2["contract_display_account"],
		Body: m2[""], // 返回原文默认用空字符串指向
	}
	resp.PlanId, _ = strconv.ParseInt(m2["plan_id"], 10, 64)
	resp.RequestSerial, _ = strconv.ParseInt(m2["request_serial"], 10, 64)
	contractState, _ := strconv.ParseInt(m2["contract_state"], 10, 64)
	resp.ContractState = int(contractState)
	contractTerminationMode, _ := strconv.ParseInt(m2["contract_termination_mode"], 10, 64)
	resp.ContractTerminationMode = int(contractTerminationMode)
	return resp, nil
}
