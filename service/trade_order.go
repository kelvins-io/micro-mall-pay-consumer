package service

import (
	"context"
	"fmt"
	"gitee.com/cristiane/micro-mall-pay-consumer/model/args"
	"gitee.com/cristiane/micro-mall-pay-consumer/model/mysql"
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
	"strings"
	"time"
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
	payRecordWhere := map[string]interface{}{
		"tx_id":     notice.PayId,
		"pay_state": 3, // 完成支付
	}
	time.Sleep(3 * time.Second) // 支付事务先提交

	payRecordList, err := repository.GetPayRecordList("id,tx_id,out_trade_no,amount", payRecordWhere)
	if err != nil {
		kelvins.ErrLogger.Info(ctx, "GetPayRecordList err: %v, where: %+v", err, payRecordWhere)
		return err
	}
	if len(payRecordList) == 0 {
		// 查询不到数据
		return fmt.Errorf("GetPayRecordList payRecordList empty tx_id: %v", notice.PayId)
	}

	go func() {
		// 检查用户身份
		userName, err := checkUserIdentity(ctx, notice)
		if err != nil {
			return
		}
		// 通知用户
		err = noticeUserPayResult(ctx, userName, payRecordList)
		if err != nil {
			return
		}
	}()

	// 通知订单服务支付结果
	go func() {
		err = noticeOrderPayCallback(ctx, notice)
		if err != nil {
			return
		}
	}()

	return nil
}

func noticeOrderPayCallback(ctx context.Context, notice args.TradePayNotice) error {
	serverName := args.RpcServiceMicroMallOrder
	conn, err := util.GetGrpcClient(ctx, serverName)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetGrpcClient %v,err: %v", serverName, err)
		return err
	}
	//defer conn.Close()
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

func noticeUserPayResult(ctx context.Context, userName string, recordList []mysql.PayRecord) error {
	total := decimal.NewFromInt(0)
	var orderCode strings.Builder
	for i := 0; i < len(recordList); i++ {
		orderCode.WriteString(fmt.Sprintf("【订单号：%v，金额：%v】", recordList[i].OutTradeNo, recordList[i].Amount))
		amount, err := decimal.NewFromString(recordList[i].Amount)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "noticeUserPayResult decimal  err: %v, value: %v", err, recordList[i].Amount)
			return fmt.Errorf(errcode.GetErrMsg(code.ErrDecimalParse))
		}
		total = util.DecimalAdd(total, amount)
	}

	// 邮件通知
	emailNotice := fmt.Sprintf(args.TradePayEmailTemp, userName, orderCode.String(), total.String())
	for _, receiver := range vars.EmailNoticeSetting.Receivers {
		err := email.SendEmailNotice(ctx, receiver, vars.AppName, emailNotice)
		if err != nil {
			kelvins.ErrLogger.Info(ctx, "noticeUserPayResult SendEmailNotice err, emailNotice: %v", emailNotice)
		}
	}
	return nil
}

func checkUserIdentity(ctx context.Context, notice args.TradePayNotice) (userName string, err error) {
	// 获取用户信息
	serverName := args.RpcServiceMicroMallUsers
	conn, err := util.GetGrpcClient(ctx, serverName)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetGrpcClient %v,err: %v", serverName, err)
		return "", err
	}
	//defer conn.Close()
	client := users.NewUsersServiceClient(conn)
	r := users.CheckUserStateRequest{
		UidList: []int64{notice.Uid},
	}
	userInfo, err := client.CheckUserState(ctx, &r)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "CheckUserState %v,err: %v, r: %+v", serverName, err, notice.Uid)
		return "", err
	}
	if userInfo.Common.Code == users.RetCode_SUCCESS {
		userInfoRsp, err := client.GetUserInfo(ctx, &users.GetUserInfoRequest{Uid: notice.Uid})
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "GetUserInfo %v,err: %v, r: %+v", serverName, err, notice.Uid)
			return "", err
		}
		if userInfoRsp.Common.Code == users.RetCode_SUCCESS {
			return userInfoRsp.Info.UserName, nil
		}
		return "", fmt.Errorf(userInfo.Common.Msg)
	}

	kelvins.ErrLogger.Errorf(ctx, "GetUserInfo  err:%v, rsp: %+v", err, userInfo)
	return "", fmt.Errorf(userInfo.Common.Msg)
}

func TradePayConsumeErr(ctx context.Context, errMsg, body string) error {
	return nil
}
