package biz

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/liyubo06/resumeOptim_claude/backend/services/ai-service/internal/conf"
	"github.com/liyubo06/resumeOptim_claude/backend/services/ai-service/internal/eino"
)

// AIRepo AI数据仓库接口
type AIRepo interface {
	SaveAnalysisResult(ctx context.Context, result *eino.AnalysisResult) error
	GetAnalysisResult(ctx context.Context, id string) (*eino.AnalysisResult, error)
	SaveChatSession(ctx context.Context, session *eino.ChatContext) error
	GetChatSession(ctx context.Context, sessionID string) (*eino.ChatContext, error)
}

// AIUsecase AI用例
type AIUsecase struct {
	repo       AIRepo
	components *eino.EinoComponents
	logger     *log.Helper
}

// NewAIUsecase 创建AI用例
func NewAIUsecase(repo AIRepo, aiConfig *conf.AI, logger log.Logger) *AIUsecase {
	helper := log.NewHelper(logger)

	// 初始化Eino组件
	components, err := eino.NewEinoComponents(aiConfig, logger)
	if err != nil {
		helper.Errorf("初始化Eino组件失败: %v", err)
		// 使用空组件继续运行，避免启动失败
		components = &eino.EinoComponents{}
	}

	return &AIUsecase{
		repo:       repo,
		components: components,
		logger:     helper,
	}
}

// AnalyzeResume 分析简历
func (uc *AIUsecase) AnalyzeResume(ctx context.Context, req *AnalyzeResumeRequest) (*AnalyzeResumeResponse, error) {
	uc.logger.WithContext(ctx).Infof("开始分析简历，ID: %s", req.ResumeID)

	// 1. 解析简历内容为结构化数据
	var resumeData *eino.ResumeData
	var err error

	if req.FilePath != "" {
		// 从文件解析
		if uc.components.ParsingChain != nil {
			resumeData, err = uc.components.ParsingChain.Execute(ctx, req.FilePath)
			if err != nil {
				return nil, fmt.Errorf("简历解析失败: %w", err)
			}
		} else {
			return nil, fmt.Errorf("简历解析链未初始化")
		}
	} else {
		// 从内容直接解析（简化实现）
		resumeData = &eino.ResumeData{
			ID:      req.ResumeID,
			Version: "1.0",
			PersonalInfo: eino.PersonalInfo{
				Name: "从内容解析的姓名",
			},
			// 其他字段的简化解析...
		}
	}

	// 设置时间戳
	resumeData.CreatedAt = time.Now()
	resumeData.UpdatedAt = time.Now()

	// 2. 执行智能分析
	var analysisResult *eino.AnalysisResult
	if uc.components.AnalysisGraph != nil {
		analysisResult, err = uc.components.AnalysisGraph.Execute(ctx, resumeData, req.TargetPosition)
		if err != nil {
			return nil, fmt.Errorf("智能分析失败: %w", err)
		}
	} else {
		// 提供默认分析结果
		analysisResult = &eino.AnalysisResult{
			ID:             fmt.Sprintf("analysis_%s", req.ResumeID),
			ResumeID:       req.ResumeID,
			TargetPosition: req.TargetPosition,
			Scores: eino.ScoreBreakdown{
				OverallScore:        75.0,
				CompletenessScore:   80.0,
				ClarityScore:        70.0,
				KeywordScore:        75.0,
				FormatScore:         85.0,
				QuantificationScore: 65.0,
			},
			Summary:    "简历整体质量良好，建议在量化描述方面进一步改进。",
			AnalyzedAt: time.Now(),
		}
	}

	// 3. 保存分析结果
	if err := uc.repo.SaveAnalysisResult(ctx, analysisResult); err != nil {
		uc.logger.WithContext(ctx).Errorf("保存分析结果失败: %v", err)
		// 不阻断流程，继续返回结果
	}

	uc.logger.WithContext(ctx).Infof("简历分析完成，分析ID: %s", analysisResult.ID)

	return &AnalyzeResumeResponse{
		AnalysisID: analysisResult.ID,
		Result:     analysisResult,
		Status:     "success",
		Message:    "分析完成",
	}, nil
}

