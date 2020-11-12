package repository

import (
	"gitee.com/cristiane/micro-mall-pay-consumer/model/mysql"
	"gitee.com/kelvins-io/kelvins"
)

func GetPayRecordList(selectSql string, where interface{}) ([]mysql.PayRecord, error) {
	var result = make([]mysql.PayRecord, 0)
	var err error
	err = kelvins.XORM_DBEngine.Table(mysql.TablePayRecord).Select(selectSql).Where(where).Find(&result)
	return result, err
}
