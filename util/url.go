package util

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// URL请求参数加解密
type UrlEncrypt struct {
	multiple     int
	addint       int
	encryptField string
}

// 构建Url参数加密器
// encryptField: 要加密的字段名，值必须是整数
// multiple: 放大倍数，推荐3倍
// add: 加的整数。20050324
func NewUrlEncrypt(encryptField string, multiple, add int) *UrlEncrypt {
	if multiple == 0 || add == 0 {
		panic("multiple and add must be greater than 0")
	}
	if encryptField == "" {
		panic("encryptField must not be empty")
	}
	return &UrlEncrypt{
		multiple:     multiple,
		addint:       add,
		encryptField: encryptField,
	}
}

func (u UrlEncrypt) encryptInt(before int) int {
	return before*u.multiple + u.addint
}

func (u UrlEncrypt) decryptInt(after int) int {
	return (after - u.addint) / u.multiple
}

func (u UrlEncrypt) Encrypt(urlQuery string) string {
	// 解析查询参数
	values, _ := url.ParseQuery(urlQuery)
	// 处理mid参数
	if intstr := values.Get(u.encryptField); intstr != "" {
		if num, err := strconv.Atoi(intstr); err == nil {
			values.Set(u.encryptField, strconv.Itoa(u.encryptInt(num)))
		}
	}

	// 生成查询字符串
	processedQuery := values.Encode()

	// URL安全的Base64编码
	urlBase64 := base64.URLEncoding.EncodeToString([]byte(processedQuery))
	return strings.TrimRight(urlBase64, "=")
}

func (u UrlEncrypt) Decrypt(urlBase64 string) (string, error) {
	// strings.HasPrefix(encrypted, "code=")
	// 补全Base64可能缺少的等号
	if pad := len(urlBase64) % 4; pad != 0 {
		urlBase64 += strings.Repeat("=", 4-pad)
	}

	// URL安全的Base64解码
	decoded, err := base64.URLEncoding.DecodeString(urlBase64)
	if err != nil {
		return "", err
	}
	querystr := string(decoded)

	// 解析查询参数
	values, err := url.ParseQuery(querystr)
	if err != nil {
		return "", err
	}
	intstr := values.Get(u.encryptField)
	if intstr == "" {
		return "", fmt.Errorf("not found %s in %s", u.encryptField, querystr)
	}
	num, err := strconv.Atoi(intstr)
	if err != nil {
		return "", err
	}
	values.Set(u.encryptField, strconv.Itoa(u.decryptInt(num)))
	return values.Encode(), nil
}
