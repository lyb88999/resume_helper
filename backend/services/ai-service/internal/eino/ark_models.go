package eino

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ARKChatModel ARK聊天模型实现
type ARKChatModel struct {
	APIKey      string
	BaseURL     string
	Model       string
	MaxTokens   int
	Temperature float32
	Timeout     time.Duration
}

// ARKRequest ARK API请求结构
type ARKRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float32   `json:"temperature,omitempty"`
	Stream      bool      `json:"stream"`
}

// ARKResponse ARK API响应结构
type ARKResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Generate 生成回复
func (a *ARKChatModel) Generate(ctx context.Context, messages []Message, options ...GenerateOption) (*GenerateResponse, error) {
	// 应用选项
	opts := &GenerateOptions{
		MaxTokens:   a.MaxTokens,
		Temperature: float64(a.Temperature),
	}
	for _, opt := range options {
		opt(opts)
	}

	// 构建请求
	reqBody := ARKRequest{
		Model:       a.Model,
		Messages:    messages,
		MaxTokens:   opts.MaxTokens,
		Temperature: float32(opts.Temperature),
		Stream:      false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "POST", a.BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+a.APIKey)

	// 设置超时
	client := &http.Client{
		Timeout: a.Timeout,
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求失败，状态码: %d", resp.StatusCode)
	}

	// 解析响应
	var arkResp ARKResponse
	if err := json.NewDecoder(resp.Body).Decode(&arkResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	// 转换为标准响应
	return &GenerateResponse{
		Choices: arkResp.Choices,
		Usage:   arkResp.Usage,
	}, nil
}

// ARKEmbeddingModel ARK嵌入模型实现
type ARKEmbeddingModel struct {
	APIKey  string
	BaseURL string
	Model   string
	Timeout time.Duration
}

// ARKEmbeddingRequest ARK嵌入请求结构
type ARKEmbeddingRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

// ARKEmbeddingResponse ARK嵌入响应结构
type ARKEmbeddingResponse struct {
	Object string             `json:"object"`
	Data   []ARKEmbeddingData `json:"data"`
	Model  string             `json:"model"`
	Usage  ARKEmbeddingUsage  `json:"usage"`
}

// ARKEmbeddingData 嵌入数据
type ARKEmbeddingData struct {
	Object    string    `json:"object"`
	Index     int       `json:"index"`
	Embedding []float64 `json:"embedding"`
}

// ARKEmbeddingUsage 嵌入使用统计
type ARKEmbeddingUsage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

// Embed 生成嵌入向量
func (a *ARKEmbeddingModel) Embed(ctx context.Context, texts []string) ([][]float64, error) {
	// 构建请求
	reqBody := ARKEmbeddingRequest{
		Model: a.Model,
		Input: texts,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "POST", a.BaseURL+"/embeddings", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+a.APIKey)

	// 设置超时
	client := &http.Client{
		Timeout: a.Timeout,
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求失败，状态码: %d", resp.StatusCode)
	}

	// 解析响应
	var arkResp ARKEmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&arkResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	// 提取嵌入向量
	embeddings := make([][]float64, len(arkResp.Data))
	for i, data := range arkResp.Data {
		embeddings[i] = data.Embedding
	}

	return embeddings, nil
}
