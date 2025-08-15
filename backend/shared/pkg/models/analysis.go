package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// AnalysisResult 分析结果模型
type AnalysisResult struct {
	ID                  uint           `gorm:"primarykey" json:"id"`
	ResumeID            uint           `gorm:"not null;index" json:"resume_id"`
	OverallScore        float64        `gorm:"type:decimal(3,1)" json:"overall_score"`
	CompletenessScore   float64        `gorm:"type:decimal(3,1)" json:"completeness_score"`
	ClarityScore        float64        `gorm:"type:decimal(3,1)" json:"clarity_score"`
	KeywordScore        float64        `gorm:"type:decimal(3,1)" json:"keyword_score"`
	FormatScore         float64        `gorm:"type:decimal(3,1)" json:"format_score"`
	QuantificationScore float64        `gorm:"type:decimal(3,1)" json:"quantification_score"`
	Suggestions         SuggestionList `gorm:"type:json" json:"suggestions"`
	AnalysisVersion     string         `gorm:"size:20" json:"analysis_version"`
	TargetPosition      string         `gorm:"size:100" json:"target_position"`
	Industry            string         `gorm:"size:50" json:"industry"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Resume Resume `gorm:"foreignKey:ResumeID" json:"resume,omitempty"`
}

// SuggestionList 建议列表
type SuggestionList []Suggestion

// Suggestion 优化建议
type Suggestion struct {
	Section     string          `json:"section"`     // 章节：experience, education, skills, etc.
	Level       SuggestionLevel `json:"level"`       // 严重程度
	Type        SuggestionType  `json:"type"`        // 建议类型
	Title       string          `json:"title"`       // 建议标题
	Description string          `json:"description"` // 详细描述
	Examples    []string        `json:"examples"`    // 示例
	Location    SuggestionLoc   `json:"location"`    // 位置信息
	Priority    int             `json:"priority"`    // 优先级 1-10
}

// SuggestionLevel 建议级别
type SuggestionLevel string

const (
	SuggestionLevelCritical SuggestionLevel = "critical" // 严重问题
	SuggestionLevelWarning  SuggestionLevel = "warning"  // 建议优化
	SuggestionLevelInfo     SuggestionLevel = "info"     // 锦上添花
)

// SuggestionType 建议类型
type SuggestionType string

const (
	SuggestionTypeContent   SuggestionType = "content"   // 内容优化
	SuggestionTypeFormat    SuggestionType = "format"    // 格式优化
	SuggestionTypeStructure SuggestionType = "structure" // 结构优化
	SuggestionTypeKeyword   SuggestionType = "keyword"   // 关键词优化
	SuggestionTypeQuantify  SuggestionType = "quantify"  // 量化优化
)

// SuggestionLoc 建议位置信息
type SuggestionLoc struct {
	Section string `json:"section"` // 章节名
	Index   int    `json:"index"`   // 在章节中的索引
	Field   string `json:"field"`   // 具体字段
}

// Value 实现 driver.Valuer 接口
func (s SuggestionList) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Scan 实现 sql.Scanner 接口
func (s *SuggestionList) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, s)
}

// TableName 表名
func (AnalysisResult) TableName() string {
	return "analysis_results"
}
