package startup

import "gitee.com/cristiane/micro-mall-pay-consumer/service"

const (
	TaskNameTradePayNotice    = "task_trade_pay_notice"
	TaskNameTradePayNoticeErr = "task_trade_pay_notice_err"
)

func GetNamedTaskFuncs() map[string]interface{} {

	var taskRegister = map[string]interface{}{
		TaskNameTradePayNotice:    service.TradePayConsume,
		TaskNameTradePayNoticeErr: service.TradePayConsumeErr,
	}
	return taskRegister
}
