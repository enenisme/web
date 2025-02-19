package api

import (
	"fmt"
	"net/http"

	"github.com/enenisme/hostalive"
	"github.com/gin-gonic/gin"
)

type HostAliveHandler struct {
	// 可以添加依赖服务
}

func NewHostAliveHandler() *HostAliveHandler {
	return &HostAliveHandler{}
}

func (h *HostAliveHandler) HandleHostAliveCheck(c *gin.Context) {
	var request struct {
		Target string `json:"target" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求参数: " + err.Error(),
		})
		return
	}

	results, err := h.processHostAliveCheck(request.Target)
	if err != nil {
		fmt.Println("ERROR: ", err)
		c.JSON(http.StatusOK, gin.H{
			"message": "检测失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "检测完成",
		"data":    results,
	})
}

func (h *HostAliveHandler) processHostAliveCheck(targets string) ([]hostalive.ScanResult, error) {
	ha := hostalive.NewHostAlive(targets, false, 3, 100)
	return ha.HostAlive()
}
