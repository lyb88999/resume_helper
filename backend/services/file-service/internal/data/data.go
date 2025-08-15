package data

import (
	"time"

	"github.com/liyubo06/resumeOptim_claude/backend/services/file-service/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewFileRepo, NewStorageRepo)

// Data .
type Data struct {
	db *gorm.DB
}

// NewData .
func NewData(c *conf.Data, l log.Logger) (*Data, func(), error) {
	helper := log.NewHelper(l)

	// 连接数据库
	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		helper.Errorf("failed to connect to database: %v", err)
		return nil, nil, err
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		helper.Errorf("failed to get database instance: %v", err)
		return nil, nil, err
	}

	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动迁移
	if err := db.AutoMigrate(&FileModel{}); err != nil {
		helper.Errorf("failed to migrate database: %v", err)
		return nil, nil, err
	}

	data := &Data{
		db: db,
	}

	cleanup := func() {
		helper.Info("closing the data resources")
		if sqlDB != nil {
			sqlDB.Close()
		}
	}

	return data, cleanup, nil
}
