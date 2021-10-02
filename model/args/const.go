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
	RpcServiceMicroMallUsers       = "micro-mall-users"
	RpcServiceMicroMallOrder       = "micro-mall-order"
	RpcServiceMicroMallLogistics   = "micro-mall-logistics"
	RpcServiceMicroMallUserTrolley = "micro-mall-trolley"
)

const (
	TradePayEmailTemp = "微商城已经收到【%v】您的订单 %v 已经完成支付，本次您共消费：%s元，欢迎你随时关注订单状态。"
)

const (
	TaskNameTradePayNotice    = "task_trade_pay_notice"
	TaskNameTradePayNoticeErr = "task_trade_pay_notice_err"
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
