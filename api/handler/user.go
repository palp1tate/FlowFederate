package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/palp1tate/FlowFederate/api/dao"
	"github.com/palp1tate/FlowFederate/api/global"
	"github.com/palp1tate/FlowFederate/api/internal/consts"
	"github.com/palp1tate/FlowFederate/api/internal/errorx"
	"github.com/palp1tate/FlowFederate/api/internal/model"
	"github.com/palp1tate/FlowFederate/api/internal/utils"
	"github.com/palp1tate/FlowFederate/api/types"

	"github.com/gin-gonic/gin"
	"github.com/palp1tate/go-crypto-guard/rsa"
	"github.com/redis/go-redis/v9"
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
	encodedPassword, _ := pwd.GenRSA(registerForm.Password, consts.RsaPublicKeyPath)
	user := model.User{
		Username: registerForm.Username,
		Mobile:   registerForm.Mobile,
		Password: encodedPassword,
		Role:     consts.User,
	}
	if err = dao.NewUserDao(context.Background()).CreateUser(&user); err != nil {
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

func UserLogin(c *gin.Context) {
	loginForm := types.LoginForm{}
	if err := c.ShouldBind(&loginForm); err != nil {
		HandleValidatorError(c, err)
		return
	}
	user := new(model.User)
	switch loginForm.Type {
	case consts.LoginByPassword:
		if ok := utils.Store.Verify(loginForm.CaptchaId, loginForm.Captcha, true); !ok {
			HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrCaptchaIncorrect)
			return
		}
		var err error
		user, err = dao.NewUserDao(context.Background()).FindUserByMobile(loginForm.Mobile)
		if err != nil {
			HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrMobileNotFound)
			return
		}
		if ok, _ := pwd.VerifyRSA(loginForm.Password, user.Password, consts.RsaPrivateKeyPath); !ok {
			HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrPasswordIncorrect)
			return
		}
	case consts.LoginByCaptcha:
		key := fmt.Sprintf("%d-%s", consts.Login, loginForm.Mobile)
		smsCode, err := global.RedisClient.Get(context.Background(), key).Result()
		if errors.Is(err, redis.Nil) {
			HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrCaptchaExpired)
			return
		} else {
			if smsCode != loginForm.Captcha {
				HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrCaptchaIncorrect)
				return
			}
			global.RedisClient.Del(context.Background(), key)
		}
		user, err = dao.NewUserDao(context.Background()).FindUserByMobile(loginForm.Mobile)
		if err != nil {
			HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrMobileNotFound)
			return
		}
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

func GetUser(c *gin.Context) {
	user, err := dao.NewUserDao(context.Background()).FindUserById(c.GetInt("id"))
	if err != nil {
		HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrUserNotFound)
		return
	}
	token := c.GetString("token")
	j := utils.NewJWT(global.ServerConfig.JWT.SigningKey, global.ServerConfig.JWT.Expiration)
	refreshedToken, _ := j.RefreshToken(token)
	HandleHttpResponse(c, http.StatusOK, errorx.Success, refreshedToken, UserModeToResponse(user))
}

func ResetPassword(c *gin.Context) {
	resetPasswordForm := types.ResetPasswordForm{}
	if err := c.ShouldBind(&resetPasswordForm); err != nil {
		HandleValidatorError(c, err)
		return
	}
	key := fmt.Sprintf("%d-%s", consts.ResetPassword, resetPasswordForm.Mobile)
	smsCode, err := global.RedisClient.Get(context.Background(), key).Result()
	if errors.Is(err, redis.Nil) {
		HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrCaptchaExpired)
		return
	} else {
		if smsCode != resetPasswordForm.Captcha {
			HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrCaptchaIncorrect)
			return
		}
		global.RedisClient.Del(context.Background(), key)
	}
	user, err := dao.NewUserDao(context.Background()).FindUserByMobile(resetPasswordForm.Mobile)
	if err != nil {
		HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrMobileNotFound)
		return
	}
	encodedPassword, _ := pwd.GenRSA(resetPasswordForm.Password, consts.RsaPublicKeyPath)
	if err := dao.NewUserDao(context.Background()).UpdatePassword(user, encodedPassword); err != nil {
		HandleHttpResponse(c, http.StatusInternalServerError, errorx.ErrPasswordReset)
		return
	}
	HandleHttpResponse(c, http.StatusOK, errorx.Success)
}

func UpdateUser(c *gin.Context) {
	updateUserForm := types.UpdateUserForm{}
	if err := c.ShouldBind(&updateUserForm); err != nil {
		HandleValidatorError(c, err)
		return
	}
	user, err := dao.NewUserDao(context.Background()).FindUserById(c.GetInt("id"))
	if err != nil {
		HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrUserNotFound)
		return
	}
	user.Username = updateUserForm.Username
	user.Avatar = updateUserForm.Avatar
	if err := dao.NewUserDao(context.Background()).UpdateUser(user); err != nil {
		HandleHttpResponse(c, http.StatusInternalServerError, errorx.ErrUpdateUserFailed)
		return
	}
	token := c.GetString("token")
	j := utils.NewJWT(global.ServerConfig.JWT.SigningKey, global.ServerConfig.JWT.Expiration)
	refreshedToken, _ := j.RefreshToken(token)
	HandleHttpResponse(c, http.StatusOK, errorx.Success, refreshedToken, nil)
}

func UserModeToResponse(user *model.User) types.User {
	return types.User{
		ID:        user.ID,
		Username:  user.Username,
		Mobile:    user.Mobile,
		Avatar:    user.Avatar,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04"),
	}
}
