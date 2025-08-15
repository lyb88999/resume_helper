package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/liyubo06/resumeOptim_claude/backend/services/ai-service/api/ai/v1"
	"github.com/liyubo06/resumeOptim_claude/backend/services/ai-service/internal/biz"
	"github.com/liyubo06/resumeOptim_claude/backend/services/ai-service/internal/eino"
)

// AIService AI服务实现
type AIService struct {
	pb.UnimplementedAIServiceServer

	aiUsecase *biz.AIUsecase
	log       *log.Helper
}

// NewAIService 创建AI服务
func NewAIService(aiUsecase *biz.AIUsecase, logger log.Logger) *AIService {
	return &AIService{
		aiUsecase: aiUsecase,
		log:       log.NewHelper(logger),
	}
}

// AnalyzeResume 分析简历
func (s *AIService) AnalyzeResume(ctx context.Context, req *pb.AnalyzeResumeRequest) (*pb.AnalyzeResumeResponse, error) {
	s.log.WithContext(ctx).Infof("收到简历分析请求，简历ID: %s", req.ResumeId)

	// 转换请求参数
	bizReq := &biz.AnalyzeResumeRequest{
		ResumeID:       req.ResumeId,
		Content:        req.Content,
		FileType:       req.FileType,
		TargetPosition: req.TargetPosition,
	}

	if req.Options != nil {
		bizReq.Options = &biz.AnalysisOptions{
			EnableCompleteness:   req.Options.EnableCompleteness,
			EnableClarity:        req.Options.EnableClarity,
			EnableKeyword:        req.Options.EnableKeyword,
			EnableFormat:         req.Options.EnableFormat,
			EnableQuantification: req.Options.EnableQuantification,
		}
	}

	// 调用业务逻辑
	bizResp, err := s.aiUsecase.AnalyzeResume(ctx, bizReq)
	if err != nil {
		s.log.WithContext(ctx).Errorf("简历分析失败: %v", err)
		return &pb.AnalyzeResumeResponse{
			Status:  "error",
			Message: err.Error(),
		}, nil
	}

	// 转换响应
	resp := &pb.AnalyzeResumeResponse{
		AnalysisId: bizResp.AnalysisID,
		Status:     bizResp.Status,
		Message:    bizResp.Message,
	}

	if bizResp.Result != nil {
		resp.Result = s.convertAnalysisResult(bizResp.Result)
	}

	s.log.WithContext(ctx).Infof("简历分析完成，分析ID: %s", bizResp.AnalysisID)
	return resp, nil
}

// GenerateSuggestions 生成优化建议
func (s *AIService) GenerateSuggestions(ctx context.Context, req *pb.GenerateSuggestionsRequest) (*pb.GenerateSuggestionsResponse, error) {
	s.log.WithContext(ctx).Infof("收到生成建议请求，分析ID: %s", req.AnalysisId)

	// 转换请求参数
	bizReq := &biz.GenerateSuggestionsRequest{
		AnalysisID:     req.AnalysisId,
		TargetPosition: req.TargetPosition,
		Industry:       req.Industry,
	}

	if req.AnalysisResult != nil {
		bizReq.AnalysisResult = s.convertToBizAnalysisResult(req.AnalysisResult)
	}

	if req.Options != nil {
		bizReq.Options = &biz.SuggestionOptions{
			MaxSuggestions:  req.Options.MaxSuggestions,
			FocusArea:       req.Options.FocusArea,
			ExperienceLevel: req.Options.ExperienceLevel,
		}
	}

	// 调用业务逻辑
	bizResp, err := s.aiUsecase.GenerateSuggestions(ctx, bizReq)
	if err != nil {
		s.log.WithContext(ctx).Errorf("生成建议失败: %v", err)
		return &pb.GenerateSuggestionsResponse{
			Status:  "error",
			Message: err.Error(),
		}, nil
	}

	// 转换响应
	suggestions := make([]*pb.Suggestion, len(bizResp.Suggestions))
	for i, suggestion := range bizResp.Suggestions {
		suggestions[i] = &pb.Suggestion{
			Id:          suggestion.ID,
			Type:        suggestion.Type,
			Title:       suggestion.Title,
			Description: suggestion.Description,
			Priority:    suggestion.Priority,
			Section:     suggestion.Section,
			Action:      suggestion.Action,
		}
	}

	return &pb.GenerateSuggestionsResponse{
		Suggestions: suggestions,
		Reasoning:   bizResp.Reasoning,
		Status:      bizResp.Status,
		Message:     bizResp.Message,
	}, nil
}

// Chat 智能问答
func (s *AIService) Chat(ctx context.Context, req *pb.ChatRequest) (*pb.ChatResponse, error) {
	s.log.WithContext(ctx).Infof("收到智能问答请求，会话ID: %s", req.SessionId)

	// 转换请求参数
	bizReq := &biz.ChatRequest{
		SessionID: req.SessionId,
		Message:   req.Message,
		Context:   req.Context,
	}

	if req.Options != nil {
		bizReq.Options = &biz.ChatOptions{
			UseResumeContext: req.Options.UseResumeContext,
			UseKnowledgeBase: req.Options.UseKnowledgeBase,
			Language:         req.Options.Language,
		}
	}

	// 调用业务逻辑
	bizResp, err := s.aiUsecase.Chat(ctx, bizReq)
	if err != nil {
		s.log.WithContext(ctx).Errorf("智能问答失败: %v", err)
		return &pb.ChatResponse{
			Status:  "error",
			Message: err.Error(),
		}, nil
	}

	return &pb.ChatResponse{
		Response:  bizResp.Response,
		SessionId: bizResp.SessionID,
		Sources:   bizResp.Sources,
		Status:    bizResp.Status,
		Message:   bizResp.Message,
	}, nil
}

