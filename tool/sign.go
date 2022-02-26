package tool

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"time"
)

const (
	Appid      = "xxxx"
	Key        = "xxxx"
	PreRequest = "https://fanyi-api.baidu.com/api/trans/vip/translate?"
)

var (
	Salt string
	from string
	to   string
)

func init() {
	from = "auto"
	to = "zh"
}
func SetFrom(s string) {
	from = s
}
func SetTo(s string) {
	to = s
}

//生成签名
func GetSign(text *string) string {
	rand.Seed(time.Now().UnixMilli())
	Salt = fmt.Sprintf("%d", rand.Uint32())
	sign := md5.Sum([]byte(strings.Join([]string{Appid, *text, Salt, Key}, "")))
	//获取32位MD5字串
	return hex.EncodeToString(sign[:])
}

//拼接参数
func Combine(text *string) string {
	sign := GetSign(text)
	query := url.QueryEscape(*text)
	return fmt.Sprintf("q=%s&from=%s&to=%s&appid=%s&salt=%s&sign=%s", query, from, to, Appid, Salt, sign)
}

//生成请求
func GetUrl(text string) string {
	return strings.Join([]string{PreRequest, Combine(&text)}, "")
}
