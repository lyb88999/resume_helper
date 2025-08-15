package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/liyubo06/resumeOptim_claude/backend/services/ai-service/internal/biz"
	"github.com/liyubo06/resumeOptim_claude/backend/services/ai-service/internal/conf"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewAIRepo)

// Data represents the data layer.
type Data struct {
	db  *gorm.DB
	rdb *redis.Client
	log *log.Helper
}

// NewData creates a new data layer.
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	helper := log.NewHelper(logger)

	// 初始化数据库连接
	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	// 自动迁移数据库表
	if err := db.AutoMigrate(&AnalysisResultModel{}, &ChatSessionModel{}); err != nil {
		return nil, nil, err
	}

	// 初始化Redis连接
	rdb := redis.NewClient(&redis.Options{
		Addr:         c.Redis.Addr,
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
	})

	// 测试Redis连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		helper.Warnf("Redis连接失败: %v", err)
	}

	d := &Data{
		db:  db,
		rdb: rdb,
		log: helper,
	}

	cleanup := func() {
		helper.Info("关闭数据层连接")
		if sqlDB, err := db.DB(); err == nil {
			_ = sqlDB.Close()
		}
		_ = rdb.Close()
	}

	return d, cleanup, nil
}

// NewAIRepo creates a new AI repository.
func NewAIRepo(data *Data, logger log.Logger) biz.AIRepo {
	return &aiRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
