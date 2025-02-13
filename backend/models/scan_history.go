package models

import "time"

type ScanHistory struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	TaskType  string    `gorm:"type:varchar(50);not null" json:"taskType"` // 扫描类型：端口/主机/漏洞等
	Target    string    `gorm:"type:text;not null" json:"target"`          // 扫描目标
	Status    string    `gorm:"type:varchar(20);not null" json:"status"`   // 扫描状态：完成/失败
	Result    string    `gorm:"type:longtext" json:"result"`               // 扫描结果（JSON字符串）
	StartTime time.Time `gorm:"not null" json:"startTime"`                 // 开始时间
	EndTime   time.Time `json:"endTime"`                                   // 结束时间
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
