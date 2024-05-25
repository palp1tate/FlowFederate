package utils

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func ResourceExtractor(c *gin.Context) string {
	// 获取路径并按"/"分割
	pathSegments := strings.Split(c.FullPath(), "/")
	// 检查路径是否至少有三部分（例如"/api/user/register"）
	if len(pathSegments) >= 3 {
		// 使用路径的第二部分作为资源名称
		return pathSegments[2]
	}
	// 如果路径不符合预期，使用默认的资源名称
	return c.FullPath()
}
