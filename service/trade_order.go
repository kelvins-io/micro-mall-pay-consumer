package service

import (
	"context"
	"fmt"
	"gitee.com/cristiane/micro-mall-pay-consumer/model/args"
	"gitee.com/cristiane/micro-mall-pay-consumer/pkg/code"
	"gitee.com/cristiane/micro-mall-pay-consumer/pkg/util"
	"gitee.com/cristiane/micro-mall-pay-consumer/pkg/util/email"
	"gitee.com/cristiane/micro-mall-pay-consumer/proto/micro_mall_order_proto/order_business"
	"gitee.com/cristiane/micro-mall-pay-consumer/proto/micro_mall_users_proto/users"
	"gitee.com/cristiane/micro-mall-pay-consumer/repository"
	"gitee.com/cristiane/micro-mall-pay-consumer/vars"
	"gitee.com/kelvins-io/common/errcode"
	"gitee.com/kelvins-io/common/json"
	"gitee.com/kelvins-io/kelvins"
	"github.com/shopspring/decimal"
)

const (
	tradePayEmailTemp = "微商城已经收到【%v】您的订单【%v】已经完成支付，本次您共消费：%s元，欢迎你随时关注订单状态。"
)

func TradePayConsume(ctx context.Context, body string) error {
	// 通知消息解码
	var businessMsg args.CommonBusinessMsg
	var err error
	err = json.Unmarshal(body, &businessMsg)
	if err != nil {
		kelvins.ErrLogger.Info(ctx, "body:%v Unmarshal err: %v", body, err)
		return err
	}
	if businessMsg.Type != args.TradePayEventTypeCreate {
		return fmt.Errorf(errcode.GetErrMsg(code.NoticeTypeNotEqual))
	}
	var notice args.TradePayNotice
	err = json.Unmarshal(businessMsg.Msg, &notice)
	if err != nil {
		kelvins.ErrLogger.Info(ctx, "businessMsg.Msg: %v Unmarshal err: %v", businessMsg.Msg, err)
		return err
	}
	// 检查用户身份
	userInfo, err := checkUserIdentity(ctx, notice)
	if err != nil {
		return err
	}
	// 通知用户
	err = noticeUserPayResult(ctx, notice, userInfo)
	if err != nil {
		return err
	}
	// 通知订单服务支付结果
	err = noticeOrderPayCallback(ctx, notice, userInfo)
	if err != nil {
		return err
	}
	return nil
}

func noticeOrderPayCallback(ctx context.Context, notice args.TradePayNotice, userInfo *users.GetUserInfoResponse) error {
	serverName := args.RpcServiceMicroMallOrder
	conn, err := util.GetGrpcClient(serverName)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetGrpcClient %v,err: %v", serverName, err)
		return err
	}
	defer conn.Close()
	client := order_business.NewOrderBusinessServiceClient(conn)
	req := order_business.OrderTradeNoticeRequest{
		Uid:         notice.Uid,
		OrderTxCode: notice.TxCode,
		PayId:       notice.PayId,
	}
	orderTradeNoticeRsp, err := client.OrderTradeNotice(ctx, &req)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetUserInfo %v,err: %v, r: %+v", serverName, err, req)
		return err
	}
	if orderTradeNoticeRsp.Common.Code != order_business.RetCode_SUCCESS {
		kelvins.ErrLogger.Errorf(ctx, "GetUserInfo %v,err: %v, rsp: %+v", serverName, err, orderTradeNoticeRsp)
		return fmt.Errorf(errcode.GetErrMsg(code.ErrorServer))
	}
	return nil
}

func noticeUserPayResult(ctx context.Context, notice args.TradePayNotice, userInfo *users.GetUserInfoResponse) error {
	// 从数据查询支付订单
	wherePayRecord := map[string]interface{}{
		"tx_id":     notice.PayId,
		"pay_state": 3, // 完成支付
	}
	recordList, err := repository.GetPayRecordList("amount", wherePayRecord)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetPayRecordList err: %v, where: %v", err, wherePayRecord)
		return fmt.Errorf(errcode.GetErrMsg(code.ErrorServer))
	}
	// todo  有毒？？
	fmt.Println("recordList len==", len(recordList))
	total := decimal.NewFromInt(0)
	for i := 0; i < len(recordList); i++ {
		amount, err := decimal.NewFromString(recordList[i].Amount)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "decimal  err: %v, value: %v", err, recordList[i].Amount)
			return fmt.Errorf(errcode.GetErrMsg(code.ErrDecimalParse))
		}
		total = util.DecimalAdd(total, amount)
	}
	// 邮件通知
	msgNotice := fmt.Sprintf(tradePayEmailTemp, userInfo.Info.UserName, notice.TxCode, total.String())
	err = email.SendEmailNotice(ctx, "565608463@qq.com", vars.AppName, msgNotice)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "SendEmailNotice err：%v, content：%v", err, msgNotice)
		return err
	}

	return nil
}

func checkUserIdentity(ctx context.Context, notice args.TradePayNotice) (*users.GetUserInfoResponse, error) {
	// 获取用户信息
	serverName := args.RpcServiceMicroMallUsers
	conn, err := util.GetGrpcClient(serverName)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetGrpcClient %v,err: %v", serverName, err)
		return nil, err
	}
	defer conn.Close()
	client := users.NewUsersServiceClient(conn)
	r := users.GetUserInfoRequest{
		Uid: notice.Uid,
	}
	userInfo, err := client.GetUserInfo(ctx, &r)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetUserInfo %v,err: %v, r: %+v", serverName, err, r)
		return nil, err
	}
	if userInfo.Common.Code != users.RetCode_SUCCESS {
		kelvins.ErrLogger.Errorf(ctx, "GetUserInfo %v,not ok : %v, rsp: %+v", serverName, err, userInfo)
		return nil, fmt.Errorf(userInfo.Common.Msg)
	}
	if userInfo.Info == nil || userInfo.Info.AccountId == "" {
		kelvins.ErrLogger.Errorf(ctx, "GetUserInfo %v,accountId null : %v, rsp: %+v", serverName, err, userInfo)
		return nil, fmt.Errorf(errcode.GetErrMsg(code.UserNotExist))
	}
	return userInfo, nil
}

func TradePayConsumeErr(ctx context.Context, errMsg, body string) error {
	return nil
}
