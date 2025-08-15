package server

import (
	"encoding/json"
	"net/http"
	"time"

	v1 "github.com/lyb88999/resume_helper/backend/services/parser-service/api/parser/v1"
	"github.com/lyb88999/resume_helper/backend/services/parser-service/internal/service"
	"github.com/lyb88999/resume_helper/backend/shared/proto/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, parserService *service.ParserService, logger log.Logger) *khttp.Server {
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

	// 注册解析服务
	v1.RegisterParserServiceHTTPServer(srv, parserService)

	// 添加健康检查端点
	srv.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"status":    "ok",
			"service":   "parser-service",
			"timestamp": time.Now().Format(time.RFC3339),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	return srv
}
