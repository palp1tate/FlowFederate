package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/gin-gonic/gin"
	"github.com/palp1tate/FlowFederate/api/dao"
	"github.com/palp1tate/FlowFederate/api/global"
	"github.com/palp1tate/FlowFederate/api/internal/consts"
	"github.com/palp1tate/FlowFederate/api/internal/errorx"
	"github.com/palp1tate/FlowFederate/api/internal/utils"
	"github.com/palp1tate/FlowFederate/api/types"
)

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
