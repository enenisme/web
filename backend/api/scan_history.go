package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"backend/global"
	"backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ScanHistoryHandler struct {
	db *gorm.DB
}

func NewScanHistoryHandler() *ScanHistoryHandler {
	return &ScanHistoryHandler{
		db: global.DB,
	}
}

// 获取扫描历史列表
func (h *ScanHistoryHandler) GetScanHistoryList(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	taskType := c.Query("taskType")

	var histories []models.ScanHistory
	var total int64
	query := h.db.Model(&models.ScanHistory{})

	// 如果指定了任务类型，添加筛选条件
	if taskType != "" {
		query = query.Where("task_type = ?", taskType)
	}

	// 计算总数
	query.Count(&total)

	// 获取分页数据
	err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&histories).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取扫描历史失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"data":  histories,
	})
}

// 获取扫描历史详情
func (h *ScanHistoryHandler) GetScanHistoryDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的ID",
		})
		return
	}

	var history models.ScanHistory
	if err := h.db.First(&history, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "未找到该记录",
		})
		return
	}

	c.JSON(http.StatusOK, history)
}

// 删除扫描历史
func (h *ScanHistoryHandler) DeleteScanHistory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的ID",
		})
		return
	}

	if err := h.db.Delete(&models.ScanHistory{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功",
	})
}

// 记录扫描历史的辅助方法
func (h *ScanHistoryHandler) RecordScanHistory(taskType, target string, result interface{}) error {
	history := models.ScanHistory{
		TaskType:  taskType,
		Target:    target,
		Status:    "completed",
		Result:    fmt.Sprintf("%v", result),
		StartTime: time.Now(),
		EndTime:   time.Now(),
	}

	return h.db.Create(&history).Error
}
