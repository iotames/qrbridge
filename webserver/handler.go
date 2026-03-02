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

type WebPage struct {
	UrlPath               string
	Title, AmisPageConfig string
}

func getWebPageList() []WebPage {
	return []WebPage{
		{"/user/po-import", "客户PO文件转换", "/api/po-import-page"},
		{"/user/cmd", "快捷命令", "/api/cmd-page"},
	}
}

func getAmisPageTitle(urlpath string) string {
	for _, p := range getWebPageList() {
		if p.UrlPath == urlpath {
			return p.Title
		}
	}
	return ""
}

func setHandler(svr *httpsvr.EasyServer) {
	svr.AddHandler("GET", "/", home)
	svr.AddHandler("GET", "/hello", hello)
	svr.AddHandler("GET", "/qrcode", qrcode)
	svr.AddHandler("GET", "/codetest"+strconv.Itoa(conf.EncryptAdd), codetest)

	// 成本核价占比
	svr.AddHandler("GET", "/pricing_percent", pricing_percent)

	// 统一Page路由配置
	for _, p := range getWebPageList() {
		svr.AddHandler("GET", p.UrlPath, func(ctx httpsvr.Context) {
			data := map[string]interface{}{
				"title":            p.Title,
				"amis_page_config": p.AmisPageConfig + "?" + ctx.Request.URL.RawQuery,
			}
			SetContentByTplFile("tpl/amis.html", ctx.Writer, data)
		})
	}

	// PO导入
	// svr.AddHandler("GET", "/user/po-import", pagePoimport)
	svr.AddHandler("POST", "/api/poimport", poimport)
	svr.AddHandler("POST", "/api/potransform", potransform)
	svr.AddHandler("POST", "/api/uploadfile", uploadfile)
	svr.AddHandler("GET", "/api/customer/list", customerList)
	svr.AddHandler("GET", "/api/po-import-page", getAmisPoImportPage)

	// 执行系统CMD命令
	// svr.AddHandler("GET", "/user/cmd", pageCmdHome)
	svr.AddHandler("GET", "/api/cmd-page", getAmisCmdConfig)
	svr.AddPostHandler("/api/cmd/exec", execmd)
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

func home(ctx httpsvr.Context) {
	ctx.Writer.Write(response.NewApiDataOk("hello home").Bytes())
}

func hello(ctx httpsvr.Context) {
	ctx.Writer.Write(response.NewApiDataOk("hello api").Bytes())
}
