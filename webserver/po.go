package webserver

import (
	"fmt"
	"slices"
	"strings"

	"github.com/iotames/easyserver/httpsvr"
	"github.com/iotames/easyserver/response"
	"github.com/iotames/qrbridge/biz"
)

func poimportCheck(ctx httpsvr.Context) (inputtpl string, inputfile string, outputfile string, err error) {
	var ok bool
	// 解析JSON数据
	var requestData map[string]string
	err = postJsonValue(ctx, &requestData)
	if err != nil {
		ctx.Writer.Write(response.NewApiDataServerError(err.Error()).Bytes())
		return
	}
	// 获取inputtpl字段
	inputtpl, ok = requestData["inputtpl"]
	if !ok || inputtpl == "" {
		// ctx.Writer.Write(response.NewApiDataFail("inputtpl字段错误", 400).Bytes())
		err = fmt.Errorf("inputtpl字段错误")
		ctx.Writer.Write(response.NewApiDataQueryArgsError(err.Error()).Bytes())
		return
	}
	// 获取inputfile字段
	inputfile, ok = requestData["inputfile"]
	if !ok || inputfile == "" {
		err = fmt.Errorf("inputfile字段错误")
		ctx.Writer.Write(response.NewApiDataQueryArgsError(err.Error()).Bytes())
		return
	}
	// 获取outputfile字段
	outputfile, ok = requestData["outputfile"]
	if !ok || outputfile == "" {
		err = fmt.Errorf("outputfile字段错误")
		ctx.Writer.Write(response.NewApiDataQueryArgsError(err.Error()).Bytes())
		return
	}
	return
}

func poimport(ctx httpsvr.Context) {
	inputtpl, inputfile, outputfile, err := poimportCheck(ctx)
	if err != nil {
		return
	}
	// 打印inputfile字段
	fmt.Printf("接收到的inputfile(%s); outputfile(%s)\n", inputfile, outputfile)
	tpllist := []string{"Rohnisch"}
	if !slices.Contains(tpllist, inputtpl) {
		err = fmt.Errorf("inputtpl参数错误: 仅支持(%s)", strings.Join(tpllist, ","))
		ctx.Writer.Write(response.NewApiDataQueryArgsError(err.Error()).Bytes())
		return
	}
	_, err = biz.PoFileTransform(inputtpl, inputfile, outputfile)
	if err != nil {
		ctx.Writer.Write(response.NewApiDataServerError(err.Error()).Bytes())
		return
	}
	ctx.Writer.Write(response.NewApiData(response.JsonObject{"inputfile": inputfile, "outputfile": outputfile}, "success", 200).Bytes())
}
