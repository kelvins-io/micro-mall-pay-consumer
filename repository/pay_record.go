package repository

import (
	"gitee.com/cristiane/micro-mall-pay-consumer/model/mysql"
	"gitee.com/kelvins-io/kelvins"
)

func GetPayRecordList(selectSql string, where interface{}, orderByDesc, orderByAsc []string, pageSize, pageNum int) ([]mysql.PayRecord, int64, error) {
	if selectSql == "" {
		selectSql = "*"
	}
	var result = make([]mysql.PayRecord, 0)
	var err error
	var total int64
	session := kelvins.XORM_DBEngine.Table(mysql.TablePayRecord).
		Select(selectSql).
		Where(where).
		Desc(orderByDesc...).
		Asc(orderByAsc...)
	if pageSize > 0 && pageNum >= 1 {
		session = session.Limit(pageSize, (pageNum-1)*pageSize)
	}
	total, err = session.FindAndCount(&result)
	return result, total, err
}
