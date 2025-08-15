package server

import (
	"context"
	"net/http"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	khttp "github.com/go-kratos/kratos/v2/transport/http"

	v1 "github.com/liyubo06/resumeOptim_claude/api/user/v1"
	"github.com/liyubo06/resumeOptim_claude/backend/services/user-service/internal/conf"
	"github.com/liyubo06/resumeOptim_claude/backend/services/user-service/internal/service"
)

// CORS 中间件
func CORS() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if httpReq, ok := req.(*http.Request); ok {
				// 检查是否是预检请求
				if httpReq.Method == "OPTIONS" {
					return nil, nil
				}
			}
			return handler(ctx, req)
		}
	}
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, userService *service.UserService, logger log.Logger) *khttp.Server {
	var opts = []khttp.ServerOption{
		khttp.Middleware(
			recovery.Recovery(),
			tracing.Server(),
			CORS(),
		),
		khttp.Filter(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// 设置CORS头
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				w.Header().Set("Access-Control-Allow-Credentials", "true")

				// 处理预检请求
				if r.Method == "OPTIONS" {
					w.WriteHeader(http.StatusOK)
					return
				}

				next.ServeHTTP(w, r)
			})
		}),
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
	v1.RegisterUserServiceHTTPServer(srv, userService)

	// 添加健康检查端点
	srv.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"user-service"}`))
	})

	return srv
}
