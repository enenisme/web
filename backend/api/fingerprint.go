package api

import (
	"bufio"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/enenisme/definger"
	"github.com/gin-gonic/gin"
)

type FingerprintHandler struct {
	// 可以添加所需的服务依赖
}

func NewFingerprintHandler() *FingerprintHandler {
	return &FingerprintHandler{}
}

func (h *FingerprintHandler) HandleFingerprintScan(c *gin.Context) {
	inputType := c.PostForm("inputType")

	var targets []string

	switch inputType {
	case "url":
		// 处理单个URL输入
		url := c.PostForm("target")
		if url != "" {
			targets = append(targets, url)
		}
	case "file":
		// 处理多文件上传
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "文件上传失败: " + err.Error(),
			})
			return
		}

		files := form.File["files"]
		for _, file := range files {
			// 检查文件扩展名
			if filepath.Ext(file.Filename) != ".txt" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "只支持txt文件格式",
				})
				return
			}

			// 打开文件
			f, err := file.Open()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "无法读取文件: " + err.Error(),
				})
				return
			}
			defer f.Close()

			// 读取文件内容
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				url := strings.TrimSpace(scanner.Text())
				if url != "" {
					targets = append(targets, url)
				}
			}

			if err := scanner.Err(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "读取文件内容失败: " + err.Error(),
				})
				return
			}
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的输入类型",
		})
		return
	}

	// 验证是否有目标URL
	if len(targets) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "未提供有效的目标URL",
		})
		return
	}

	// TODO: 处理指纹识别逻辑
	results := h.processFingerprintScan(targets)

	c.JSON(http.StatusOK, gin.H{
		"message": "扫描完成",
		"data":    results,
	})
}

func (h *FingerprintHandler) processFingerprintScan(targets []string) []string {
	df := definger.NewDefinger(targets[0])
	result, err := df.Definger("C:\\Users\\张裕波\\Desktop\\project\\WEB-dev\\backend\\rule.json")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return result
}
