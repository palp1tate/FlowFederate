package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/palp1tate/FlowFederate/api/dao"
	"github.com/palp1tate/FlowFederate/api/internal/errorx"
	"github.com/palp1tate/FlowFederate/api/internal/model"
	"github.com/palp1tate/FlowFederate/api/internal/utils"
	"github.com/palp1tate/FlowFederate/api/types"

	"github.com/gin-gonic/gin"
)

//func Train(c *gin.Context) {
//	var trainForm form.TrainForm
//	if err := c.ShouldBind(&trainForm); err != nil {
//		HandleValidatorError(c, err)
//		return
//	}
//	if _, err := global.EdgeServiceClient.TrainTask(context.Background(), &pb.TrainRequest{
//		UserName: trainForm.UserName,
//		Conf:     trainForm.Conf,
//	}); err != nil {
//		HandleGrpcErrorToHttp(c, err)
//		return
//	}
//	HandleHttpResponse(c, http.StatusOK, errorx.Success, nil)
//}

func GetTaskList(c *gin.Context) {
	u := c.Query("u")
	if u == "" {
		HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrParamsInvalid, nil)
		return
	}
	page, pageSize := utils.ParsePageAndPageSize(c.Query("page"), c.Query("pageSize"))
	tasks, pages, totalCount, err := dao.NewTrainDao(context.Background()).FindTaskList(u, page, pageSize)
	if err != nil {
		HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrGetTaskList, nil)
		return
	}
	taskList := make([]types.TaskWithFormattedTime, len(tasks))
	for i, task := range tasks {
		taskList[i] = FormatTaskTime(task)
	}
	HandleHttpResponse(c, http.StatusOK, errorx.Success, gin.H{
		"tasks":      taskList,
		"pages":      pages,
		"totalCount": totalCount,
	})
}

func GetTask(c *gin.Context) {
	tid := c.Query("tid")
	if tid == "" {
		HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrParamsInvalid, nil)
		return
	}
	id, _ := strconv.Atoi(tid)
	task, err := dao.NewTrainDao(context.Background()).FindTask(id)
	if err != nil {
		HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrGetTask, nil)
		return
	}
	HandleHttpResponse(c, http.StatusOK, errorx.Success, FormatTaskTime(task))
}

func GetServerProgress(c *gin.Context) {
	tid := c.Query("tid")
	if tid == "" {
		HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrParamsInvalid, nil)
		return
	}
	id, _ := strconv.Atoi(tid)
	page, pageSize := utils.ParsePageAndPageSize(c.Query("page"), c.Query("pageSize"))
	servers, pages, totalCount, err := dao.NewTrainDao(context.Background()).FindServerList(id, page, pageSize)
	if err != nil {
		HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrGetServerList, nil)
		return
	}
	serverList := make([]types.ServerWithFormattedTime, len(servers))
	for i, server := range servers {
		serverList[i] = FormatServerTime(server)
	}
	HandleHttpResponse(c, http.StatusOK, errorx.Success, gin.H{
		"servers":    serverList,
		"pages":      pages,
		"totalCount": totalCount,
	})
}

func GetClientProgress(c *gin.Context) {
	tid := c.Query("tid")
	if tid == "" {
		HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrParamsInvalid, nil)
		return
	}
	id, _ := strconv.Atoi(tid)
	page, pageSize := utils.ParsePageAndPageSize(c.Query("page"), c.Query("pageSize"))
	clients, pages, totalCount, err := dao.NewTrainDao(context.Background()).FindClientList(id, page, pageSize)
	if err != nil {
		HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrGetClientList, nil)
		return
	}
	clientList := make([]types.ClientWithFormattedTime, len(clients))
	for i, client := range clients {
		clientList[i] = FormatClientTime(client)
	}
	HandleHttpResponse(c, http.StatusOK, errorx.Success, gin.H{
		"clients":    clientList,
		"pages":      pages,
		"totalCount": totalCount,
	})
}

func FormatTaskTime(t *model.Task) types.TaskWithFormattedTime {
	return types.TaskWithFormattedTime{
		ID:        t.ID,
		UserName:  t.UserName,
		CreatedAt: t.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: t.UpdatedAt.Format("2006-01-02 15:04:05"),
		Model:     t.Model,
		Dataset:   t.Dataset,
		Type:      t.Type,
		Status:    t.Status,
		Progress:  t.Progress,
		Accuracy:  t.Accuracy,
		Loss:      t.Loss,
	}
}

func FormatServerTime(s *model.Server) types.ServerWithFormattedTime {
	return types.ServerWithFormattedTime{
		ServerID:  s.ServerID,
		TaskID:    s.TaskID,
		CreatedAt: s.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: s.UpdatedAt.Format("2006-01-02 15:04:05"),
		Model:     s.Model,
		Dataset:   s.Dataset,
		Type:      s.Type,
		Status:    s.Status,
		Progress:  s.Progress,
		Accuracy:  s.Accuracy,
		Loss:      s.Loss,
		Cpu:       s.Cpu,
		Memory:    s.Memory,
		Disk:      s.Disk,
	}
}

func FormatClientTime(c *model.Client) types.ClientWithFormattedTime {
	return types.ClientWithFormattedTime{
		ClientID:  c.ClientID,
		TaskID:    c.TaskID,
		CreatedAt: c.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: c.UpdatedAt.Format("2006-01-02 15:04:05"),
		Model:     c.Model,
		Dataset:   c.Dataset,
		Type:      c.Type,
		Status:    c.Status,
		Progress:  c.Progress,
		Accuracy:  c.Accuracy,
		Loss:      c.Loss,
		Cpu:       c.Cpu,
		Memory:    c.Memory,
		Disk:      c.Disk,
	}
}
