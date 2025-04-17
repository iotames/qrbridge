package webserver

import (
	"strconv"

	"github.com/iotames/easyserver/httpsvr"
	"github.com/iotames/easyserver/response"
	"github.com/iotames/miniutils"
	"github.com/iotames/qrbridge/conf"
	"github.com/iotames/qrbridge/service"
	"github.com/iotames/qrbridge/util"
)

func setHandler(svr *httpsvr.EasyServer) {
	svr.AddHandler("GET", "/hello", hello)
	svr.AddHandler("GET", "/qrcode", qrcode)
	// 成本核价占比
	svr.AddHandler("GET", "/pricing_percent", pricing_percent)
	svr.AddHandler("GET", "/codetest"+strconv.Itoa(conf.EncryptAdd), codetest)
}

func hello(ctx httpsvr.Context) {
	ctx.Writer.Write(response.NewApiDataOk("hello api").Bytes())
}

func qrcode(ctx httpsvr.Context) {
	lg := util.GetLogger()
	code := ctx.Request.URL.Query().Get("code")
	decoder := util.NewUrlEncrypt("mid", conf.EncryptMultiple, conf.EncryptAdd)
	codeParsed, err := decoder.Decrypt(code)
	var qrid int
	var qrToUrl string
	queryErr := service.GetOneQrcode(&qrid, &qrToUrl, code)
	if queryErr != nil {
		// 数据库查询失败
		ctx.Writer.Write(response.NewApiDataServerError(queryErr.Error()).Bytes())
		return
	}
	isNew := true
	if qrid > 0 {
		isNew = false
	}
	var status int
	var toUrl string
	if err != nil {
		// 解密失败。在数据库中查询code. status = 0. 生成新的数据或更新数据
		status = 0
		service.UpdateQrcode(*ctx.Request, qrToUrl, status, isNew, "")
		if qrToUrl == "" {
			ctx.Writer.Write(response.NewApiDataQueryArgsError("非法请求").Bytes())
			return
		}
		toUrl = miniutils.GetUrl(qrToUrl, conf.ToBaseUrl)
	} else {
		// 解密成功。在数据库中查询code. status = 1. 生成新数据或更新数据
		status = 1
		toUrl = "?code=" + code
		service.UpdateQrcode(*ctx.Request, toUrl, status, isNew, codeParsed)
		toUrl = miniutils.GetUrl(toUrl, conf.ToBaseUrl)
	}
	lg.Debugf("---qrcode---code(%s)---codeParsed(%s)--", code, codeParsed)
	if toUrl != "" {
		// 重定向
		ctx.Writer.Header().Set("Location", toUrl)
		ctx.Writer.WriteHeader(302)
		return
	}
	ctx.Writer.Write(response.NewApiDataQueryArgsError("未知请求错误").Bytes())
}

func codetest(ctx httpsvr.Context) {
	// 获取整个query参数，并转成字符串
	query := ctx.Request.URL.Query()
	mid := query.Get("mid")
	if mid == "" {
		ctx.Writer.Write(response.NewApiDataQueryArgsError("缺少mid参数").Bytes())
		return
	}
	queryStr := query.Encode()
	encoder := util.NewUrlEncrypt("mid", conf.EncryptMultiple, conf.EncryptAdd)
	code := encoder.Encrypt(queryStr)
	ctx.Writer.Write(response.NewApiDataOk("code=" + code).Bytes())
}