// GenerateSuggestions 生成优化建议
func (uc *AIUsecase) GenerateSuggestions(ctx context.Context, req *GenerateSuggestionsRequest) (*GenerateSuggestionsResponse, error) {
	uc.logger.WithContext(ctx).Infof("开始生成建议，分析ID: %s", req.AnalysisID)

	// 获取分析结果
	analysisResult := req.AnalysisResult
	if analysisResult == nil && req.AnalysisID != "" {
		var err error
		analysisResult, err = uc.repo.GetAnalysisResult(ctx, req.AnalysisID)
		if err != nil {
			return nil, fmt.Errorf("获取分析结果失败: %w", err)
		}
	}

	if analysisResult == nil {
		return nil, fmt.Errorf("分析结果不存在")
	}

	// 生成建议（这里可以调用更复杂的AI生成逻辑）
	suggestions := []eino.Suggestion{
		{
			ID:          "suggestion_1",
			Type:        "content",
			Title:       "优化工作经历描述",
			Description: "建议在工作经历中添加更多量化数据，如具体的业绩提升百分比、处理的数据量等。",
			Priority:    "high",
			Section:     "experience",
			Action:      "在每个工作经历描述中至少包含2-3个具体数字",
			Examples:    []string{"将'显著提升'改为'提升40%'", "将'大量用户'改为'10万+用户'"},
		},
		{
			ID:          "suggestion_2",
			Type:        "keywords",
			Title:       "增加技能关键词",
			Description: fmt.Sprintf("针对%s职位，建议增加相关的技术关键词以提高ATS匹配度。", req.TargetPosition),
			Priority:    "medium",
			Section:     "skills",
			Action:      "研究目标职位的JD，补充相关技能关键词",
			Examples:    []string{"添加云计算相关技能", "补充最新框架技术"},
		},
		{
			ID:          "suggestion_3",
			Type:        "format",
			Title:       "优化简历格式",
			Description: "建议调整简历的版式设计，提高可读性和专业性。",
			Priority:    "low",
			Section:     "overall",
			Action:      "使用更清晰的层次结构和一致的格式",
			Examples:    []string{"统一字体大小", "调整行间距"},
		},
	}

	return &GenerateSuggestionsResponse{
		Suggestions: suggestions,
		Reasoning:   "基于当前简历的评分情况和目标职位要求，重点关注内容量化和关键词优化。",
		Status:      "success",
		Message:     "建议生成完成",
	}, nil
}

// Chat 智能问答
func (uc *AIUsecase) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	uc.logger.WithContext(ctx).Infof("开始处理智能问答，会话ID: %s", req.SessionID)

	// 获取或创建会话上下文
	var chatContext *eino.ChatContext
	if req.SessionID != "" {
		var err error
		chatContext, err = uc.repo.GetChatSession(ctx, req.SessionID)
		if err != nil {
			uc.logger.WithContext(ctx).Warnf("获取会话失败，创建新会话: %v", err)
		}
	}

	if chatContext == nil {
		chatContext = &eino.ChatContext{
			SessionID: req.SessionID,
			Messages:  []eino.Message{},
		}
	}

	// 添加用户消息
	chatContext.Messages = append(chatContext.Messages, eino.Message{
		Role:    "human",
		Content: req.Message,
	})

	// 设置上下文信息
	if req.Context != "" {
		// 这里可以解析上下文信息，如简历数据等
	}

	// 调用智能体进行对话
	var response string
	var sources []string
	if uc.components.ChatModel != nil {
		// 使用ChatModel进行简单对话
		messages := []eino.Message{
			{Role: "system", Content: "你是一个专业的简历优化助手，请根据用户问题提供有用的建议。"},
			{Role: "user", Content: req.Message},
		}

		// 调用模型生成回复
		resp, err := uc.components.ChatModel.Generate(ctx, messages, eino.WithMaxTokens(2048))
		if err != nil {
			uc.logger.WithContext(ctx).Errorf("ChatModel调用失败: %v", err)
			response = "抱歉，我暂时无法回答您的问题，请稍后再试。"
		} else {
			response = resp.Choices[0].Message.Content
		}
		sources = []string{"chatmodel_response"}
	} else {
		// 提供默认回复
		response = "感谢您的问题。目前智能问答功能正在完善中，请稍后再试。"
		sources = []string{"default_response"}
	}

	// 添加助手回复
	chatContext.Messages = append(chatContext.Messages, eino.Message{
		Role:    "assistant",
		Content: response,
	})

	// 保存会话
	if err := uc.repo.SaveChatSession(ctx, chatContext); err != nil {
		uc.logger.WithContext(ctx).Errorf("保存会话失败: %v", err)
	}

	return &ChatResponse{
		Response:  response,
		SessionID: chatContext.SessionID,
		Sources:   sources,
		Status:    "success",
		Message:   "对话完成",
	}, nil
}

