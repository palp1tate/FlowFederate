package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/palp1tate/FlowFederate/api/form"
	"github.com/palp1tate/FlowFederate/api/global"
	"github.com/palp1tate/FlowFederate/internal/errorx"
	"github.com/palp1tate/FlowFederate/service/edge/pb"
	"net/http"
)

func Train(c *gin.Context) {
	var trainForm form.TrainForm
	if err := c.ShouldBind(&trainForm); err != nil {
		HandleValidatorError(c, err)
		return
	}
	if _, err := global.EdgeServiceClient.TrainTask(context.Background(), &pb.TrainRequest{
		Conf: trainForm.Conf,
	}); err != nil {
		HandleGrpcErrorToHttp(c, err)
		return
	}
	HandleHttpResponse(c, http.StatusOK, errorx.Success, nil)
}
