package service

import (
	"context"
	"fmt"
	"gitee.com/cristiane/micro-mall-pay-consumer/model/args"
	"gitee.com/cristiane/micro-mall-pay-consumer/pkg/code"
	"gitee.com/cristiane/micro-mall-pay-consumer/pkg/util"
	"gitee.com/cristiane/micro-mall-pay-consumer/proto/micro_mall_users_proto/users"
	"gitee.com/cristiane/micro-mall-pay-consumer/repository"
	"gitee.com/kelvins-io/common/errcode"
	"gitee.com/kelvins-io/common/json"
	"gitee.com/kelvins-io/kelvins"
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
	// 获取用户信息
	serverName := args.RpcServiceMicroMallUsers
	conn, err := util.GetGrpcClient(serverName)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetGrpcClient %v,err: %v", serverName, err)
		return err
	}
	defer conn.Close()
	client := users.NewUsersServiceClient(conn)
	r := users.GetUserInfoRequest{
		Uid: notice.Uid,
	}
	rsp, err := client.GetUserInfo(ctx, &r)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetUserInfo %v,err: %v, r: %+v", serverName, err, r)
		return err
	}
	if rsp.Common.Code != users.RetCode_SUCCESS {
		kelvins.ErrLogger.Errorf(ctx, "GetUserInfo %v,not ok : %v, rsp: %+v", serverName, err, rsp)
		return fmt.Errorf(rsp.Common.Msg)
	}
	if rsp.Info == nil || rsp.Info.AccountId == "" {
		kelvins.ErrLogger.Errorf(ctx, "GetUserInfo %v,accountId null : %v, rsp: %+v", serverName, err, rsp)
		return fmt.Errorf(errcode.GetErrMsg(code.UserNotExist))
	}

	// 从数据查询
	wherePayRecord := map[string]interface{}{}
	recordList, _, err := repository.GetPayRecordList("*", wherePayRecord, nil, nil, 0, 0)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetPayRecordList err: %v, where: %v", err, wherePayRecord)
		return fmt.Errorf(errcode.GetErrMsg(code.ErrorServer))
	}
	amount :=
	for i := 0; i < len(recordList); i++ {

	}

	return nil
}

func TradePayConsumeErr(ctx context.Context, errMsg, body string) error {
	return nil
}
