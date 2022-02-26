package tool

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
	"translate/model"
)

func init() { //设置timeout
	http.DefaultClient.Timeout = 1 * time.Second
}

//读取文件数据
func GetText(src string) (res string, err error) {
	data, err := ioutil.ReadFile(src)
	res = string(data)
	return
}

//发送请求
func GetRes(text string) (res string) {
	url := GetUrl(text)
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return fmt.Sprint(err)
	}
	defer resp.Body.Close()
	data := model.Data{}
	json.NewDecoder(resp.Body).Decode(&data)
	if data.Error_code != model.SUCCESS && data.Error_code != 0 {
		return fmt.Sprintf("失败 错误码: %d", data.Error_code)
	}
	for _, v := range data.Trans_result {
		res = strings.Join([]string{res, v.String()}, "\n")
	}
	return
}

//结果另存为,暂未应用
// func outfile(data string, dst string) (bool, error) {
// 	if err := ioutil.WriteFile(dst, []byte(data), 0666); err != nil {
// 		return false, err
// 	}
// 	return true, nil
// }
// func Deal(src string, dst string) (bool, error) {
// 	text, err := GetText(src)
// 	if err != nil {
// 		log.Println(err)
// 		return false, err
// 	}
// 	data := GetRes(text)
// 	return outfile(data, dst)
// }

//通过文件名获取翻译
func Deal(path string) (string, string, bool) {
	buf, err := GetText(path)
	if err != nil {
		return "", err.Error(), false
	}
	return buf, Deal_(buf), true
}

//通过文本获取翻译
func Deal_(src string) (res string) {
	done := make(chan struct{})
	go func() {
		res = GetRes(src)
		close(done)
	}()
	select {
	case <-done:
		return
	case <-time.After(1200 * time.Millisecond):
		return "网络延迟大，失败"
	}

}
