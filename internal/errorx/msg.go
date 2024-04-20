package errorx

const (
	Error                 = "未知错误!"
	SuccessMsg            = "Success!"
	ErrorParamsInvalidMsg = "参数错误!"
)

var MsgFlags = map[int]string{
	Err:              Error,
	Success:          SuccessMsg,
	ErrParamsInvalid: ErrorParamsInvalidMsg,
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[Err]
}
