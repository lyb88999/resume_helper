package data

import (
	"github.com/lyb88999/resume_helper/backend/services/parser-service/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewParseTaskRepo)

// Data .
type Data struct {
	db *gorm.DB
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	log := log.NewHelper(logger)

	// 连接数据库
	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// 自动迁移数据表
	if err := db.AutoMigrate(&ParseTaskModel{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	cleanup := func() {
		log.Info("closing the data resources")
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
	}

	return &Data{
		db: db,
	}, cleanup, nil
}
