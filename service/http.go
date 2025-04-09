package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/iotames/qrbridge/db"
	dbtable "github.com/iotames/qrbridge/dbtable"
	"github.com/iotames/qrbridge/util"
)

func GetOneQrcode(qrid *int, qrToUrl *string, code string) error {
	qrcode := dbtable.Qrcode{}
	querySQL := fmt.Sprintf("select id, to_url from %s where code = $1", qrcode.TableName())
	return db.GetOne(querySQL, []interface{}{qrid, qrToUrl}, code)
}

// UpdateQrcode 更新数据库的二维码查询记录
// TODO 待优化，应使用事务，把更新二维码表和二维码查询日志表的操作放在一个事务中
func UpdateQrcode(r http.Request, toUrl string, status int, isNew bool, codeParsed string) {
	lg := util.GetLogger()
	code := r.URL.Query().Get("code")
	requestIp := util.GetHttpClientIP(r)
	userAgent := r.Header.Get("User-Agent")
	// 将请求头转换为JSON字符串
	requestHeaders, err := json.Marshal(r.Header)
	if err != nil {
		lg.Errorf("convert headers to json failed: %v", err)
		requestHeaders = []byte("{}")
	}

	// 更新二维码表
	qr := dbtable.Qrcode{}
	qrTable := qr.TableName()
	if isNew {
		cols := []string{"code", "to_url", "code_parsed", "pv", "status"}
		colsVal := []interface{}{code, toUrl, codeParsed, 1, status}
		db.ExecInsert(qrTable, cols, colsVal)
	} else {
		updateSQL := fmt.Sprintf("update %s set pv = pv + 1, updated_at = CURRENT_TIMESTAMP where code = '%s'", qrTable, code)
		db.ExecSqlText(updateSQL)
	}
	requestUrl := r.URL.String()

	// 插入二维码查询日志表
	qrlog := dbtable.QrcodeQueryLog{}
	cols := []string{"code", "request_url", "user_agent", "request_headers", "request_ip"}
	colsVal := []interface{}{code, requestUrl, userAgent, string(requestHeaders), requestIp}
	db.ExecInsert(qrlog.TableName(), cols, colsVal)
	// lg.Debugf("---UpdateQrcode--requestUrl(%s)-code: %s, request_ip: %s, user_agent(%s)---hdr(%s)", requestUrl, code, requestIp, userAgent, string(requestHeaders))
}
