package service

import (
	"context"
	"fmt"
	"gitee.com/cristiane/micro-mall-pay-consumer/model/args"
	"gitee.com/cristiane/micro-mall-pay-consumer/pkg/code"
	"gitee.com/cristiane/micro-mall-pay-consumer/pkg/util"
	"gitee.com/cristiane/micro-mall-pay-consumer/pkg/util/email"
	"gitee.com/cristiane/micro-mall-pay-consumer/proto/micro_mall_logistics_proto/logistics_business"
	"gitee.com/cristiane/micro-mall-pay-consumer/proto/micro_mall_order_proto/order_business"
	"gitee.com/cristiane/micro-mall-pay-consumer/proto/micro_mall_users_proto/users"
	"gitee.com/cristiane/micro-mall-pay-consumer/repository"
	"gitee.com/cristiane/micro-mall-pay-consumer/vars"
	"gitee.com/kelvins-io/common/errcode"
	"gitee.com/kelvins-io/common/json"
	"gitee.com/kelvins-io/kelvins"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc"
	"time"
)

const (
	tradePayEmailTemp = "微商城已经收到【%v】您的订单【%v】，本次您共消费：%s元，欢迎你随时关注订单状态。"
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
	userInfo, err := client.GetUserInfo(ctx, &r)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetUserInfo %v,err: %v, r: %+v", serverName, err, r)
		return err
	}
	if userInfo.Common.Code != users.RetCode_SUCCESS {
		kelvins.ErrLogger.Errorf(ctx, "GetUserInfo %v,not ok : %v, rsp: %+v", serverName, err, userInfo)
		return fmt.Errorf(userInfo.Common.Msg)
	}
	if userInfo.Info == nil || userInfo.Info.AccountId == "" {
		kelvins.ErrLogger.Errorf(ctx, "GetUserInfo %v,accountId null : %v, rsp: %+v", serverName, err, userInfo)
		return fmt.Errorf(errcode.GetErrMsg(code.UserNotExist))
	}

	// 从数据查询支付订单
	wherePayRecord := map[string]interface{}{
		"tx_id": notice.TxCode,
	}
	recordList, _, err := repository.GetPayRecordList("*", wherePayRecord, nil, nil, 0, 0)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetPayRecordList err: %v, where: %v", err, wherePayRecord)
		return fmt.Errorf(errcode.GetErrMsg(code.ErrorServer))
	}
	total := decimal.NewFromInt(0)
	for i := 0; i < len(recordList); i++ {
		amount, err := decimal.NewFromString(recordList[i].Amount)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "decimal  err: %v, value: %v", err, recordList[i].Amount)
			return fmt.Errorf(errcode.GetErrMsg(code.ErrDecimalParse))
		}
		total = util.DecimalAdd(total, amount)
	}
	// 获取订单
	serverName = args.RpcServiceMicroMallOrder
	conn, err = util.GetGrpcClient(serverName)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetGrpcClient %v,err: %v", serverName, err)
		return err
	}
	defer conn.Close()
	orderClient := order_business.NewOrderBusinessServiceClient(conn)
	reqOrderSku := order_business.GetOrderSkuRequest{TxCode: notice.TxCode}
	orderSkuRsp, err := orderClient.GetOrderSku(ctx, &reqOrderSku)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetOrderSku %v,err: %v, r: %+v", orderSkuRsp, err, reqOrderSku)
		return err
	}
	if orderSkuRsp == nil || orderSkuRsp.Common == nil || orderSkuRsp.Common.Code == order_business.RetCode_ERROR {
		kelvins.ErrLogger.Errorf(ctx, "GetOrderSku %v,err: %v, r: %+v", orderSkuRsp, err, reqOrderSku)
		return err
	}
	for i := 0; i < len(orderSkuRsp.OrderList); i++ {
		row := orderSkuRsp.OrderList[i]
		// 关联物流消息
		serverName = args.RpcServiceMicroMallLogistics
		conn, err = util.GetGrpcClient(serverName)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "GetGrpcClient %v,err: %v", serverName, err)
			return err
		}
		closeGRPCConn(conn)
		clientLogistics := logistics_business.NewLogisticsBusinessServiceClient(conn)
		goods := make([]*logistics_business.GoodsInfo, len(row.Goods))
		for k := 0; k < len(row.Goods); k++ {
			g := &logistics_business.GoodsInfo{
				SkuCode: row.Goods[k].SkuCode,
				Name:    row.Goods[k].Name,
				Kind:    row.Goods[k].Name,
				Count:   row.Goods[k].Amount,
			}
			goods[k] = g
		}
		reqLogistics := logistics_business.ApplyLogisticsRequest{
			OutTradeNo:  row.OrderCode,
			Courier:     "微商城快递",
			CourierType: 1,
			ReceiveType: 1,
			SendTime:    util.ParseTimeOfStr(time.Now().Unix()),
			Customer: &logistics_business.CustomerInfo{
				SendUser:     fmt.Sprintf("%v", notice.Uid),
				SendAddr:     "深圳市宝安区宝源二区",
				SendPhone:    "18319430520",
				SendTime:     util.ParseTimeOfStr(time.Now().Unix()),
				ReceiveUser:  "张起灵",
				ReceiveAddr:  "南海国人民南路128号枫林别院12栋10单元1009",
				ReceivePhone: "15501707783",
			},
			Goods: goods,
		}
		applyRsp, err := clientLogistics.ApplyLogistics(ctx, &reqLogistics)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "ApplyLogistics %v,err: %v, r: %+v", applyRsp, err, reqLogistics)
			return err
		}
		if applyRsp == nil || applyRsp.Common == nil || applyRsp.Common.Code == logistics_business.RetCode_ERROR {
			kelvins.ErrLogger.Errorf(ctx, "ApplyLogistics %v,err: %v, r: %+v", applyRsp, err, reqLogistics)
			return err
		}
		kelvins.BusinessLogger.Infof(ctx, "ApplyLogistics code: %v", applyRsp.LogisticsCode)
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

func closeGRPCConn(conn *grpc.ClientConn) {
	defer conn.Close()
}

func TradePayConsumeErr(ctx context.Context, errMsg, body string) error {
	return nil
}
