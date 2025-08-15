package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
	"gorm.io/gorm"
)

// Resume 简历模型
type Resume struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	UserID        uint           `gorm:"not null;index" json:"user_id"`
	Title         string         `gorm:"size:200;not null" json:"title"`
	FilePath      string         `gorm:"size:500" json:"file_path"`
	FileType      FileType       `gorm:"type:varchar(20);not null" json:"file_type"`
	FileSize      int64          `json:"file_size"`
	ParsedContent ResumeContent  `gorm:"type:json" json:"parsed_content"`
	Status        ResumeStatus   `gorm:"type:tinyint;default:0" json:"status"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	
	// 关联
	User            User              `gorm:"foreignKey:UserID" json:"user,omitempty"`
	AnalysisResults []AnalysisResult  `gorm:"foreignKey:ResumeID" json:"analysis_results,omitempty"`
}

// FileType 文件类型
type FileType string

const (
	FileTypePDF      FileType = "pdf"
	FileTypeMarkdown FileType = "markdown"
)

// ResumeStatus 简历状态
type ResumeStatus int

const (
	ResumeStatusUploading ResumeStatus = 0 // 上传中
	ResumeStatusParsing   ResumeStatus = 1 // 解析中
	ResumeStatusParsed    ResumeStatus = 2 // 已解析
	ResumeStatusFailed    ResumeStatus = 3 // 解析失败
)

// ResumeContent 简历内容结构
type ResumeContent struct {
	PersonalInfo PersonalInfo `json:"personal_info"`
	Education    []Education  `json:"education"`
	Experience   []Experience `json:"experience"`
	Projects     []Project    `json:"projects"`
	Skills       Skills       `json:"skills"`
	Others       []Other      `json:"others"`
}

// PersonalInfo 个人信息
type PersonalInfo struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Location string `json:"location"`
	Website  string `json:"website"`
	LinkedIn string `json:"linkedin"`
	GitHub   string `json:"github"`
}

// Education 教育背景
type Education struct {
	School    string    `json:"school"`
	Degree    string    `json:"degree"`
	Major     string    `json:"major"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	GPA       string    `json:"gpa"`
	Honors    []string  `json:"honors"`
}

// Experience 工作经历
type Experience struct {
	Company      string    `json:"company"`
	Position     string    `json:"position"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	Description  string    `json:"description"`
	Achievements []string  `json:"achievements"`
}

// Project 项目经历
type Project struct {
	Name         string    `json:"name"`
	Role         string    `json:"role"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	Description  string    `json:"description"`
	Technologies []string  `json:"technologies"`
	Achievements []string  `json:"achievements"`
	URL          string    `json:"url"`
}

// Skills 技能
type Skills struct {
	Technical []string `json:"technical"`
	Languages []string `json:"languages"`
	Others    []string `json:"others"`
}

// Other 其他信息
type Other struct {
	Type        string `json:"type"`        // 类型：证书、获奖、志愿活动等
	Title       string `json:"title"`       // 标题
	Description string `json:"description"` // 描述
	Date        string `json:"date"`        // 日期
}

// Value 实现 driver.Valuer 接口
func (r ResumeContent) Value() (driver.Value, error) {
	return json.Marshal(r)
}

// Scan 实现 sql.Scanner 接口
func (r *ResumeContent) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	
	return json.Unmarshal(bytes, r)
}

// TableName 表名
func (Resume) TableName() string {
	return "resumes"
}