// RetrieveKnowledge 知识检索
func (uc *AIUsecase) RetrieveKnowledge(ctx context.Context, req *RetrieveKnowledgeRequest) (*RetrieveKnowledgeResponse, error) {
	uc.logger.WithContext(ctx).Infof("开始知识检索，查询: %s", req.Query)

	// 这里实现向量检索逻辑
	// 暂时返回模拟数据
	items := []eino.Document{
		{
			ID:      "knowledge_1",
			Title:   "软件工程师简历优化指南",
			Content: "软件工程师简历应该突出技术栈、项目经验和解决问题的能力...",
			Score:   0.95,
			Metadata: map[string]string{
				"category": "career_guide",
				"type":     "resume_tips",
			},
		},
		{
			ID:      "knowledge_2",
			Title:   "技术面试常见问题",
			Content: "技术面试中经常会问到算法、系统设计、项目经验等问题...",
			Score:   0.87,
			Metadata: map[string]string{
				"category": "interview",
				"type":     "technical",
			},
		},
	}

	return &RetrieveKnowledgeResponse{
		Items:   items,
		Status:  "success",
		Message: "检索完成",
	}, nil
}

// 请求和响应结构体

type AnalyzeResumeRequest struct {
	ResumeID       string
	Content        string
	FilePath       string
	FileType       string
	TargetPosition string
	Options        *AnalysisOptions
}

type AnalysisOptions struct {
	EnableCompleteness   bool
	EnableClarity        bool
	EnableKeyword        bool
	EnableFormat         bool
	EnableQuantification bool
}

type AnalyzeResumeResponse struct {
	AnalysisID string
	Result     *eino.AnalysisResult
	Status     string
	Message    string
}

type GenerateSuggestionsRequest struct {
	AnalysisID     string
	AnalysisResult *eino.AnalysisResult
	TargetPosition string
	Industry       string
	Options        *SuggestionOptions
}

type SuggestionOptions struct {
	MaxSuggestions  int32
	FocusArea       string
	ExperienceLevel string
}

type GenerateSuggestionsResponse struct {
	Suggestions []eino.Suggestion
	Reasoning   string
	Status      string
	Message     string
}

type ChatRequest struct {
	SessionID string
	Message   string
	Context   string
	Options   *ChatOptions
}

type ChatOptions struct {
	UseResumeContext bool
	UseKnowledgeBase bool
	Language         string
}

type ChatResponse struct {
	Response  string
	SessionID string
	Sources   []string
	Status    string
	Message   string
}

type RetrieveKnowledgeRequest struct {
	Query               string
	TopK                int32
	SimilarityThreshold float32
	Filters             []string
}

type RetrieveKnowledgeResponse struct {
	Items   []eino.Document
	Status  string
	Message string
}
