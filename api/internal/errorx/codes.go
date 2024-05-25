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
)
