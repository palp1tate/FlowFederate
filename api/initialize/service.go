package initialize

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/palp1tate/FlowFederate/api/global"
	"github.com/palp1tate/FlowFederate/internal/consts"
	"github.com/palp1tate/FlowFederate/internal/util"
	"github.com/palp1tate/FlowFederate/service/edge/pb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func InitServiceConn() {
	credentials, err := util.GetClientCredentials()
	if err != nil {
		zap.S().Error(err)
	}
	consul := global.ServerConfig.Consul
	edgeConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s",
			consul.Host, consul.Port, global.ServerConfig.Service.Edge),
		grpc.WithTransportCredentials(credentials),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(consts.MaxMsgSize), grpc.MaxCallSendMsgSize(consts.MaxMsgSize)))
	if err != nil {
		zap.S().Error(err)
	}
	global.EdgeServiceClient = pb.NewEdgeServiceClient(edgeConn)
}
