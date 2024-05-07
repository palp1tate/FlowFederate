package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/palp1tate/FlowFederate/internal/model"

	"github.com/palp1tate/FlowFederate/api/dao"
	"github.com/palp1tate/FlowFederate/api/form"
	"github.com/palp1tate/FlowFederate/api/global"
	"github.com/palp1tate/FlowFederate/internal/errorx"
	"github.com/palp1tate/FlowFederate/service/edge/pb"

	"github.com/gin-gonic/gin"
)

func Train(c *gin.Context) {
	var trainForm form.TrainForm
	if err := c.ShouldBind(&trainForm); err != nil {
		HandleValidatorError(c, err)
		return
	}
	if _, err := global.EdgeServiceClient.TrainTask(context.Background(), &pb.TrainRequest{
		UserName: trainForm.UserName,
		Conf:     trainForm.Conf,
	}); err != nil {
		HandleGrpcErrorToHttp(c, err)
		return
	}
	HandleHttpResponse(c, http.StatusOK, errorx.Success, nil)
}

func tasksEqual(a, b []*model.Task) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if *a[i] != *b[i] {
			return false
		}
	}
	return true
}

func GetTaskList(c *gin.Context) {
	u := c.Query("u")
	if u == "" {
		HandleHttpResponse(c, http.StatusBadRequest, errorx.ErrParamsInvalid, nil)
		return
	}
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		HandleHttpResponse(c, http.StatusInternalServerError, errorx.ErrInternalServer, nil)
		return
	}
	defer ws.Close()

	var taskList []*model.Task
	done := make(chan struct{})

	go func() {
		defer close(done)
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				tasks, pages, totalCount, err := dao.NewTrainDao(context.Background()).FindTaskList(u, 1, 10)
				if err != nil {
					return
				}

				if !tasksEqual(taskList, tasks) {
					// 更新用户的任务列表
					taskList = cloneTasks(tasks)

					// 返回新的任务列表给前端
					if err := ws.WriteJSON(gin.H{
						"code": errorx.Success,
						"msg":  errorx.GetMsg(errorx.Success),
						"data": gin.H{
							"tasks":      tasks,
							"pages":      pages,
							"totalCount": totalCount,
						},
					}); err != nil {
						return
					}
				}
			case <-done:
				return
			}
		}
	}()

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			done <- struct{}{}
			break
		}
	}
}

func cloneTasks(tasks []*model.Task) []*model.Task {
	cloned := make([]*model.Task, len(tasks))
	for i, task := range tasks {
		taskCopy := *task
		cloned[i] = &taskCopy
	}
	return cloned
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
	HandleHttpResponse(c, http.StatusOK, errorx.Success, task)
}
