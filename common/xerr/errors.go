package xerr

import "fmt"

type CodeError struct {
	errCode uint32
	errMsg  string
}

// 1. 获取错误码
func (e *CodeError) GetErrCode() uint32 {
	return e.errCode
}

// 2. 获取错误信息
func (e *CodeError) GetErrMsg() string {
	return e.errMsg
}

// 3. 实现 error 接口
func (e *CodeError) Error() string {
	return fmt.Sprintf("ErrCode:%d，ErrMsg:%s", e.errCode, e.errMsg)
}

// 4. 工厂方法：创建一个新的业务异常
func NewErrCode(errCode uint32) *CodeError {
	return &CodeError{errCode: errCode, errMsg: MapErrMsg(errCode)}
}

// 5. 工厂方法：创建一个带自定义消息的业务异常
func NewErrMsg(errMsg string) *CodeError {
	return &CodeError{errCode: SERVER_COMMON_ERROR, errMsg: errMsg}
}

// 新增：专门用于返回 JSON 的结构体
type CodeErrorResponse struct {
	Code uint32 `json:"code"`
	Msg  string `json:"msg"`
}

// 新增：给 CodeError 加个方法，方便转换
func (e *CodeError) Data() *CodeErrorResponse {
	return &CodeErrorResponse{
		Code: e.errCode,
		Msg:  e.errMsg,
	}
}
