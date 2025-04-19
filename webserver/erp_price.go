package webserver

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/iotames/easyserver/httpsvr"
	"github.com/iotames/easyserver/response"
	"github.com/iotames/qrbridge/db"
	"github.com/iotames/qrbridge/service"
	"github.com/iotames/qrbridge/util"
)

func pricing_percent(ctx httpsvr.Context) {
	sqlText, err := service.GetSqlText("pricing_percent.sql")
	if err != nil {
		// 数据库查询SQL获取失败
		ctx.Writer.Write(response.NewApiDataServerError(err.Error()).Bytes())
		return
	}
	sqlText = strings.ReplaceAll(sqlText, "---%s---", "%s")
	customer_name := ctx.Request.URL.Query().Get("customer_name")
	clothing_material_type_detail_name := ctx.Request.URL.Query().Get("clothing_material_type_detail_name")
	// TODO 微服务暂时不考虑SQL注入的情况
	var whereList []string
	var queryArgs []interface{}
	var where1 []string
	if customer_name != "" {
		for i, v := range strings.Split(customer_name, ",") {
			where1 = append(where1, fmt.Sprintf("$%d", i+1))
			queryArgs = append(queryArgs, v)
		}
		whereList = append(whereList, fmt.Sprintf("and cp.customer_name in (%s)", strings.Join(where1, ",")))
	}
	if clothing_material_type_detail_name != "" {
		var where2 []string
		where2add := len(where1) + 1
		for i, v := range strings.Split(clothing_material_type_detail_name, ",") {
			where2 = append(where2, fmt.Sprintf("$%d", i+where2add))
			queryArgs = append(queryArgs, v)
		}
		whereList = append(whereList, fmt.Sprintf("and cp.clothing_material_type_detail_name in(%s)", strings.Join(where2, ",")))
	}
	queryStr := strings.Join(whereList, " ")
	sqlText = fmt.Sprintf(sqlText, queryStr)
	//  ".$where1.$where2."
	// $where1 = " and cp.customer_name in('".$sCustomerShortName."') ";
	// $where2 = " and cp.clothing_material_type_detail_name in('".$sStyleArchivesMaterialTypeDetailName."') ";

	lg := util.GetLogger()
	lg.Debugf("------pricing_percent---sqlText(%s)---queryArgs(%v)--", sqlText, queryArgs)

	data := make(map[string]interface{}, 21)
	// querySQL := fmt.Sprintf("select id, to_url from %s where code = $1", qrcode.TableName())
	queryErr := db.GetOneData(sqlText, &data, queryArgs...)
	if queryErr != nil {
		// 数据库查询失败
		ctx.Writer.Write(response.NewApiDataServerError(queryErr.Error()).Bytes())
		return
	}

	ctx.Writer.Write(response.NewApiData(base64Decode(data), "success", 200).Bytes())
}

// base64Decode
// data 参数类型实际是 map[string][]byte
func base64Decode(data map[string]interface{}) map[string]interface{} {
	// 创建新的map存储解码后的数据
	result := make(map[string]interface{})
	// 遍历原始数据
	for key, value := range data {
		// 将interface{}转换为string
		if strValue, ok := value.([]byte); ok {
			// 对base64编码的字符串进行解码
			strValue := string(strValue)
			val, errconv := strconv.ParseFloat(strValue, 64)
			if errconv != nil {
				result[key] = strValue
			} else {
				result[key] = val
			}
			// fmt.Printf("-----解码-key(%s)---value(%v)--(%T)--\n", key, result[key], result[key])
		} else {
			// fmt.Printf("-----解码--Fail-key(%s)---(%v)\n", key, value)
			// 如果不是字符串类型，保留原始值
			result[key] = value
		}
	}
	return result
}
