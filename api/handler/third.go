package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/palp1tate/FlowFederate/api/dao"
	"github.com/palp1tate/FlowFederate/api/global"
	"github.com/palp1tate/FlowFederate/api/internal/consts"
	"github.com/palp1tate/FlowFederate/api/internal/errorx"
	"github.com/palp1tate/FlowFederate/api/internal/utils"
	"github.com/palp1tate/FlowFederate/api/types"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetPicCaptcha(ctx *gin.Context) {
	id, b64s, err := utils.GeneratePicCaptcha()
	if err != nil {
		zap.S().Error(err)
		HandleHttpResponse(ctx, http.StatusInternalServerError, errorx.ErrPicCaptchaGenerate)
		return

	}
	HandleHttpResponse(ctx, http.StatusOK, errorx.Success, nil, gin.H{
		"captchaId": id,
		"picPath":   b64s,
	})
}

func SendSms(c *gin.Context) {
	sendSmsForm := types.SendSmsForm{}
	if err := c.ShouldBind(&sendSmsForm); err != nil {
		HandleValidatorError(c, err)
		return
	}
	e, b := sentinel.Entry(consts.SMSResource, sentinel.WithTrafficType(base.Inbound), sentinel.WithArgs(sendSmsForm.Mobile))
	if b != nil {
		HandleHttpResponse(c, http.StatusTooManyRequests, errorx.TooManyRequests)
		return
	}
	defer e.Exit()
	if sendSmsForm.Role == consts.User {
		switch sendSmsForm.Type {
		case consts.Register:
			if _, err := dao.NewUserDao(context.Background()).FindUserByMobile(sendSmsForm.Mobile); err == nil {
				HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrMobileExists)
				return
			}
			if captcha, err := utils.SendSms(sendSmsForm.Mobile); err != nil {
				zap.S().Error(err)
				HandleHttpResponse(c, http.StatusInternalServerError, errorx.ErrSendSmsFailed)
				return
			} else {
				global.RedisClient.Set(context.Background(), fmt.Sprintf("%d-%s", consts.Register, sendSmsForm.Mobile), captcha,
					time.Duration(global.ServerConfig.Redis.Expiration)*time.Minute)
			}

		case consts.Login:
			if _, err := dao.NewUserDao(context.Background()).FindUserByMobile(sendSmsForm.Mobile); err != nil {
				HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrMobileNotFound)
				return
			}
			if captcha, err := utils.SendSms(sendSmsForm.Mobile); err != nil {
				zap.S().Error(err)
				HandleHttpResponse(c, http.StatusInternalServerError, errorx.ErrSendSmsFailed)
				return
			} else {
				global.RedisClient.Set(context.Background(), fmt.Sprintf("%d-%s", consts.Login, sendSmsForm.Mobile), captcha,
					time.Duration(global.ServerConfig.Redis.Expiration)*time.Minute)
			}

		case consts.ResetPassword:
			if _, err := dao.NewUserDao(context.Background()).FindUserByMobile(sendSmsForm.Mobile); err != nil {
				HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrMobileNotFound)
				return
			}
			if captcha, err := utils.SendSms(sendSmsForm.Mobile); err != nil {
				zap.S().Error(err)
				HandleHttpResponse(c, http.StatusInternalServerError, errorx.ErrSendSmsFailed)
				return
			} else {
				global.RedisClient.Set(context.Background(), fmt.Sprintf("%d-%s", consts.ResetPassword, sendSmsForm.Mobile), captcha,
					time.Duration(global.ServerConfig.Redis.Expiration)*time.Minute)
			}
		}
	} else {
		if sendSmsForm.Type == consts.Login {
			if _, err := dao.NewAdminDao(context.Background()).FindAdminByMobile(sendSmsForm.Mobile); err != nil {
				HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrMobileNotFound)
				return
			}
			if captcha, err := utils.SendSms(sendSmsForm.Mobile); err != nil {
				zap.S().Error(err)
				HandleHttpResponse(c, http.StatusInternalServerError, errorx.ErrSendSmsFailed)
				return
			} else {
				global.RedisClient.Set(context.Background(), fmt.Sprintf("%d-%s", consts.Login, sendSmsForm.Mobile), captcha,
					time.Duration(global.ServerConfig.Redis.Expiration)*time.Minute)
			}
		} else if sendSmsForm.Type == consts.ResetPassword {
			if _, err := dao.NewAdminDao(context.Background()).FindAdminByMobile(sendSmsForm.Mobile); err != nil {
				HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrMobileNotFound)
				return
			}
			if captcha, err := utils.SendSms(sendSmsForm.Mobile); err != nil {
				zap.S().Error(err)
				HandleHttpResponse(c, http.StatusInternalServerError, errorx.ErrSendSmsFailed)
				return
			} else {
				global.RedisClient.Set(context.Background(), fmt.Sprintf("%d-%s", consts.ResetPassword, sendSmsForm.Mobile), captcha,
					time.Duration(global.ServerConfig.Redis.Expiration)*time.Minute)
			}
		} else {
			HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrParamsInvalid)
			return
		}
	}

	HandleHttpResponse(c, http.StatusOK, errorx.Success)
}

func UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		zap.S().Info(err)
		HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrParamsInvalid)
		return
	}
	url, err := utils.UploadFile(file, header)
	if err != nil {
		zap.S().Error(err)
		HandleHttpResponse(c, http.StatusInternalServerError, errorx.ErrUploadFileFailed)
	}
	token := c.GetString("token")
	j := utils.NewJWT(global.ServerConfig.JWT.SigningKey, global.ServerConfig.JWT.Expiration)
	refreshedToken, _ := j.RefreshToken(token)
	HandleHttpResponse(c, http.StatusOK, errorx.Success, refreshedToken, url)
}

func DeleteFile(c *gin.Context) {
	urlForm := types.UrlForm{}
	if err := c.ShouldBind(&urlForm); err != nil {
		HandleValidatorError(c, err)
		return
	}
	if err := utils.DeleteFile(urlForm.Url); err != nil {
		zap.S().Error(err)
		HandleHttpResponse(c, http.StatusInternalServerError, errorx.ErrDeleteFileFailed)
		return
	}
	token := c.GetString("token")
	j := utils.NewJWT(global.ServerConfig.JWT.SigningKey, global.ServerConfig.JWT.Expiration)
	refreshedToken, _ := j.RefreshToken(token)
	HandleHttpResponse(c, http.StatusOK, errorx.Success, refreshedToken)
}
