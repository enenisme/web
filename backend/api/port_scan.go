package api

import (
	"backend/global"
	"backend/models"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/enenisme/port_scan"
	"github.com/gin-gonic/gin"
)

type PortScanHandler struct {
	scanner *port_scan.Scanner
}

func NewPortScanHandler() *PortScanHandler {
	return &PortScanHandler{
		scanner: port_scan.NewScanner(2 * time.Second), // 设置超时时间为2秒
	}
}

func (h *PortScanHandler) HandlePortScan(c *gin.Context) {
	var request struct {
		Target    string `json:"target" binding:"required"`
		PortRange string `json:"portRange"` // 例如: "80,443" 或 "1-1000"
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求参数: " + err.Error(),
		})
		return
	}

	// 验证目标地址
	if request.Target == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "目标地址不能为空",
		})
		return
	}

	// 解析端口范围
	ports, err := h.parsePortRange(request.PortRange)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的端口范围: " + err.Error(),
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

	// 执行端口扫描
	results := h.scanner.ScanPorts(request.Target, ports)

	// 格式化返回结果
	formattedResults := h.formatResults(request.Target, results)

	global.DB.Create(&models.ScanHistory{
		TaskType:  "port_scan",
		Target:    request.Target,
		Status:    "completed",
		Result:    fmt.Sprintf("%v", formattedResults),
		StartTime: time.Now(),
		EndTime:   time.Now(),
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "扫描完成",
		"data":    formattedResults,
	})
}

// parsePortRange 解析端口范围字符串
func (h *PortScanHandler) parsePortRange(portRange string) ([]int, error) {
	if portRange == "" {
		// 如果未指定端口范围，使用常用端口列表
		return h.scanner.LoadPortsFromFile("port_scan/common_ports.txt")
	}

	var ports []int
	// 处理逗号分隔的端口列表，如 "80,443,8080"
	if strings.Contains(portRange, ",") {
		for _, p := range strings.Split(portRange, ",") {
			port, err := strconv.Atoi(strings.TrimSpace(p))
			if err != nil {
				return nil, err
			}
			ports = append(ports, port)
		}
		return ports, nil
	}

	// 处理范围格式，如 "1-1000"
	if strings.Contains(portRange, "-") {
		parts := strings.Split(portRange, "-")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid port range format")
		}

		start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return nil, err
		}

		end, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return nil, err
		}

		for i := start; i <= end; i++ {
			ports = append(ports, i)
		}
		return ports, nil
	}

	// 单个端口
	port, err := strconv.Atoi(strings.TrimSpace(portRange))
	if err != nil {
		return nil, err
	}
	return []int{port}, nil
}

// formatResults 格式化扫描结果
func (h *PortScanHandler) formatResults(target string, results []port_scan.PortScan) map[string]interface{} {
	openPorts := make([]map[string]interface{}, 0)
	for _, result := range results {
		if result.State {
			portInfo := map[string]interface{}{
				"port":    result.Port,
				"state":   "open",
				"service": result.Service,
			}
			openPorts = append(openPorts, portInfo)
		}
	}

	return map[string]interface{}{
		"host":      target,
		"portCount": len(openPorts),
		"ports":     openPorts,
	}
}
