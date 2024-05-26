package errorx

const (
	Error                      = "未知错误!"
	SuccessMsg                 = "Success!"
	ErrorParamsInvalidMsg      = "参数错误!"
	ErrorMustAdminMsg          = "需要管理员权限!"
	ErrorGetTaskListMsg        = "获取任务列表失败!"
	ErrorGetTaskMsg            = "获取任务失败!"
	ErrorInternalServerMsg     = "内部服务器错误!"
	ErrorGetServerListMsg      = "获取边缘服务器列表失败!"
	ErrorGetClientListMsg      = "获取客户端列表失败!"
	ErrorTokenNeedMsg          = "需要Token!"
	ErrorTokenExpiredMsg       = "Token已过期!"
	ErrorTokenParseFailedMsg   = "Token解析失败!"
	ErrorMobileExistsMsg       = "该手机号已被注册!"
	ErrorRegisterFailedMsg     = "注册失败!"
	ErrorGenTokenFailedMsg     = "生成Token失败!"
	ErrorSendSmsFailedMsg      = "发送短信失败!"
	ErrorCaptchaExpiredMsg     = "验证码已过期!"
	ErrorCaptchaIncorrectMsg   = "验证码错误!"
	ErrorTooManyRequestsMsg    = "请求过于频繁!"
	ErrorMobileNotFoundMsg     = "该手机号未注册!"
	ErrorPasswordIncorrectMsg  = "密码错误!"
	ErrorPicCaptchaGenerateMsg = "生成图片验证码失败!"
	ErrorUploadFileFailedMsg   = "上传文件失败!"
	ErrorDeleteFileFailedMsg   = "删除文件失败!"
	ErrorUserNotFoundMsg       = "用户不存在!"
	ErrorPasswordResetMsg      = "密码重置失败!"
	ErrorUpdateUserFailedMsg   = "更新用户信息失败!"
	ErrorAdminNotFoundMsg      = "管理员不存在!"
	ErrorUpdateAdminFailedMsg  = "更新管理员信息失败!"
	ErrorGetUserListFailedMsg  = "获取用户列表失败!"
	ErrorDeleteUserFailedMsg   = "删除用户失败!"
)

var MsgFlags = map[int]string{
	Err:                   Error,
	Success:               SuccessMsg,
	ErrParamsInvalid:      ErrorParamsInvalidMsg,
	ErrMustAdmin:          ErrorMustAdminMsg,
	ErrGetTaskList:        ErrorGetTaskListMsg,
	ErrGetTask:            ErrorGetTaskMsg,
	ErrInternalServer:     ErrorInternalServerMsg,
	ErrGetServerList:      ErrorGetServerListMsg,
	ErrGetClientList:      ErrorGetClientListMsg,
	ErrTokenNeed:          ErrorTokenNeedMsg,
	ErrTokenExpired:       ErrorTokenExpiredMsg,
	ErrTokenParseFailed:   ErrorTokenParseFailedMsg,
	ErrMobileExists:       ErrorMobileExistsMsg,
	ErrRegisterFailed:     ErrorRegisterFailedMsg,
	ErrGenTokenFailed:     ErrorGenTokenFailedMsg,
	ErrSendSmsFailed:      ErrorSendSmsFailedMsg,
	ErrCaptchaExpired:     ErrorCaptchaExpiredMsg,
	ErrCaptchaIncorrect:   ErrorCaptchaIncorrectMsg,
	TooManyRequests:       ErrorTooManyRequestsMsg,
	ErrMobileNotFound:     ErrorMobileNotFoundMsg,
	ErrPasswordIncorrect:  ErrorPasswordIncorrectMsg,
	ErrPicCaptchaGenerate: ErrorPicCaptchaGenerateMsg,
	ErrUploadFileFailed:   ErrorUploadFileFailedMsg,
	ErrDeleteFileFailed:   ErrorDeleteFileFailedMsg,
	ErrUserNotFound:       ErrorUserNotFoundMsg,
	ErrPasswordReset:      ErrorPasswordResetMsg,
	ErrUpdateUserFailed:   ErrorUpdateUserFailedMsg,
	ErrAdminNotFound:      ErrorAdminNotFoundMsg,
	ErrUpdateAdminFailed:  ErrorUpdateAdminFailedMsg,
	ErrGetUserListFailed:  ErrorGetUserListFailedMsg,
	ErrDeleteUserFailed:   ErrorDeleteUserFailedMsg,
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[Err]
}
