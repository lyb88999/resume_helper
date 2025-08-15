package biz

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

// 错误定义
var (
	ErrTaskNotFound      = errors.New("task not found")
	ErrUnsupportedType   = errors.New("unsupported file type")
	ErrFileNotFound      = errors.New("file not found")
	ErrParseTimeout      = errors.New("parse timeout")
	ErrFileTooLarge      = errors.New("file too large")
	ErrDocumentProtected = errors.New("document is password protected")
	ErrEmptyContent      = errors.New("document content is empty")
)

// ParseOptions 解析选项
type ParseOptions struct {
	ExtractImages  bool     `json:"extract_images"`
	CleanText      bool     `json:"clean_text"`
	TargetLanguage string   `json:"target_language"`
	SkipSections   []string `json:"skip_sections"`
}

// ParseTask 解析任务业务模型
type ParseTask struct {
	ID          string         `json:"id"`
	ResumeID    string         `json:"resume_id"`
	UserID      string         `json:"user_id"`
	FilePath    string         `json:"file_path"`
	FileType    string         `json:"file_type"`
	Status      string         `json:"status"` // pending, processing, completed, failed
	Progress    int            `json:"progress"`
	Result      *ParsedContent `json:"result,omitempty"`
	ErrorMsg    string         `json:"error_msg,omitempty"`
	Options     *ParseOptions  `json:"options,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	CompletedAt *time.Time     `json:"completed_at,omitempty"`
}

// ParsedContent 解析内容
type ParsedContent struct {
	PersonalInfo *PersonalInfo  `json:"personal_info,omitempty"`
	Education    []*Education   `json:"education,omitempty"`
	Experience   []*Experience  `json:"experience,omitempty"`
	Projects     []*Project     `json:"projects,omitempty"`
	Skills       *Skills        `json:"skills,omitempty"`
	RawText      string         `json:"raw_text,omitempty"`
	Metadata     *ParseMetadata `json:"metadata,omitempty"`
}

// PersonalInfo 个人信息
type PersonalInfo struct {
	Name        string   `json:"name,omitempty"`
	Phone       string   `json:"phone,omitempty"`
	Email       string   `json:"email,omitempty"`
	Address     string   `json:"address,omitempty"`
	BirthDate   string   `json:"birth_date,omitempty"`
	Gender      string   `json:"gender,omitempty"`
	Nationality string   `json:"nationality,omitempty"`
	SocialLinks []string `json:"social_links,omitempty"`
	AvatarURL   string   `json:"avatar_url,omitempty"`
}

// Education 教育背景
type Education struct {
	School      string   `json:"school,omitempty"`
	Degree      string   `json:"degree,omitempty"`
	Major       string   `json:"major,omitempty"`
	StartDate   string   `json:"start_date,omitempty"`
	EndDate     string   `json:"end_date,omitempty"`
	GPA         string   `json:"gpa,omitempty"`
	Description string   `json:"description,omitempty"`
	Courses     []string `json:"courses,omitempty"`
}

// Experience 工作经历
type Experience struct {
	Company          string   `json:"company,omitempty"`
	Position         string   `json:"position,omitempty"`
	StartDate        string   `json:"start_date,omitempty"`
	EndDate          string   `json:"end_date,omitempty"`
	Location         string   `json:"location,omitempty"`
	Department       string   `json:"department,omitempty"`
	Responsibilities []string `json:"responsibilities,omitempty"`
	Achievements     []string `json:"achievements,omitempty"`
	Technologies     []string `json:"technologies,omitempty"`
}

// Project 项目经历
type Project struct {
	Name         string   `json:"name,omitempty"`
	Role         string   `json:"role,omitempty"`
	StartDate    string   `json:"start_date,omitempty"`
	EndDate      string   `json:"end_date,omitempty"`
	Description  string   `json:"description,omitempty"`
	Technologies []string `json:"technologies,omitempty"`
	Achievements []string `json:"achievements,omitempty"`
	URL          string   `json:"url,omitempty"`
	Company      string   `json:"company,omitempty"`
}

// Skills 技能
type Skills struct {
	Categories     []*SkillCategory `json:"categories,omitempty"`
	Languages      []string         `json:"languages,omitempty"`
	Certifications []string         `json:"certifications,omitempty"`
	Awards         []string         `json:"awards,omitempty"`
}

// SkillCategory 技能分类
type SkillCategory struct {
	Category string       `json:"category,omitempty"`
	Skills   []*SkillItem `json:"skills,omitempty"`
}

// SkillItem 技能项
type SkillItem struct {
	Name  string `json:"name,omitempty"`
	Level string `json:"level,omitempty"`
	Years int32  `json:"years,omitempty"`
}

// ParseMetadata 解析元数据
type ParseMetadata struct {
	FileSize        string   `json:"file_size,omitempty"`
	PageCount       int32    `json:"page_count,omitempty"`
	ParseDuration   string   `json:"parse_duration,omitempty"`
	ParserVersion   string   `json:"parser_version,omitempty"`
	Warnings        []string `json:"warnings,omitempty"`
	ConfidenceScore int32    `json:"confidence_score,omitempty"`
}

// ParseTaskRepo 解析任务仓库接口
type ParseTaskRepo interface {
	CreateTask(ctx context.Context, task *ParseTask) (*ParseTask, error)
	GetTask(ctx context.Context, taskID string) (*ParseTask, error)
	UpdateTask(ctx context.Context, task *ParseTask) error
	ListTasksByUser(ctx context.Context, userID string, limit, offset int) ([]*ParseTask, error)
	DeleteTask(ctx context.Context, taskID string) error
}

// DocumentParser 文档解析器接口
type DocumentParser interface {
	Parse(ctx context.Context, filePath string, options *ParseOptions) (*ParsedContent, error)
	SupportedTypes() []string
}

// ParserUsecase 解析用例
type ParserUsecase struct {
	repo    ParseTaskRepo
	parsers map[string]DocumentParser
	log     *log.Helper
}

// NewParserUsecase 创建解析用例
func NewParserUsecase(repo ParseTaskRepo, logger log.Logger) *ParserUsecase {
	return &ParserUsecase{
		repo:    repo,
		parsers: make(map[string]DocumentParser),
		log:     log.NewHelper(logger),
	}
}

// RegisterParser 注册解析器
func (uc *ParserUsecase) RegisterParser(fileType string, parser DocumentParser) {
	uc.parsers[fileType] = parser
}

// ParseDocument 解析文档
func (uc *ParserUsecase) ParseDocument(ctx context.Context, filePath, fileType, resumeID, userID string, options *ParseOptions) (*ParseTask, error) {
	// 验证文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, ErrFileNotFound
	}

	// 检查是否支持该文件类型
	parser, exists := uc.parsers[fileType]
	if !exists {
		return nil, ErrUnsupportedType
	}

	// 创建解析任务
	task := &ParseTask{
		ID:        uuid.New().String(),
		ResumeID:  resumeID,
		UserID:    userID,
		FilePath:  filePath,
		FileType:  fileType,
		Status:    "pending",
		Progress:  0,
		Options:   options,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 保存任务
	task, err := uc.repo.CreateTask(ctx, task)
	if err != nil {
		return nil, err
	}

	// 异步处理解析
	go uc.processParseTask(context.Background(), task, parser)

	return task, nil
}

// processParseTask 处理解析任务
func (uc *ParserUsecase) processParseTask(ctx context.Context, task *ParseTask, parser DocumentParser) {
	// 更新状态为处理中
	task.Status = "processing"
	task.Progress = 10
	task.UpdatedAt = time.Now()
	uc.repo.UpdateTask(ctx, task)

	startTime := time.Now()

	// 执行解析
	content, err := parser.Parse(ctx, task.FilePath, task.Options)
	if err != nil {
		// 解析失败
		task.Status = "failed"
		task.ErrorMsg = err.Error()
		task.Progress = 0
		task.UpdatedAt = time.Now()
		uc.repo.UpdateTask(ctx, task)
		uc.log.Errorf("Parse failed for task %s: %v", task.ID, err)
		return
	}

	// 添加元数据
	if content.Metadata == nil {
		content.Metadata = &ParseMetadata{}
	}
	content.Metadata.ParseDuration = time.Since(startTime).String()
	content.Metadata.ParserVersion = "1.0.0"

	// 文件信息
	if fileInfo, err := os.Stat(task.FilePath); err == nil {
		content.Metadata.FileSize = fmt.Sprintf("%d", fileInfo.Size())
	}

	// 计算置信度
	content.Metadata.ConfidenceScore = uc.calculateConfidence(content)

	// 解析成功
	task.Status = "completed"
	task.Progress = 100
	task.Result = content
	task.UpdatedAt = time.Now()

	if err := uc.repo.UpdateTask(ctx, task); err != nil {
		uc.log.Errorf("Failed to update task %s: %v", task.ID, err)
	}

	uc.log.Infof("Parse completed for task %s", task.ID)
}

// calculateConfidence 计算解析置信度
func (uc *ParserUsecase) calculateConfidence(content *ParsedContent) int32 {
	score := int32(0)

	// 个人信息完整度
	if content.PersonalInfo != nil {
		if content.PersonalInfo.Name != "" {
			score += 20
		}
		if content.PersonalInfo.Email != "" {
			score += 15
		}
		if content.PersonalInfo.Phone != "" {
			score += 10
		}
	}

	// 工作经历
	if len(content.Experience) > 0 {
		score += 25
		for _, exp := range content.Experience {
			if exp.Company != "" && exp.Position != "" {
				score += 5
				break
			}
		}
	}

	// 教育背景
	if len(content.Education) > 0 {
		score += 15
		for _, edu := range content.Education {
			if edu.School != "" && edu.Degree != "" {
				score += 5
				break
			}
		}
	}

	// 技能信息
	if content.Skills != nil && len(content.Skills.Categories) > 0 {
		score += 10
	}

	// 原始文本长度
	if len(content.RawText) > 100 {
		score += 5
	}

	if score > 100 {
		score = 100
	}

	return score
}

// GetParseStatus 获取解析状态
func (uc *ParserUsecase) GetParseStatus(ctx context.Context, taskID string) (*ParseTask, error) {
	return uc.repo.GetTask(ctx, taskID)
}

// ListUserTasks 获取用户的解析任务列表
func (uc *ParserUsecase) ListUserTasks(ctx context.Context, userID string, limit, offset int) ([]*ParseTask, error) {
	return uc.repo.ListTasksByUser(ctx, userID, limit, offset)
}

// CleanText 清洗文本
func (uc *ParserUsecase) CleanText(text string) string {
	// 移除多余空白
	text = strings.TrimSpace(text)

	// 统一换行符
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")

	// 移除多余的换行符
	for strings.Contains(text, "\n\n\n") {
		text = strings.ReplaceAll(text, "\n\n\n", "\n\n")
	}

	// 移除多余的空格
	for strings.Contains(text, "  ") {
		text = strings.ReplaceAll(text, "  ", " ")
	}

	return text
}

// ExtractFileExtension 提取文件扩展名
func ExtractFileExtension(filePath string) string {
	ext := filepath.Ext(filePath)
	if len(ext) > 0 {
		return strings.ToLower(ext[1:]) // 移除点号
	}
	return ""
}
