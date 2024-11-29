package App

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

const CodeSuccess = 0
const CodeError = 1
const CodeToken = 2
const CodeSign = 3
const CodeValidate = 4
const CodeSql = 5

type ResponseJson struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func ReturnJson(r *ghttp.Request, code gcode.Code) {
	returnDetail := code.Detail()
	if g.IsEmpty(returnDetail) {
		returnDetail = make([]interface{}, 0)
	}
	r.Response.WriteJsonExit(ResponseJson{
		Code: code.Code(),
		Msg:  code.Message(),
		Data: returnDetail,
	})
}
