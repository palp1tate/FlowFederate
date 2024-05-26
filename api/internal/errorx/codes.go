package errorx

const (
	Err     = 500
	Success = 200

	ErrParamsInvalid      = 10001
	ErrMustAdmin          = 10002
	ErrInternalServer     = 10003
	ErrGetServerList      = 10004
	ErrGetClientList      = 10005
	ErrGetTaskList        = 10006
	ErrGetTask            = 10007
	ErrTokenNeed          = 10008
	ErrTokenExpired       = 10009
	ErrTokenParseFailed   = 10010
	ErrMobileExists       = 10011
	ErrRegisterFailed     = 10012
	ErrGenTokenFailed     = 10013
	ErrSendSmsFailed      = 10014
	ErrCaptchaExpired     = 10015
	ErrCaptchaIncorrect   = 10016
	TooManyRequests       = 10017
	ErrMobileNotFound     = 10018
	ErrPasswordIncorrect  = 10019
	ErrPicCaptchaGenerate = 10020
	ErrUploadFileFailed   = 10021
	ErrDeleteFileFailed   = 10022
	ErrUserNotFound       = 10023
	ErrPasswordReset      = 10024
	ErrUpdateUserFailed   = 10025
	ErrAdminNotFound      = 10026
	ErrUpdateAdminFailed  = 10027
	ErrGetUserListFailed  = 10028
	ErrDeleteUserFailed   = 10029
)
