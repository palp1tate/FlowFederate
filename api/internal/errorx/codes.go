package errorx

const (
	Err     = 500
	Success = 200

	ErrParamsInvalid    = 10001
	ErrMustAdmin        = 10002
	ErrInternalServer   = 10003
	ErrGetServerList    = 10004
	ErrGetClientList    = 10005
	ErrGetTaskList      = 10006
	ErrGetTask          = 10007
	ErrTokenNeed        = 10008
	ErrTokenExpired     = 10009
	ErrTokenParseFailed = 10010
	ErrMobileExists     = 10011
	ErrRegisterFailed   = 10012
	ErrGenTokenFailed   = 10013
	ErrSendSmsFailed    = 10014
	ErrCaptchaExpired   = 10015
	ErrCaptchaIncorrect = 10016
	TooManyRequests     = 10017
	ErrMobileNotFound   = 10018
)
