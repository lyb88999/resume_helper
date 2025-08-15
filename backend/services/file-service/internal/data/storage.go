package data

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/lyb88999/resume_helper/backend/services/file-service/internal/biz"
	"github.com/lyb88999/resume_helper/backend/shared/proto/conf"
)

type localStorage struct {
	basePath string
	log      *log.Helper
}

// NewStorageRepo 创建存储仓库
func NewStorageRepo(config *conf.Storage, logger log.Logger) biz.StorageRepo {
	helper := log.NewHelper(logger)

	switch config.Type {
	case "local":
		return &localStorage{
			basePath: config.Local.Path,
			log:      helper,
		}
	default:
		// 默认使用本地存储
		return &localStorage{
			basePath: "/tmp/uploads",
			log:      helper,
		}
	}
}

func (s *localStorage) Upload(ctx context.Context, filename string, content io.Reader) (string, error) {
	// 确保目录存在
	if err := os.MkdirAll(s.basePath, 0755); err != nil {
		s.log.WithContext(ctx).Errorf("failed to create directory: %v", err)
		return "", err
	}

	filePath := filepath.Join(s.basePath, filename)

	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		s.log.WithContext(ctx).Errorf("failed to create file: %v", err)
		return "", err
	}
	defer file.Close()

	// 写入内容
	if _, err := io.Copy(file, content); err != nil {
		s.log.WithContext(ctx).Errorf("failed to write file: %v", err)
		os.Remove(filePath) // 清理失败的文件
		return "", err
	}

	// 返回文件URL（这里简化为文件路径）
	return fmt.Sprintf("/files/%s", filename), nil
}

func (s *localStorage) Download(ctx context.Context, filename string) (io.ReadCloser, error) {
	filePath := filepath.Join(s.basePath, filename)

	file, err := os.Open(filePath)
	if err != nil {
		s.log.WithContext(ctx).Errorf("failed to open file: %v", err)
		return nil, err
	}

	return file, nil
}

func (s *localStorage) Delete(ctx context.Context, filename string) error {
	filePath := filepath.Join(s.basePath, filename)

	if err := os.Remove(filePath); err != nil {
		s.log.WithContext(ctx).Errorf("failed to delete file: %v", err)
		return err
	}

	return nil
}

func (s *localStorage) GetURL(filename string) string {
	return fmt.Sprintf("/files/%s", filename)
}
