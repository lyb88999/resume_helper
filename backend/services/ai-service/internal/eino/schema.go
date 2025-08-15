package eino

import (
	"time"
)

// ResumeData 简历数据结构
type ResumeData struct {
	ID           string            `json:"id"`
	Version      string            `json:"version"`
	PersonalInfo PersonalInfo      `json:"personal_info"`
	Education    []Education       `json:"education"`
	Experience   []Experience      `json:"experience"`
	Projects     []Project         `json:"projects"`
	Skills       Skills            `json:"skills"`
	Others       map[string]string `json:"others"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

// PersonalInfo 个人信息
type PersonalInfo struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Location string `json:"location"`
	LinkedIn string `json:"linkedin"`
	GitHub   string `json:"github"`
	Website  string `json:"website"`
}

// Education 教育背景
type Education struct {
	School    string    `json:"school"`
	Degree    string    `json:"degree"`
	Major     string    `json:"major"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	GPA       string    `json:"gpa"`
	Courses   []string  `json:"courses"`
	Honors    []string  `json:"honors"`
}

// Experience 工作经历
type Experience struct {
	Company      string    `json:"company"`
	Position     string    `json:"position"`
	Location     string    `json:"location"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	Description  []string  `json:"description"`
	Achievements []string  `json:"achievements"`
	Technologies []string  `json:"technologies"`
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

// Skills 技能特长
type Skills struct {
	Technical  []string `json:"technical"`
	Languages  []string `json:"languages"`
	Frameworks []string `json:"frameworks"`
	Tools      []string `json:"tools"`
	Soft       []string `json:"soft"`
}

// AnalysisResult已在factory.go中定义，避免重复

// SectionResult 章节分析结果
type SectionResult struct {
	Content       string  `json:"content"`
	ExtractedInfo string  `json:"extracted_info"`
	QualityScore  float64 `json:"quality_score"`
	Issues        []Issue `json:"issues"`
}

// Issue 问题
type Issue struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
	Suggestion  string `json:"suggestion"`
}

// Suggestion和ScoreBreakdown已在factory.go中定义

// Improvement已合并到Suggestion类型中

// Document 文档结构
type Document struct {
	ID       string            `json:"id"`
	Title    string            `json:"title"`
	Content  string            `json:"content"`
	Metadata map[string]string `json:"metadata"`
	Score    float64           `json:"score"`
}

// Message已在factory.go中定义

// ChatContext 聊天上下文
type ChatContext struct {
	SessionID      string      `json:"session_id"`
	Messages       []Message   `json:"messages"`
	ResumeData     *ResumeData `json:"resume_data,omitempty"`
	TargetPosition string      `json:"target_position,omitempty"`
	Knowledge      []Document  `json:"knowledge,omitempty"`
}

// AgentInput Agent输入
type AgentInput struct {
	Message string      `json:"message"`
	Context ChatContext `json:"context"`
}

// AgentOutput Agent输出
type AgentOutput struct {
	Message   string      `json:"message"`
	Sources   []string    `json:"sources"`
	Reasoning string      `json:"reasoning"`
	Context   ChatContext `json:"context"`
}

// WorkflowInput 工作流输入
type WorkflowInput struct {
	FilePath       string                 `json:"file_path"`
	FileType       string                 `json:"file_type"`
	TargetPosition string                 `json:"target_position"`
	Options        map[string]interface{} `json:"options"`
}

// WorkflowOutput 工作流输出
type WorkflowOutput struct {
	ResumeData     *ResumeData     `json:"resume_data"`
	AnalysisResult *AnalysisResult `json:"analysis_result"`
	Status         string          `json:"status"`
	Message        string          `json:"message"`
}
