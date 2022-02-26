package model

import "fmt"

const SUCCESS = 52000

//百度api要求的数据格式
type Data struct {
	From         string   `json:"from"`
	To           string   `json:"to"`
	Trans_result []Result `json:"trans_result"`
	Error_code   int      `json:"error_code"`
}
type Result struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

func (r *Result) String() string {
	return fmt.Sprintf("%s\n", r.Dst)
}
