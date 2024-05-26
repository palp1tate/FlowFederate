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

func AdminLogin(c *gin.Context) {
	loginForm := types.LoginForm{}
	if err := c.ShouldBind(&loginForm); err != nil {
		HandleValidatorError(c, err)
		return
	}
	admin := new(model.User)
	switch loginForm.Type {
	case consts.LoginByPassword:
		if ok := utils.Store.Verify(loginForm.CaptchaId, loginForm.Captcha, true); !ok {
			HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrCaptchaIncorrect)
			return
		}
		var err error
		admin, err = dao.NewAdminDao(context.Background()).FindAdminByMobile(loginForm.Mobile)
		if err != nil {
			HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrMobileNotFound)
			return

		}
		if ok, _ := pwd.VerifyRSA(loginForm.Password, admin.Password, consts.RsaPrivateKeyPath); !ok {
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
		admin, err = dao.NewAdminDao(context.Background()).FindAdminByMobile(loginForm.Mobile)
		if err != nil {
			HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrMobileNotFound)
			return
		}
	}
	claims := utils.CustomClaims{
		ID:   strconv.Itoa(admin.ID),
		Role: strconv.Itoa(consts.Admin),
	}
	j := utils.NewJWT(global.ServerConfig.JWT.SigningKey, global.ServerConfig.JWT.Expiration)
	token, err := j.CreateToken(claims)
	if err != nil {
		HandleHttpResponse(c, http.StatusInternalServerError, errorx.ErrGenTokenFailed)
		return
	}
	HandleHttpResponse(c, http.StatusOK, errorx.Success, token)
}

func GetAdmin(c *gin.Context) {
	admin, err := dao.NewAdminDao(context.Background()).FindAdminById(c.GetInt("id"))
	if err != nil {
		HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrAdminNotFound)
		return

	}
	token := c.GetString("token")
	j := utils.NewJWT(global.ServerConfig.JWT.SigningKey, global.ServerConfig.JWT.Expiration)
	refreshedToken, _ := j.RefreshToken(token)
	HandleHttpResponse(c, http.StatusOK, errorx.Success, refreshedToken, UserModeToResponse(admin))
}

func UpdateAdmin(c *gin.Context) {
	updateAdminForm := types.UpdateUserForm{}
	if err := c.ShouldBind(&updateAdminForm); err != nil {
		HandleValidatorError(c, err)
		return
	}
	admin, err := dao.NewAdminDao(context.Background()).FindAdminById(c.GetInt("id"))
	if err != nil {
		HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrAdminNotFound)
		return
	}
	admin.Username = updateAdminForm.Username
	admin.Avatar = updateAdminForm.Avatar
	if err = dao.NewUserDao(context.Background()).UpdateUser(admin); err != nil {
		HandleHttpResponse(c, http.StatusInternalServerError, errorx.ErrUpdateAdminFailed)
		return

	}
	token := c.GetString("token")
	j := utils.NewJWT(global.ServerConfig.JWT.SigningKey, global.ServerConfig.JWT.Expiration)
	refreshedToken, _ := j.RefreshToken(token)
	HandleHttpResponse(c, http.StatusOK, errorx.Success, refreshedToken)
}

func AdminResetPassword(c *gin.Context) {
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
	admin, err := dao.NewAdminDao(context.Background()).FindAdminByMobile(resetPasswordForm.Mobile)
	if err != nil {
		HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrMobileNotFound)
		return

	}
	encodedPassword, _ := pwd.GenRSA(resetPasswordForm.Password, consts.RsaPublicKeyPath)
	if err = dao.NewUserDao(context.Background()).UpdatePassword(admin, encodedPassword); err != nil {
		HandleHttpResponse(c, http.StatusInternalServerError, errorx.ErrPasswordReset)
		return

	}
	HandleHttpResponse(c, http.StatusOK, errorx.Success)
}

func AdminGetUser(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Query("uid"))
	if userId == 0 {
		HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrParamsInvalid)
		return
	}
	user, err := dao.NewUserDao(context.Background()).FindUserById(userId)
	if err != nil {
		HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrUserNotFound)
		return
	}
	token := c.GetString("token")
	j := utils.NewJWT(global.ServerConfig.JWT.SigningKey, global.ServerConfig.JWT.Expiration)
	refreshedToken, _ := j.RefreshToken(token)
	HandleHttpResponse(c, http.StatusOK, errorx.Success, refreshedToken, UserModeToResponse(user))
}

func GetUserList(c *gin.Context) {
	page, pageSize := utils.ParsePageAndPageSize(c.Query("page"), c.Query("pageSize"))
	users, pages, totalCount, err := dao.NewAdminDao(context.Background()).FindUserList(page, pageSize)
	if err != nil {
		HandleHttpResponse(c, http.StatusInternalServerError, errorx.ErrGetUserListFailed)
		return
	}
	userList := make([]types.User, len(users))
	for i, user := range users {
		userList[i] = UserModeToResponse(user)
	}
	token := c.GetString("token")
	j := utils.NewJWT(global.ServerConfig.JWT.SigningKey, global.ServerConfig.JWT.Expiration)
	refreshedToken, _ := j.RefreshToken(token)
	HandleHttpResponse(c, http.StatusOK, errorx.Success, refreshedToken, gin.H{
		"users":      userList,
		"pages":      pages,
		"totalCount": totalCount,
	})
}

func AdminUpdateUser(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Query("uid"))
	if userId == 0 {
		HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrParamsInvalid)
		return
	}
	updateUserForm := types.UpdateUserForm{}
	if err := c.ShouldBind(&updateUserForm); err != nil {
		HandleValidatorError(c, err)
		return
	}
	user, err := dao.NewUserDao(context.Background()).FindUserById(userId)
	if err != nil {
		HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrUserNotFound)
		return
	}
	user.Username = updateUserForm.Username
	user.Avatar = updateUserForm.Avatar
	if err = dao.NewUserDao(context.Background()).UpdateUser(user); err != nil {
		HandleHttpResponse(c, http.StatusInternalServerError, errorx.ErrUpdateUserFailed)
		return
	}
	token := c.GetString("token")
	j := utils.NewJWT(global.ServerConfig.JWT.SigningKey, global.ServerConfig.JWT.Expiration)
	refreshedToken, _ := j.RefreshToken(token)
	HandleHttpResponse(c, http.StatusOK, errorx.Success, refreshedToken)
}

func DeleteUser(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Query("uid"))
	if userId == 0 {
		HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrParamsInvalid)
		return
	}
	user, err := dao.NewUserDao(context.Background()).FindUserById(userId)
	if err != nil {
		HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrUserNotFound)
		return
	}
	if err = dao.NewAdminDao(context.Background()).DeleteUser(user); err != nil {
		HandleHttpResponse(c, http.StatusInternalServerError, errorx.ErrDeleteUserFailed)
		return
	}
	token := c.GetString("token")
	j := utils.NewJWT(global.ServerConfig.JWT.SigningKey, global.ServerConfig.JWT.Expiration)
	refreshedToken, _ := j.RefreshToken(token)
	HandleHttpResponse(c, http.StatusOK, errorx.Success, refreshedToken)
}
