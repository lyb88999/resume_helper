package data

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"

	"github.com/lyb88999/resume_helper/backend/services/ai-service/internal/eino"
)

// aiRepo AI数据仓库实现
type aiRepo struct {
	data *Data
	log  *log.Helper
}

// AnalysisResultModel 分析结果数据模型
type AnalysisResultModel struct {
	ID             string    `gorm:"primaryKey;size:64" json:"id"`
	ResumeID       string    `gorm:"index;size:64;not null" json:"resume_id"`
	TargetPosition string    `gorm:"size:100" json:"target_position"`
	ResultData     string    `gorm:"type:longtext" json:"result_data"` // JSON格式存储
	OverallScore   float64   `gorm:"index" json:"overall_score"`
	Status         string    `gorm:"size:20;default:completed" json:"status"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// ChatSessionModel 聊天会话数据模型
type ChatSessionModel struct {
	ID         string    `gorm:"primaryKey;size:64" json:"id"`
	SessionID  string    `gorm:"uniqueIndex;size:64;not null" json:"session_id"`
	Messages   string    `gorm:"type:longtext" json:"messages"` // JSON格式存储消息列表
	Context    string    `gorm:"type:text" json:"context"`      // JSON格式存储上下文
	LastActive time.Time `gorm:"index" json:"last_active"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 设置表名
func (AnalysisResultModel) TableName() string {
	return "ai_analysis_results"
}

func (ChatSessionModel) TableName() string {
	return "ai_chat_sessions"
}

// SaveAnalysisResult 保存分析结果
func (r *aiRepo) SaveAnalysisResult(ctx context.Context, result *eino.AnalysisResult) error {
	r.log.WithContext(ctx).Infof("保存分析结果: %s", result.ID)

	// 序列化结果数据
	resultData, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("序列化分析结果失败: %w", err)
	}

	model := &AnalysisResultModel{
		ID:             result.ID,
		ResumeID:       result.ResumeID,
		TargetPosition: result.TargetPosition,
		ResultData:     string(resultData),
		OverallScore:   result.Scores.OverallScore,
		Status:         "completed",
	}

	// 使用GORM保存到数据库
	if err := r.data.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("保存分析结果到数据库失败: %w", err)
	}

	// 同时缓存到Redis（1小时过期）
	cacheKey := fmt.Sprintf("analysis_result:%s", result.ID)
	if err := r.data.rdb.Set(ctx, cacheKey, string(resultData), time.Hour).Err(); err != nil {
		r.log.WithContext(ctx).Warnf("缓存分析结果失败: %v", err)
		// 缓存失败不影响主流程
	}

	return nil
}

// GetAnalysisResult 获取分析结果
func (r *aiRepo) GetAnalysisResult(ctx context.Context, id string) (*eino.AnalysisResult, error) {
	r.log.WithContext(ctx).Infof("获取分析结果: %s", id)

	// 首先尝试从Redis缓存获取
	cacheKey := fmt.Sprintf("analysis_result:%s", id)
	cachedData, err := r.data.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var result eino.AnalysisResult
		if err := json.Unmarshal([]byte(cachedData), &result); err == nil {
			r.log.WithContext(ctx).Info("从缓存获取分析结果成功")
			return &result, nil
		}
		r.log.WithContext(ctx).Warnf("缓存数据反序列化失败: %v", err)
	}

	// 从数据库获取
	var model AnalysisResultModel
	if err := r.data.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("分析结果不存在: %s", id)
		}
		return nil, fmt.Errorf("查询分析结果失败: %w", err)
	}

	// 反序列化结果数据
	var result eino.AnalysisResult
	if err := json.Unmarshal([]byte(model.ResultData), &result); err != nil {
		return nil, fmt.Errorf("反序列化分析结果失败: %w", err)
	}

	// 更新缓存
	if err := r.data.rdb.Set(ctx, cacheKey, model.ResultData, time.Hour).Err(); err != nil {
		r.log.WithContext(ctx).Warnf("更新缓存失败: %v", err)
	}

	return &result, nil
}

