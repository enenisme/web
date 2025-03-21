package api

import (
	"backend/global"
	"backend/models"
	"fmt"
	"net/http"
	"time"

	"github.com/enenisme/host_alive/scanner"
	"github.com/gin-gonic/gin"
)

type HostAliveHandler struct {
	scanner *scanner.HostScanner
}

func NewHostAliveHandler() *HostAliveHandler {
	return &HostAliveHandler{
		scanner: scanner.NewHostScanner(3*time.Second, 100),
	}
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

		global.DB.Create(&models.ScanHistory{
			TaskType:  "host_alive",
			Target:    request.Target,
			Status:    "failed",
			Result:    fmt.Sprintf("%v", err),
			StartTime: time.Now(),
			EndTime:   time.Now(),
		})
		return
	}

	global.DB.Create(&models.ScanHistory{
		TaskType:  "host_alive",
		Target:    request.Target,
		Status:    "completed",
		Result:    fmt.Sprintf("%v", results),
		StartTime: time.Now(),
		EndTime:   time.Now(),
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "检测完成",
		"data":    results,
	})
}

func (h *HostAliveHandler) processHostAliveCheck(targets string) ([]scanner.ScanResult, error) {
	return h.scanner.Scan(targets)
}
