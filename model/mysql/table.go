package mysql

import (
	"time"
)

const (
	TablePayRecord = "pay_record"
)

type PayRecord struct {
	Id          int64     `xorm:"pk autoincr comment('自增ID') BIGINT"`
	TxId        string    `xorm:"not null comment('批次交易号') index CHAR(40)"`
	OutTradeNo  string    `xorm:"not null comment('外部商户订单号') index CHAR(40)"`
	NotifyUrl   string    `xorm:"comment('交易结果通知地址') VARCHAR(255)"`
	Description string    `xorm:"comment('交易描述') VARCHAR(255)"`
	Merchant    string    `xorm:"not null comment('交易商户ID') index CHAR(40)"`
	Attach      string    `xorm:"comment('交易留言') TEXT"`
	User        string    `xorm:"not null comment('交易用户ID') index CHAR(40)"`
	Amount      string    `xorm:"not null comment('交易数量') DECIMAL(64,4)"`
	CoinType    int       `xorm:"not null default 0 comment('交易币种，0-cny,1-usd') TINYINT"`
	Reduction   string    `xorm:"comment('满减优惠') DECIMAL(64,4)"`
	PayType     int       `xorm:"not null comment('交易类型，1入账，2退款') TINYINT"`
	PayState    int       `xorm:"comment('支付状态，0-未支付，1-支付中，2-支付失败，3-支付成功') TINYINT"`
	CreateTime  time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('创建时间') DATETIME"`
	UpdateTime  time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('修改时间') DATETIME"`
}
