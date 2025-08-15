package eino

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/lyb88999/resume_helper/backend/services/ai-service/internal/conf"
)

// EinoComponents Eino组件集合（简化实现）
type EinoComponents struct {
	ChatModel     ChatModel
	Embedding     EmbeddingModel
	ParsingChain  *ResumeParsingChain
	AnalysisGraph *AnalysisGraph
	logger        *log.Helper
}

// ChatModel 聊天模型接口
type ChatModel interface {
	Generate(ctx context.Context, messages []Message, options ...GenerateOption) (*GenerateResponse, error)
}

// EmbeddingModel 嵌入模型接口
type EmbeddingModel interface {
	Embed(ctx context.Context, texts []string) ([][]float64, error)
}

// Message 消息结构
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// GenerateResponse 生成响应
type GenerateResponse struct {
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice 选择
type Choice struct {
	Message Message `json:"message"`
}

// Usage 使用统计
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// GenerateOption 生成选项
type GenerateOption func(*GenerateOptions)

// GenerateOptions 生成选项
type GenerateOptions struct {
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
}

// WithMaxTokens 设置最大token数
func WithMaxTokens(maxTokens int) GenerateOption {
	return func(opts *GenerateOptions) {
		opts.MaxTokens = maxTokens
	}
}

// WithTemperature 设置温度
func WithTemperature(temperature float64) GenerateOption {
	return func(opts *GenerateOptions) {
		opts.Temperature = temperature
	}
}

// AnalysisResult 分析结果
type AnalysisResult struct {
	ID             string         `json:"id"`
	ResumeID       string         `json:"resume_id"`
	TargetPosition string         `json:"target_position"`
	Scores         ScoreBreakdown `json:"scores"`
	Suggestions    []Suggestion   `json:"suggestions"`
	Summary        string         `json:"summary"`
	AnalyzedAt     time.Time      `json:"analyzed_at"`
}

// ScoreBreakdown 评分详情
type ScoreBreakdown struct {
	OverallScore        float64            `json:"overall_score"`
	CompletenessScore   float64            `json:"completeness_score"`
	ClarityScore        float64            `json:"clarity_score"`
	KeywordScore        float64            `json:"keyword_score"`
	FormatScore         float64            `json:"format_score"`
	QuantificationScore float64            `json:"quantification_score"`
	DimensionScores     map[string]float64 `json:"dimension_scores"`
}

// Suggestion 建议
type Suggestion struct {
	ID          string   `json:"id"`
	Type        string   `json:"type"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Priority    string   `json:"priority"`
	Section     string   `json:"section"`
	Action      string   `json:"action"`
	Examples    []string `json:"examples"`
}

// NewEinoComponents 创建Eino组件集合
func NewEinoComponents(aiConfig *conf.AI, logger log.Logger) (*EinoComponents, error) {
	helper := log.NewHelper(logger)

	components := &EinoComponents{
		logger: helper,
	}

	// 初始化ChatModel
	if err := components.initChatModel(aiConfig.Model); err != nil {
		return nil, fmt.Errorf("初始化ChatModel失败: %w", err)
	}

	// 初始化Embedding
	if err := components.initEmbedding(aiConfig.Embedding); err != nil {
		return nil, fmt.Errorf("初始化Embedding失败: %w", err)
	}

	// 初始化文档处理组件
	if err := components.initDocumentComponents(); err != nil {
		return nil, fmt.Errorf("初始化文档组件失败: %w", err)
	}

	// 初始化高级组件
	if err := components.initAdvancedComponents(aiConfig); err != nil {
		return nil, fmt.Errorf("初始化高级组件失败: %w", err)
	}

	helper.Info("Eino组件初始化完成")
	return components, nil
}

// initChatModel 初始化聊天模型
func (c *EinoComponents) initChatModel(config *conf.ModelConfig) error {
	switch config.Provider {
	case "ark":
		chatModel := &ARKChatModel{
			APIKey:      config.ApiKey,
			BaseURL:     config.BaseUrl,
			Model:       config.ModelName,
			MaxTokens:   int(config.MaxTokens),
			Temperature: config.Temperature,
			Timeout:     time.Duration(config.TimeoutSeconds) * time.Second,
		}
		c.ChatModel = chatModel
		c.logger.Infof("已初始化ARK ChatModel: %s", config.ModelName)

	default:
		return fmt.Errorf("不支持的模型提供商: %s", config.Provider)
	}

	return nil
}

// initEmbedding 初始化嵌入模型
func (c *EinoComponents) initEmbedding(config *conf.EmbeddingConfig) error {
	switch config.Provider {
	case "ark":
		embedder := &ARKEmbeddingModel{
			APIKey:  config.ApiKey,
			BaseURL: config.BaseUrl,
			Model:   config.ModelName,
			Timeout: 30 * time.Second,
		}
		c.Embedding = embedder
		c.logger.Infof("已初始化ARK Embedding: %s", config.ModelName)

	default:
		return fmt.Errorf("不支持的嵌入提供商: %s", config.Provider)
	}

	return nil
}

// initDocumentComponents 初始化文档处理组件（暂时跳过）
func (c *EinoComponents) initDocumentComponents() error {
	c.logger.Info("文档处理组件暂时跳过初始化")
	return nil
}

// initAdvancedComponents 初始化高级组件
func (c *EinoComponents) initAdvancedComponents(config *conf.AI) error {
	// 初始化简历解析Chain
	c.ParsingChain = NewResumeParsingChain(
		c.ChatModel,
		c.logger,
	)

	// 初始化分析Graph
	c.AnalysisGraph = NewAnalysisGraph(
		c.ChatModel,
		c.logger,
	)

	c.logger.Info("已初始化高级组件")
	return nil
}

// ResumeParsingChain 简历解析链
type ResumeParsingChain struct {
	chatModel ChatModel
	logger    *log.Helper
}

// NewResumeParsingChain 创建简历解析链
func NewResumeParsingChain(
	chatModel ChatModel,
	logger *log.Helper,
) *ResumeParsingChain {
	return &ResumeParsingChain{
		chatModel: chatModel,
		logger:    logger,
	}
}

// Execute 执行简历解析
func (c *ResumeParsingChain) Execute(ctx context.Context, filePath string) (*ResumeData, error) {
	c.logger.WithContext(ctx).Infof("开始解析简历文件: %s", filePath)

	// 暂时使用模拟内容，实际应该读取文件
	content := "这是一个示例简历内容..."

	// 使用大模型进行结构化提取
	resumeData, err := c.extractStructuredData(ctx, content)
	if err != nil {
		return nil, fmt.Errorf("结构化提取失败: %w", err)
	}

	c.logger.WithContext(ctx).Info("简历解析完成")
	return resumeData, nil
}

// extractStructuredData 提取结构化数据
func (c *ResumeParsingChain) extractStructuredData(ctx context.Context, content string) (*ResumeData, error) {
	// 构建提示
	messages := c.buildExtractionPrompt(content)

	// 调用大模型
	resp, err := c.chatModel.Generate(ctx, messages, WithMaxTokens(4096))
	if err != nil {
		return nil, fmt.Errorf("大模型调用失败: %w", err)
	}

	// 解析响应
	resumeData, err := c.parseModelResponse(resp.Choices[0].Message.Content)
	if err != nil {
		return nil, fmt.Errorf("响应解析失败: %w", err)
	}

	return resumeData, nil
}

// buildExtractionPrompt 构建提取提示
func (c *ResumeParsingChain) buildExtractionPrompt(content string) []Message {
	systemPrompt := `你是一个专业的简历解析专家。请将以下简历内容解析为结构化的JSON格式。

请识别并提取以下章节：
1. 个人信息 (personal_info)：姓名、联系方式、邮箱、电话、地址、LinkedIn、GitHub等
2. 教育背景 (education)：学校、专业、学历、时间、GPA、相关课程、荣誉等  
3. 工作经历 (experience)：公司、职位、时间、地点、工作内容、成就、使用技术等
4. 项目经历 (projects)：项目名称、角色、时间、描述、技术栈、成果、项目链接等
5. 技能特长 (skills)：技术技能、编程语言、框架、工具、软技能等
6. 其他信息 (others)：获奖经历、证书、兴趣爱好、志愿经历等

请严格按照JSON格式输出，确保数据的准确性和完整性。时间格式使用ISO 8601标准。`

	userPrompt := fmt.Sprintf("请解析以下简历内容：\n\n%s", content)

	return []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}
}

// parseModelResponse 解析模型响应
func (c *ResumeParsingChain) parseModelResponse(content string) (*ResumeData, error) {
	// 这里需要实现JSON解析逻辑
	// 简化实现，实际应该完整解析JSON
	resumeData := &ResumeData{
		ID:      generateID(),
		Version: "1.0",
		PersonalInfo: PersonalInfo{
			Name: "从模型提取的姓名",
		},
		// 其他字段的解析...
	}

	return resumeData, nil
}

// AnalysisGraph 分析图
type AnalysisGraph struct {
	chatModel ChatModel
	logger    *log.Helper
}

// NewAnalysisGraph 创建分析图
func NewAnalysisGraph(
	chatModel ChatModel,
	logger *log.Helper,
) *AnalysisGraph {
	return &AnalysisGraph{
		chatModel: chatModel,
		logger:    logger,
	}
}

// buildGraph 构建分析图（简化实现）
func (g *AnalysisGraph) buildGraph() {
	// 简化的图构建逻辑
	g.logger.Info("分析图构建完成")
}

// Execute 执行分析图
func (g *AnalysisGraph) Execute(ctx context.Context, resumeData *ResumeData, targetPosition string) (*AnalysisResult, error) {
	g.logger.WithContext(ctx).Infof("开始执行智能分析，目标职位: %s", targetPosition)

	// 简化的分析流程
	// 1. 完整性分析
	completenessScore := g.analyzeCompleteness(ctx, resumeData)

	// 2. 清晰度分析
	clarityScore := g.analyzeClarity(ctx, resumeData)

	// 3. 关键词分析
	keywordScore := g.analyzeKeywords(ctx, resumeData, targetPosition)

	// 4. 格式分析
	formatScore := g.analyzeFormat(ctx, resumeData)

	// 5. 量化分析
	quantificationScore := g.analyzeQuantification(ctx, resumeData)

	// 计算总分
	overallScore := (completenessScore + clarityScore + keywordScore + formatScore + quantificationScore) / 5

	// 生成建议
	suggestions := g.generateSuggestions(ctx, overallScore, targetPosition)

	// 构建分析结果
	analysisResult := &AnalysisResult{
		ID:             fmt.Sprintf("analysis_%s", resumeData.ID),
		ResumeID:       resumeData.ID,
		TargetPosition: targetPosition,
		Scores: ScoreBreakdown{
			OverallScore:        overallScore,
			CompletenessScore:   completenessScore,
			ClarityScore:        clarityScore,
			KeywordScore:        keywordScore,
			FormatScore:         formatScore,
			QuantificationScore: quantificationScore,
		},
		Suggestions: suggestions,
		Summary:     fmt.Sprintf("简历整体质量为%.1f分，建议重点关注%s方面的优化。", overallScore, g.getWeakestArea(completenessScore, clarityScore, keywordScore, formatScore, quantificationScore)),
		AnalyzedAt:  time.Now(),
	}

	g.logger.WithContext(ctx).Info("智能分析执行完成")
	return analysisResult, nil
}

// 分析方法实现
func (g *AnalysisGraph) analyzeCompleteness(ctx context.Context, resumeData *ResumeData) float64 {
	// 完整性分析：检查必要字段是否存在
	score := 100.0

	if resumeData.PersonalInfo.Name == "" {
		score -= 20
	}
	if resumeData.PersonalInfo.Email == "" {
		score -= 15
	}
	if len(resumeData.Experience) == 0 {
		score -= 25
	}
	if len(resumeData.Education) == 0 {
		score -= 20
	}
	if len(resumeData.Skills.Technical) == 0 {
		score -= 20
	}

	return score
}

func (g *AnalysisGraph) analyzeClarity(ctx context.Context, resumeData *ResumeData) float64 {
	// 清晰度分析：检查描述的清晰程度
	score := 80.0

	// 检查工作经历描述
	for _, exp := range resumeData.Experience {
		if len(exp.Description) == 0 {
			score -= 10
		}
	}

	return score
}

func (g *AnalysisGraph) analyzeKeywords(ctx context.Context, resumeData *ResumeData, targetPosition string) float64 {
	// 关键词分析：与目标职位的匹配度
	if targetPosition == "" {
		return 70.0
	}

	positionLower := strings.ToLower(targetPosition)
	matchCount := 0
	totalSkills := len(resumeData.Skills.Technical) + len(resumeData.Skills.Frameworks)

	for _, skill := range resumeData.Skills.Technical {
		if strings.Contains(positionLower, strings.ToLower(skill)) {
			matchCount++
		}
	}

	if totalSkills == 0 {
		return 50.0
	}

	return float64(matchCount) / float64(totalSkills) * 100
}

func (g *AnalysisGraph) analyzeFormat(ctx context.Context, resumeData *ResumeData) float64 {
	// 格式分析：检查数据结构的规范性
	score := 85.0

	// 检查联系方式格式
	if resumeData.PersonalInfo.Email != "" && !strings.Contains(resumeData.PersonalInfo.Email, "@") {
		score -= 10
	}

	return score
}

func (g *AnalysisGraph) analyzeQuantification(ctx context.Context, resumeData *ResumeData) float64 {
	// 量化分析：检查描述中的数字化程度
	quantifiedCount := 0
	totalDescriptions := 0

	for _, exp := range resumeData.Experience {
		totalDescriptions += len(exp.Description)
		for _, desc := range exp.Description {
			if strings.ContainsAny(desc, "0123456789%") {
				quantifiedCount++
			}
		}
	}

	if totalDescriptions == 0 {
		return 60.0
	}

	return float64(quantifiedCount) / float64(totalDescriptions) * 100
}

func (g *AnalysisGraph) generateSuggestions(ctx context.Context, overallScore float64, targetPosition string) []Suggestion {
	suggestions := []Suggestion{
		{
			ID:          "suggestion_1",
			Type:        "content",
			Title:       "优化工作经历描述",
			Description: "建议在工作经历中添加更多量化数据来展示具体成果",
			Priority:    "high",
			Section:     "experience",
			Action:      "在每个工作经历中添加具体的数字和百分比",
		},
	}

	if overallScore < 70 {
		suggestions = append(suggestions, Suggestion{
			ID:          "suggestion_2",
			Type:        "structure",
			Title:       "完善简历结构",
			Description: "建议补充缺失的关键信息，如联系方式、教育背景等",
			Priority:    "high",
			Section:     "overall",
			Action:      "检查并补充所有必要的简历章节",
		})
	}

	if targetPosition != "" {
		suggestions = append(suggestions, Suggestion{
			ID:          "suggestion_3",
			Type:        "keywords",
			Title:       "增加相关技能关键词",
			Description: fmt.Sprintf("针对%s职位，建议增加更多相关的技能关键词", targetPosition),
			Priority:    "medium",
			Section:     "skills",
			Action:      "研究目标职位要求，补充相关技能",
		})
	}

	return suggestions
}

func (g *AnalysisGraph) getWeakestArea(completeness, clarity, keyword, format, quantification float64) string {
	scores := map[string]float64{
		"完整性": completeness,
		"清晰度": clarity,
		"关键词": keyword,
		"格式":  format,
		"量化":  quantification,
	}

	weakest := "完整性"
	minScore := completeness

	for area, score := range scores {
		if score < minScore {
			minScore = score
			weakest = area
		}
	}

	return weakest
}

// 工具函数
func generateID() string {
	return fmt.Sprintf("resume_%d", time.Now().Unix())
}
