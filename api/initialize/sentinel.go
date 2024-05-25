package initialize

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/core/hotspot"
	"github.com/palp1tate/FlowFederate/api/internal/consts"
	"go.uber.org/zap"
)

func InitSentinel() {
	err := sentinel.InitDefault()
	if err != nil {
		zap.S().Fatalf(err.Error())
	}
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               consts.UserResource,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              3000,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               consts.ThirdResource,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              2000,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               consts.TrainResource,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              1500,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               consts.AdminResource,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              2000,
			StatIntervalInMs:       1000,
		},
	})

	if err != nil {
		zap.S().Fatalf(err.Error())
	}
	_, err = hotspot.LoadRules([]*hotspot.Rule{
		{
			Resource:        consts.SMSResource,
			MetricType:      hotspot.QPS,
			ControlBehavior: hotspot.Reject,
			ParamIndex:      0, // 根据手机号限流
			Threshold:       1,
			DurationInSec:   60,
		},
	})

	if err != nil {
		zap.S().Fatalf(err.Error())
	}
}
