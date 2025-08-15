package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
	"gorm.io/gorm"
)

// KnowledgeBase 知识库模型
type KnowledgeBase struct {
	ID        uint              `gorm:"primarykey" json:"id"`
	Title     string            `gorm:"size:200;not null" json:"title"`
	Content   string            `gorm:"type:text;not null" json:"content"`
	Category  string            `gorm:"size:50" json:"category"`
	Tags      TagList           `gorm:"type:json" json:"tags"`
	FilePath  string            `gorm:"size:500" json:"file_path"`
	VectorID  string            `gorm:"size:100" json:"vector_id"`
	Status    KnowledgeStatus   `gorm:"type:tinyint;default:1" json:"status"`
	CreatedBy uint              `gorm:"index" json:"created_by"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	DeletedAt gorm.DeletedAt    `gorm:"index" json:"-"`
	
	// 关联
	Creator User `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

// TagList 标签列表
type TagList []string

// KnowledgeStatus 知识状态
type KnowledgeStatus int

const (
	KnowledgeStatusInactive KnowledgeStatus = 0 // 未激活
	KnowledgeStatusActive   KnowledgeStatus = 1 // 激活
)

// KnowledgeCategory 知识分类常量
const (
	CategoryResumeTips     = "resume_tips"      // 简历技巧
	CategoryIndustryGuide  = "industry_guide"   // 行业指南
	CategoryPositionReq    = "position_req"     // 岗位要求
	CategoryBestPractice   = "best_practice"    // 最佳实践
	CategoryCommonMistake  = "common_mistake"   // 常见错误
	CategoryTemplateGuide  = "template_guide"   // 模板指南
)

// Value 实现 driver.Valuer 接口
func (t TagList) Value() (driver.Value, error) {
	return json.Marshal(t)
}

// Scan 实现 sql.Scanner 接口
func (t *TagList) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	
	return json.Unmarshal(bytes, t)
}

// TableName 表名
func (KnowledgeBase) TableName() string {
	return "knowledge_base"
}

// KnowledgeChunk 知识块（用于向量化存储）
type KnowledgeChunk struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	KnowledgeID  uint      `gorm:"not null;index" json:"knowledge_id"`
	ChunkText    string    `gorm:"type:text;not null" json:"chunk_text"`
	ChunkIndex   int       `gorm:"not null" json:"chunk_index"`
	VectorID     string    `gorm:"size:100" json:"vector_id"`
	TokenCount   int       `json:"token_count"`
	CreatedAt    time.Time `json:"created_at"`
	
	// 关联
	Knowledge KnowledgeBase `gorm:"foreignKey:KnowledgeID" json:"knowledge,omitempty"`
}

// TableName 表名
func (KnowledgeChunk) TableName() string {
	return "knowledge_chunks"
}
