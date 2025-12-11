package webserver

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/iotames/easyserver/httpsvr"
	"github.com/iotames/easyserver/response"
	"github.com/iotames/qrbridge/conf"
)

func setHandler(svr *httpsvr.EasyServer) {
	svr.AddHandler("GET", "/hello", hello)
	svr.AddHandler("GET", "/qrcode", qrcode)
	svr.AddHandler("GET", "/codetest"+strconv.Itoa(conf.EncryptAdd), codetest)

	// 成本核价占比
	svr.AddHandler("GET", "/pricing_percent", pricing_percent)
	// PO导入
	svr.AddHandler("POST", "/api/poimport", poimport)
}

func hello(ctx httpsvr.Context) {
	ctx.Writer.Write(response.NewApiDataOk("hello api").Bytes())
}

func postJsonValue(ctx httpsvr.Context, v any) error {
	// 读取请求体中的数据
	var err error
	var b []byte
	b, err = io.ReadAll(ctx.Request.Body)
	if err != nil {
		return fmt.Errorf("读取请求体失败io.ReadAll error: %w", err)
	}
	// 解析JSON数据
	err = json.Unmarshal(b, v)
	if err != nil {
		return fmt.Errorf("解析JSON失败json.Unmarshal error: %w", err)
	}
	return err
}
