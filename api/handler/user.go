package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/redis/go-redis/v9"

	"github.com/palp1tate/go-crypto-guard/rsa"

	"github.com/palp1tate/FlowFederate/api/internal/model"

	"github.com/palp1tate/FlowFederate/api/dao"

	"github.com/palp1tate/FlowFederate/api/internal/utils"

	"github.com/palp1tate/FlowFederate/api/types"

	"github.com/palp1tate/FlowFederate/api/global"
	"github.com/palp1tate/FlowFederate/api/internal/consts"
	"github.com/palp1tate/FlowFederate/api/internal/errorx"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	registerForm := types.RegisterForm{}
	if err := c.ShouldBind(&registerForm); err != nil {
		HandleValidatorError(c, err)
		return
	}
	key := fmt.Sprintf("%d-%s", consts.Register, registerForm.Mobile)
	smsCode, err := global.RedisClient.Get(context.Background(), key).Result()
	if errors.Is(err, redis.Nil) {
		HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrCaptchaExpired)
		return
	} else {
		if smsCode != registerForm.Captcha {
			HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrCaptchaIncorrect)
			return
		}
		global.RedisClient.Del(context.Background(), key)
	}
	encodedPassword, _ := pwd.GenRSA(registerForm.Password, "publicKey.pem")
	user := model.User{
		Username: registerForm.Username,
		Mobile:   registerForm.Mobile,
		Password: encodedPassword,
		Role:     consts.User,
	}
	if err := dao.NewUserDao(context.Background()).CreateUser(&user); err != nil {
		HandleHttpResponse(c, http.StatusInternalServerError, errorx.ErrRegisterFailed)
		return
	}
	claims := utils.CustomClaims{
		ID:   strconv.Itoa(user.ID),
		Role: strconv.Itoa(consts.User),
	}
	j := utils.NewJWT(global.ServerConfig.JWT.SigningKey, global.ServerConfig.JWT.Expiration)
	token, err := j.CreateToken(claims)
	if err != nil {
		HandleHttpResponse(c, http.StatusInternalServerError, errorx.ErrGenTokenFailed)
		return
	}
	HandleHttpResponse(c, http.StatusOK, errorx.Success, token)
}

//func UserLogin(c *gin.Context) {
//	loginForm := form.LoginForm{}
//	if err := c.ShouldBind(&loginForm); err != nil {
//		HandleValidatorError(c, err)
//		return
//	}
//	var token string
//	switch loginForm.Type {
//	case consts.LoginByPassword:
//		_, err := global.ThirdPartyServiceClient.CheckPicCaptcha(context.Background(), &thirdPb.CheckPicCaptchaRequest{
//			CaptchaId: loginForm.CaptchaId,
//			Captcha:   loginForm.Captcha,
//		})
//		if err != nil {
//			HandleGrpcErrorToHttp(c, err)
//			return
//		}
//		res, err := global.UserServiceClient.LoginByPassword(context.Background(), &userPb.LoginByPasswordRequest{
//			Mobile:   loginForm.Mobile,
//			Password: loginForm.Password,
//		})
//		if err != nil {
//			HandleGrpcErrorToHttp(c, err)
//			return
//		}
//		claims := util.CustomClaims{
//			ID:   int(res.Id),
//			Role: consts.User,
//		}
//		j := util.NewJWT(global.ServerConfig.JWT.SigningKey, global.ServerConfig.JWT.Expiration)
//		token, err = j.CreateToken(claims)
//		if err != nil {
//			HandleHttpResponse(c, http.StatusInternalServerError, errorx.ErrGenTokenFailed, nil, nil)
//			return
//		}
//	case consts.LoginByCaptcha:
//		_, err := global.ThirdPartyServiceClient.CheckSmsCaptcha(context.Background(), &thirdPb.CheckSmsCaptchaRequest{
//			Mobile:  loginForm.Mobile,
//			Captcha: loginForm.Captcha,
//			Type:    int64(consts.Login),
//		})
//		if err != nil {
//			HandleGrpcErrorToHttp(c, err)
//			return
//		}
//		res, err := global.UserServiceClient.LoginBySMS(context.Background(), &userPb.LoginBySMSRequest{
//			Mobile: loginForm.Mobile,
//		})
//		if err != nil {
//			HandleGrpcErrorToHttp(c, err)
//			return
//		}
//		claims := util.CustomClaims{
//			ID:   int(res.Id),
//			Role: consts.User,
//		}
//		j := util.NewJWT(global.ServerConfig.JWT.SigningKey, global.ServerConfig.JWT.Expiration)
//		token, err = j.CreateToken(claims)
//		if err != nil {
//			HandleHttpResponse(c, http.StatusInternalServerError, errorx.ErrGenTokenFailed, nil, nil)
//			return
//		}
//	}
//	HandleHttpResponse(c, http.StatusOK, errorx.Success, token, nil)
//	return
//}
//
//func GetUser(c *gin.Context) {
//	token := c.GetString("token")
//	res, err := global.UserServiceClient.GetUser(NewCtxWithToken(context.Background(), token), &empty.Empty{})
//	if err != nil {
//		HandleGrpcErrorToHttp(c, err)
//		return
//	}
//	j := util.NewJWT(global.ServerConfig.JWT.SigningKey, global.ServerConfig.JWT.Expiration)
//	refreshedToken, _ := j.RefreshToken(token)
//	HandleHttpResponse(c, http.StatusOK, errorx.Success, refreshedToken, res)
//	return
//}
//
//func ResetPassword(c *gin.Context) {
//	resetPasswordForm := form.ResetPasswordForm{}
//	if err := c.ShouldBind(&resetPasswordForm); err != nil {
//		HandleValidatorError(c, err)
//		return
//	}
//	_, err := global.ThirdPartyServiceClient.CheckSmsCaptcha(context.Background(), &thirdPb.CheckSmsCaptchaRequest{
//		Mobile:  resetPasswordForm.Mobile,
//		Captcha: resetPasswordForm.Captcha,
//		Type:    int64(consts.ResetPassword),
//	})
//	if err != nil {
//		HandleGrpcErrorToHttp(c, err)
//		return
//	}
//	_, err = global.UserServiceClient.ResetPassword(context.Background(), &userPb.ResetPasswordRequest{
//		Mobile:   resetPasswordForm.Mobile,
//		Password: resetPasswordForm.Password,
//	})
//	if err != nil {
//		HandleGrpcErrorToHttp(c, err)
//		return
//	}
//	HandleHttpResponse(c, http.StatusOK, errorx.Success, nil, nil)
//	return
//}
//
//func UpdateUser(c *gin.Context) {
//	updateUserForm := form.UpdateUserForm{}
//	if err := c.ShouldBind(&updateUserForm); err != nil {
//		HandleValidatorError(c, err)
//		return
//	}
//	token := c.GetString("token")
//	_, err := global.UserServiceClient.UpdateUser(NewCtxWithToken(context.Background(), token), &userPb.UpdateUserRequest{
//		Username: updateUserForm.Username,
//		Avatar:   updateUserForm.Avatar,
//	})
//	if err != nil {
//		HandleGrpcErrorToHttp(c, err)
//		return
//	}
//	j := util.NewJWT(global.ServerConfig.JWT.SigningKey, global.ServerConfig.JWT.Expiration)
//	refreshedToken, _ := j.RefreshToken(token)
//	HandleHttpResponse(c, http.StatusOK, errorx.Success, refreshedToken, nil)
//	return
//}
