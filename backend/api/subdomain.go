package api

import (
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
		DomainFile string `json:"domain_file" `
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

type SubdomainResult struct {
	Domain     string `json:"domain"`
	IP         string `json:"ip"`
	Source     string `json:"source"`
	RecordType string `json:"record_type"`
}

func (h *SubdomainHandler) processSubdomainScan(domain string, domainFile string) []SubdomainResult {
	// TODO: 实现具体的子域名扫描逻辑
	sf := subfinder.NewSubfinder(domain, domainFile, 120)
	if err := sf.Run(); err != nil {
		fmt.Println("subfinder run error: ", err)
	}
	subdomains := make([]SubdomainResult, 0)
	for _, subdomain := range sf.Subdomains {
		subdomains = append(subdomains, SubdomainResult{
			Domain:     subdomain.Domain,
			IP:         subdomain.IP,
			Source:     subdomain.Source,
			RecordType: subdomain.RecordType,
		})
	}

	return subdomains
}
