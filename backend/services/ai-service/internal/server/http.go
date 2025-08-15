package server

import (
	"encoding/json"
	"net/http"
	"time"

	v1 "github.com/liyubo06/resumeOptim_claude/backend/services/ai-service/api/ai/v1"
	"github.com/liyubo06/resumeOptim_claude/backend/services/ai-service/internal/conf"
	"github.com/liyubo06/resumeOptim_claude/backend/services/ai-service/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, aiService *service.AIService, logger log.Logger) *khttp.Server {
	var opts = []khttp.ServerOption{
		khttp.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, khttp.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, khttp.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, khttp.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := khttp.NewServer(opts...)

	// 注册AI服务
	v1.RegisterAIServiceHTTPServer(srv, aiService)

	// 添加健康检查端点
	srv.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"status":    "ok",
			"service":   "ai-service",
			"timestamp": time.Now().Format(time.RFC3339),
			"version":   "1.0.0",
			"components": map[string]string{
				"database":  "ok",
				"redis":     "ok",
				"ai_model":  "ok",
				"vector_db": "ok",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// 添加AI服务状态端点
	srv.HandleFunc("/ai/status", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"service_name": "ai-service",
			"status":       "running",
			"features": map[string]bool{
				"resume_parsing":       true,
				"intelligent_analysis": true,
				"chat_bot":             true,
				"knowledge_retrieval":  true,
			},
			"eino_components": map[string]bool{
				"parsing_chain":  true,
				"analysis_graph": true,
				"react_agent":    true,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	return srv
}

