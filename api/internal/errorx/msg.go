package errorx

const (
	Error                    = "未知错误!"
	SuccessMsg               = "Success!"
	ErrorParamsInvalidMsg    = "参数错误!"
	ErrorMustAdminMsg        = "需要管理员权限!"
	ErrorGetTaskListMsg      = "获取任务列表失败!"
	ErrorGetTaskMsg          = "获取任务失败!"
	ErrorInternalServerMsg   = "内部服务器错误!"
	ErrorGetServerListMsg    = "获取边缘服务器列表失败!"
	ErrorGetClientListMsg    = "获取客户端列表失败!"
	ErrorTokenNeedMsg        = "需要Token!"
	ErrorTokenExpiredMsg     = "Token已过期!"
	ErrorTokenParseFailedMsg = "Token解析失败!"
	ErrorMobileExistsMsg     = "该手机号已被注册!"
	ErrorRegisterFailedMsg   = "注册失败!"
	ErrorGenTokenFailedMsg   = "生成Token失败!"
	ErrorSendSmsFailedMsg    = "发送短信失败!"
	ErrorCaptchaExpiredMsg   = "验证码已过期!"
	ErrorCaptchaIncorrectMsg = "验证码错误!"
	ErrorTooManyRequestsMsg  = "请求过于频繁!"
	ErrorMobileNotFoundMsg   = "该手机号未注册!"
)

var MsgFlags = map[int]string{
	Err:                 Error,
	Success:             SuccessMsg,
	ErrParamsInvalid:    ErrorParamsInvalidMsg,
	ErrMustAdmin:        ErrorMustAdminMsg,
	ErrGetTaskList:      ErrorGetTaskListMsg,
	ErrGetTask:          ErrorGetTaskMsg,
	ErrInternalServer:   ErrorInternalServerMsg,
	ErrGetServerList:    ErrorGetServerListMsg,
	ErrGetClientList:    ErrorGetClientListMsg,
	ErrTokenNeed:        ErrorTokenNeedMsg,
	ErrTokenExpired:     ErrorTokenExpiredMsg,
	ErrTokenParseFailed: ErrorTokenParseFailedMsg,
	ErrMobileExists:     ErrorMobileExistsMsg,
	ErrRegisterFailed:   ErrorRegisterFailedMsg,
	ErrGenTokenFailed:   ErrorGenTokenFailedMsg,
	ErrSendSmsFailed:    ErrorSendSmsFailedMsg,
	ErrCaptchaExpired:   ErrorCaptchaExpiredMsg,
	ErrCaptchaIncorrect: ErrorCaptchaIncorrectMsg,
	TooManyRequests:     ErrorTooManyRequestsMsg,
	ErrMobileNotFound:   ErrorMobileNotFoundMsg,
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[Err]
}
