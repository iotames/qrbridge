package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	// "github.com/iotames/qrbridge/db"
	"github.com/iotames/qrbridge/db"
	dbtable "github.com/iotames/qrbridge/dbtable"
	"github.com/iotames/qrbridge/util"
)

func GetOneQrcode(r http.Request) (qrcode dbtable.Qrcode, err error) {
	code := r.URL.Query().Get("code")
	err = db.GetOne(fmt.Sprintf("select * from %s where code = $1", qrcode.TableName()), &qrcode, code)
	return qrcode, err
}
func UpdateQrcode(r http.Request, isNew bool) {
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

	// 更新数据库中的数据
	qr := dbtable.Qrcode{}
	if isNew {
		// TODO
		db.ExecInsert(qr.TableName(), []string{}, []interface{}{code, r.URL.Query().Get("to_url"), 1, 1, requestIp, userAgent, string(requestHeaders)})
	} else {
		// TODO
		// db.ExecUpdate(qr.TableName(), "pv = pv + 1", "code = $1", code)
	}

	// 插入数据
	qrlg := dbtable.QrcodeQueryLog{}
	db.ExecInsert(qrlg.TableName(), []string{}, []interface{}{code, requestIp, userAgent, string(requestHeaders)})

	lg.Debugf("code: %s, request_ip: %s, user_agent(%s)---hdr(%s)", code, requestIp, userAgent, string(requestHeaders))

}
