package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/palp1tate/FlowFederate/api/global"
	"github.com/palp1tate/FlowFederate/internal/errorx"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/status"
)

type Response struct {
	Code int         `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func HandleHttpResponse(c *gin.Context, code int, xcode int, data interface{}) {
	response := Response{
		Code: xcode,
		Msg:  errorx.GetMsg(xcode),
		Data: data,
	}
	c.JSON(code, response)
}

func HandleGrpcErrorToHttp(c *gin.Context, err error) {
	if e, ok := status.FromError(err); ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.Code(),
			"msg":  e.Message(),
		})
		return
	}
}

func removeTopStruct(fields map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fields {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

func HandleValidatorError(c *gin.Context, err error) {
	var errs validator.ValidationErrors
	if ok := errors.As(err, &errs); !ok {
		HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrParamsInvalid, nil)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": errorx.ErrParamsInvalid,
			"msg":  removeTopStruct(errs.Translate(global.Translator)),
		})
	}
}
