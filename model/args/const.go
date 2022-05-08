package args

type MerchantsMaterialInfo struct {
	Uid          int64
	MaterialId   int64
	RegisterAddr string
	HealthCardNo string
	Identity     int32
	State        int32
	TaxCardNo    string
}

const (
	RpcServiceMicroMallUsers = "micro-mall-users"
	RpcServiceMicroMallOrder = "micro-mall-order"
)

const (
	TradePayEmailTemp = "尊敬的用户【%v】你好，您的订单【%v】已经完成支付，本次共消费：%s元，请留意订单状态。"
)

type CommonBusinessMsg struct {
	Type    int    `json:"type"`
	Tag     string `json:"tag"`
	UUID    string `json:"uuid"`
	Content string `json:"content"`
}

type TradePayNotice struct {
	Uid    int64  `json:"uid"`
	Time   string `json:"time"`
	TxCode string `json:"tx_code"`
	PayId  string `json:"pay_id"`
}

const (
	Unknown                 = 0
	TradePayEventTypeCreate = 10016
	TradePayEventTypeExpire = 10017
)

var MsgFlags = map[int]string{
	Unknown:                 "未知",
	TradePayEventTypeCreate: "支付订单创建",
	TradePayEventTypeExpire: "支付订单过期",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[Unknown]
}
