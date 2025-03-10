package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/enenisme/subfinder"
	"github.com/gin-gonic/gin"
)

type SubdomainHandler struct {
	// 可以添加依赖服务
}

func NewSubdomainHandler() *SubdomainHandler {
	return &SubdomainHandler{}
}

func (h *SubdomainHandler) HandleSubdomainScan(c *gin.Context) {
	var request struct {
		Domain     string `json:"domain" binding:"required"`
		DomainFile string `json:"domain_file" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求参数: " + err.Error(),
		})
		return
	}

	results := h.processSubdomainScan(request.Domain, request.DomainFile)
	c.JSON(http.StatusOK, gin.H{
		"message": "扫描完成",
		"data":    results,
	})
}

func (h *SubdomainHandler) processSubdomainScan(domain string, domainFile string) map[string]interface{} {
	// TODO: 实现具体的子域名扫描逻辑
	sf := subfinder.NewSubfinder(domain, domainFile)
	if err := sf.Run(); err != nil {
		fmt.Println("subfinder run error: ", err)
	}
	jsonData, err := json.Marshal(sf.Subdomains)
	if err != nil {
		fmt.Println("json marshal error: ", err)
	}
	fmt.Println("subdomains: ", string(jsonData))

	return nil
}
