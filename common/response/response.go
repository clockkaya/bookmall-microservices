package response

import (
	"net/http"

	"bookNew/common/xerr" // 引入刚才定义的 xerr

	"github.com/zeromicro/go-zero/rest/httpx"
)

type Body struct {
	Code uint32      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func Response(w http.ResponseWriter, resp interface{}, err error) {
	var body Body

	// 1. 如果 Logic 层返回了错误
	if err != nil {
		switch e := err.(type) {
		case *xerr.CodeError: // 如果是自定义的业务错误
			body.Code = e.GetErrCode()
			body.Msg = e.GetErrMsg()
		default: // 如果是系统未知的错误 (如 panic, 数据库连接失败等)
			body.Code = xerr.SERVER_COMMON_ERROR
			body.Msg = xerr.MapErrMsg(xerr.SERVER_COMMON_ERROR)
		}
	} else {
		// 2. 如果成功
		body.Code = xerr.OK
		body.Msg = "success"
		body.Data = resp
	}

	// 3. 统一返回 HTTP 200，具体错误看 Body.Code
	httpx.OkJson(w, body)
}