// SaveChatSession 保存聊天会话
func (r *aiRepo) SaveChatSession(ctx context.Context, session *eino.ChatContext) error {
	r.log.WithContext(ctx).Infof("保存聊天会话: %s", session.SessionID)

	// 序列化消息和上下文
	messagesData, err := json.Marshal(session.Messages)
	if err != nil {
		return fmt.Errorf("序列化消息失败: %w", err)
	}

	contextData, err := json.Marshal(map[string]interface{}{
		"resume_data":     session.ResumeData,
		"target_position": session.TargetPosition,
		"knowledge":       session.Knowledge,
	})
	if err != nil {
		return fmt.Errorf("序列化上下文失败: %w", err)
	}

	model := &ChatSessionModel{
		SessionID:  session.SessionID,
		Messages:   string(messagesData),
		Context:    string(contextData),
		LastActive: time.Now(),
	}

	// 使用UPSERT操作
	if err := r.data.db.WithContext(ctx).
		Where("session_id = ?", session.SessionID).
		Assign(map[string]interface{}{
			"messages":    model.Messages,
			"context":     model.Context,
			"last_active": model.LastActive,
		}).
		FirstOrCreate(model).Error; err != nil {
		return fmt.Errorf("保存聊天会话失败: %w", err)
	}

	// 缓存会话（30分钟过期）
	cacheKey := fmt.Sprintf("chat_session:%s", session.SessionID)
	sessionData, _ := json.Marshal(session)
	if err := r.data.rdb.Set(ctx, cacheKey, string(sessionData), 30*time.Minute).Err(); err != nil {
		r.log.WithContext(ctx).Warnf("缓存会话失败: %v", err)
	}

	return nil
}

// GetChatSession 获取聊天会话
func (r *aiRepo) GetChatSession(ctx context.Context, sessionID string) (*eino.ChatContext, error) {
	r.log.WithContext(ctx).Infof("获取聊天会话: %s", sessionID)

	// 首先尝试从Redis缓存获取
	cacheKey := fmt.Sprintf("chat_session:%s", sessionID)
	cachedData, err := r.data.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var session eino.ChatContext
		if err := json.Unmarshal([]byte(cachedData), &session); err == nil {
			r.log.WithContext(ctx).Info("从缓存获取会话成功")
			return &session, nil
		}
		r.log.WithContext(ctx).Warnf("缓存会话数据反序列化失败: %v", err)
	}

	// 从数据库获取
	var model ChatSessionModel
	if err := r.data.db.WithContext(ctx).Where("session_id = ?", sessionID).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("会话不存在: %s", sessionID)
		}
		return nil, fmt.Errorf("查询会话失败: %w", err)
	}

	// 反序列化数据
	var messages []eino.Message
	if err := json.Unmarshal([]byte(model.Messages), &messages); err != nil {
		return nil, fmt.Errorf("反序列化消息失败: %w", err)
	}

	var contextData map[string]interface{}
	if err := json.Unmarshal([]byte(model.Context), &contextData); err != nil {
		return nil, fmt.Errorf("反序列化上下文失败: %w", err)
	}

	session := &eino.ChatContext{
		SessionID: sessionID,
		Messages:  messages,
	}

	// 恢复上下文数据
	if resumeData, ok := contextData["resume_data"]; ok && resumeData != nil {
		if resumeBytes, err := json.Marshal(resumeData); err == nil {
			var resume eino.ResumeData
			if err := json.Unmarshal(resumeBytes, &resume); err == nil {
				session.ResumeData = &resume
			}
		}
	}

	if targetPosition, ok := contextData["target_position"].(string); ok {
		session.TargetPosition = targetPosition
	}

	if knowledge, ok := contextData["knowledge"]; ok && knowledge != nil {
		if knowledgeBytes, err := json.Marshal(knowledge); err == nil {
			var docs []eino.Document
			if err := json.Unmarshal(knowledgeBytes, &docs); err == nil {
				session.Knowledge = docs
			}
		}
	}

	// 更新缓存
	sessionData, _ := json.Marshal(session)
	if err := r.data.rdb.Set(ctx, cacheKey, string(sessionData), 30*time.Minute).Err(); err != nil {
		r.log.WithContext(ctx).Warnf("更新会话缓存失败: %v", err)
	}

	return session, nil
}

// CleanupExpiredSessions 清理过期会话（可以定时调用）
func (r *aiRepo) CleanupExpiredSessions(ctx context.Context, expiredBefore time.Time) error {
	result := r.data.db.WithContext(ctx).
		Where("last_active < ?", expiredBefore).
		Delete(&ChatSessionModel{})

	if result.Error != nil {
		return fmt.Errorf("清理过期会话失败: %w", result.Error)
	}

	r.log.WithContext(ctx).Infof("清理了 %d 个过期会话", result.RowsAffected)
	return nil
}

// GetRecentAnalysisResults 获取最近的分析结果
func (r *aiRepo) GetRecentAnalysisResults(ctx context.Context, limit int) ([]*eino.AnalysisResult, error) {
	var models []AnalysisResultModel
	if err := r.data.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("查询最近分析结果失败: %w", err)
	}

	results := make([]*eino.AnalysisResult, len(models))
	for i, model := range models {
		var result eino.AnalysisResult
		if err := json.Unmarshal([]byte(model.ResultData), &result); err != nil {
			r.log.WithContext(ctx).Warnf("反序列化分析结果失败: %v", err)
			continue
		}
		results[i] = &result
	}

	return results, nil
}
