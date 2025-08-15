package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Milvus   MilvusConfig   `mapstructure:"milvus"`
	Eino     EinoConfig     `mapstructure:"eino"`
	Log      LogConfig      `mapstructure:"log"`
	Tracing  TracingConfig  `mapstructure:"tracing"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	HTTP struct {
		Port    int    `mapstructure:"port"`
		Timeout int    `mapstructure:"timeout"`
	} `mapstructure:"http"`
	GRPC struct {
		Port    int    `mapstructure:"port"`
		Timeout int    `mapstructure:"timeout"`
	} `mapstructure:"grpc"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver       string `mapstructure:"driver"`
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	Database     string `mapstructure:"database"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// MilvusConfig Milvus向量数据库配置
type MilvusConfig struct {
	Host       string `mapstructure:"host"`
	Port       int    `mapstructure:"port"`
	Username   string `mapstructure:"username"`
	Password   string `mapstructure:"password"`
	Database   string `mapstructure:"database"`
	Collection string `mapstructure:"collection"`
}

// EinoConfig Eino AI框架配置
type EinoConfig struct {
	ModelProvider string            `mapstructure:"model_provider"`
	OpenAI        OpenAIConfig      `mapstructure:"openai"`
	Claude        ClaudeConfig      `mapstructure:"claude"`
	Qwen          QwenConfig        `mapstructure:"qwen"`
	Embeddings    EmbeddingsConfig  `mapstructure:"embeddings"`
	Workflows     WorkflowsConfig   `mapstructure:"workflows"`
}

// OpenAIConfig OpenAI配置
type OpenAIConfig struct {
	APIKey      string  `mapstructure:"api_key"`
	BaseURL     string  `mapstructure:"base_url"`
	Model       string  `mapstructure:"model"`
	Temperature float64 `mapstructure:"temperature"`
	MaxTokens   int     `mapstructure:"max_tokens"`
}

// ClaudeConfig Claude配置
type ClaudeConfig struct {
	APIKey      string  `mapstructure:"api_key"`
	BaseURL     string  `mapstructure:"base_url"`
	Model       string  `mapstructure:"model"`
	Temperature float64 `mapstructure:"temperature"`
	MaxTokens   int     `mapstructure:"max_tokens"`
}

// QwenConfig 通义千问配置
type QwenConfig struct {
	APIKey      string  `mapstructure:"api_key"`
	BaseURL     string  `mapstructure:"base_url"`
	Model       string  `mapstructure:"model"`
	Temperature float64 `mapstructure:"temperature"`
	MaxTokens   int     `mapstructure:"max_tokens"`
}

// EmbeddingsConfig 向量化配置
type EmbeddingsConfig struct {
	Provider string `mapstructure:"provider"`
	Model    string `mapstructure:"model"`
	APIKey   string `mapstructure:"api_key"`
}

// WorkflowsConfig 工作流配置
type WorkflowsConfig struct {
	ResumeParsingTimeout  int `mapstructure:"resume_parsing_timeout"`
	AnalysisTimeout       int `mapstructure:"analysis_timeout"`
	KnowledgeRetrieval    int `mapstructure:"knowledge_retrieval"`
	MaxConcurrency        int `mapstructure:"max_concurrency"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Encoding   string `mapstructure:"encoding"`
	OutputPath string `mapstructure:"output_path"`
}

// TracingConfig 链路追踪配置
type TracingConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	ServiceName string `mapstructure:"service_name"`
	Endpoint    string `mapstructure:"endpoint"`
}

// LoadConfig 加载配置
func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")
	
	// 设置环境变量前缀
	viper.SetEnvPrefix("RESUME_OPTIM")
	viper.AutomaticEnv()
	
	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	
	return &config, nil
}

// GetLogger 获取日志记录器
func GetLogger(config LogConfig) (*zap.Logger, error) {
	var logger *zap.Logger
	var err error
	
	if config.Level == "production" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
	
	if err != nil {
		return nil, err
	}
	
	return logger, nil
}
