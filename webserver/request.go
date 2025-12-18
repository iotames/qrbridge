package webserver

import (
	"fmt"

	"github.com/iotames/easyserver/httpsvr"
)

// checkJsonField 检查JSON数据中是否缺少必填字段，并检查数据类型。
// 目前仅支持int、string类型。int不允许0，string不允许空字符串。
func checkJsonField(ctx httpsvr.Context, fields ...string) (requestData map[string]any, err error) {
	// var err error
	var ok bool
	var val any
	// 解析JSON数据
	// var requestData map[string]string
	err = postJsonValue(ctx, &requestData)
	if err != nil {
		// ctx.Writer.Write(response.NewApiDataServerError(err.Error()).Bytes())
		return
	}
	for _, field := range fields {
		val, ok = requestData[field]
		if !ok {
			return nil, fmt.Errorf("缺少字段：%s", field)
		}
		// 判断val类型，如果是字符串类型，字符串不能为空
		switch v := val.(type) {
		case string:
			if v == "" {
				return nil, fmt.Errorf("字段%s不能为空", field)
			}
		case int:
			if v == 0 {
				return nil, fmt.Errorf("字段%s不能为0", field)
			}
		default:
			return nil, fmt.Errorf("字段%s类型错误", field)
		}
	}
	return requestData, nil
}