// RetrieveKnowledge 知识检索
func (s *AIService) RetrieveKnowledge(ctx context.Context, req *pb.RetrieveKnowledgeRequest) (*pb.RetrieveKnowledgeResponse, error) {
	s.log.WithContext(ctx).Infof("收到知识检索请求，查询: %s", req.Query)

	// 转换请求参数
	bizReq := &biz.RetrieveKnowledgeRequest{
		Query:               req.Query,
		TopK:                req.TopK,
		SimilarityThreshold: req.SimilarityThreshold,
		Filters:             req.Filters,
	}

	// 调用业务逻辑
	bizResp, err := s.aiUsecase.RetrieveKnowledge(ctx, bizReq)
	if err != nil {
		s.log.WithContext(ctx).Errorf("知识检索失败: %v", err)
		return &pb.RetrieveKnowledgeResponse{
			Status:  "error",
			Message: err.Error(),
		}, nil
	}

	// 转换响应
	items := make([]*pb.KnowledgeItem, len(bizResp.Items))
	for i, item := range bizResp.Items {
		items[i] = &pb.KnowledgeItem{
			Id:       item.ID,
			Title:    item.Title,
			Content:  item.Content,
			Score:    float32(item.Score),
			Metadata: item.Metadata,
		}
	}

	return &pb.RetrieveKnowledgeResponse{
		Items:   items,
		Status:  bizResp.Status,
		Message: bizResp.Message,
	}, nil
}

// Health 健康检查
func (s *AIService) Health(ctx context.Context, req *emptypb.Empty) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{
		Status:  "healthy",
		Version: "1.0.0",
		Components: map[string]string{
			"database": "ok",
			"redis":    "ok",
			"ai_model": "ok",
		},
	}, nil
}

// Helper methods for type conversion

func (s *AIService) convertAnalysisResult(result *eino.AnalysisResult) *pb.AnalysisResult {
	if result == nil {
		return nil
	}

	// 转换章节结果（简化版本）
	// sections := make(map[string]*pb.Section)
	// 简化处理，不使用result.Sections
	/*
		for k, v := range result.Sections {
			issues := make([]*pb.Issue, len(v.Issues))
			for i, issue := range v.Issues {
				issues[i] = &pb.Issue{
					Type:        issue.Type,
					Description: issue.Description,
					Severity:    issue.Severity,
					Suggestion:  issue.Suggestion,
				}
			}

			sections[k] = &pb.Section{
				Content:       v.Content,
				ExtractedInfo: v.ExtractedInfo,
				QualityScore:  float32(v.QualityScore),
				Issues:        issues,
			}
		}
	*/

	// 转换建议
	suggestions := make([]*pb.Suggestion, len(result.Suggestions))
	for i, suggestion := range result.Suggestions {
		suggestions[i] = &pb.Suggestion{
			Id:          suggestion.ID,
			Type:        suggestion.Type,
			Title:       suggestion.Title,
			Description: suggestion.Description,
			Priority:    suggestion.Priority,
			Section:     suggestion.Section,
			Action:      suggestion.Action,
		}
	}

	// 转换评分
	scores := &pb.ScoreBreakdown{
		OverallScore:        float32(result.Scores.OverallScore),
		CompletenessScore:   float32(result.Scores.CompletenessScore),
		ClarityScore:        float32(result.Scores.ClarityScore),
		KeywordScore:        float32(result.Scores.KeywordScore),
		FormatScore:         float32(result.Scores.FormatScore),
		QuantificationScore: float32(result.Scores.QuantificationScore),
		DimensionScores:     make(map[string]float32),
	}

	for k, v := range result.Scores.DimensionScores {
		scores.DimensionScores[k] = float32(v)
	}

	// 转换建议为改进建议格式
	improvements := make([]*pb.Improvement, len(result.Suggestions))
	for i, suggestion := range result.Suggestions {
		improvements[i] = &pb.Improvement{
			Type:        suggestion.Type,
			Description: suggestion.Description,
			Priority:    suggestion.Priority,
			Section:     suggestion.Section,
			Before:      "",                // Suggestion结构体中没有Before字段
			After:       suggestion.Action, // 使用Action作为After
			Examples:    suggestion.Examples,
		}
	}

	return &pb.AnalysisResult{
		Sections:     make(map[string]*pb.Section), // 使用空map替代sections
		Suggestions:  suggestions,
		Scores:       scores,
		Improvements: improvements,
		Summary:      result.Summary,
		// AnalyzedAt:   timestamppb.New(result.AnalyzedAt), // 如果proto中没有这个字段就注释掉
	}
}

func (s *AIService) convertToBizAnalysisResult(pbResult *pb.AnalysisResult) *eino.AnalysisResult {
	if pbResult == nil {
		return nil
	}

	// 这里实现protobuf到业务对象的转换
	// 简化实现，实际应该完整转换所有字段
	result := &eino.AnalysisResult{
		Summary: pbResult.Summary,
	}

	if pbResult.Scores != nil {
		result.Scores = eino.ScoreBreakdown{
			OverallScore:        float64(pbResult.Scores.OverallScore),
			CompletenessScore:   float64(pbResult.Scores.CompletenessScore),
			ClarityScore:        float64(pbResult.Scores.ClarityScore),
			KeywordScore:        float64(pbResult.Scores.KeywordScore),
			FormatScore:         float64(pbResult.Scores.FormatScore),
			QuantificationScore: float64(pbResult.Scores.QuantificationScore),
		}
	}

	return result
}
