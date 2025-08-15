package biz

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"mime"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	v1 "github.com/liyubo06/resumeOptim_claude/backend/services/file-service/api/file/v1"
	"github.com/liyubo06/resumeOptim_claude/backend/services/file-service/internal/conf"
)

var (
	ErrFileNotFound     = errors.New("file not found")
	ErrInvalidFileType  = errors.New("invalid file type")
	ErrFileSizeExceeded = errors.New("file size exceeded")
	ErrPermissionDenied = errors.New("permission denied")
)

// File 文件业务模型
type File struct {
	ID           int64
	FileID       string
	Filename     string
	OriginalName string
	Size         int64
	MimeType     string
	URL          string
	Status       string
	UserID       int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// ListFilesRequest 文件列表请求
type ListFilesRequest struct {
	Page          int
	PerPage       int
	Type          string
	Status        string
	UserID        int64
	CreatedAfter  string
	CreatedBefore string
}

// FileRepo 文件仓库接口
type FileRepo interface {
	Save(ctx context.Context, file *File) (*File, error)
	Update(ctx context.Context, file *File) (*File, error)
	FindByID(ctx context.Context, fileID string, userID int64) (*File, error)
	List(ctx context.Context, req *ListFilesRequest) ([]*File, int64, error)
	Delete(ctx context.Context, fileID string, userID int64) error
}

// StorageRepo 存储仓库接口
type StorageRepo interface {
	Upload(ctx context.Context, filename string, content io.Reader) (string, error)
	Download(ctx context.Context, filename string) (io.ReadCloser, error)
	Delete(ctx context.Context, filename string) error
	GetURL(filename string) string
}

// FileUsecase 文件用例
type FileUsecase struct {
	repo    FileRepo
	storage StorageRepo
	config  *conf.Storage
	log     *log.Helper
}

// NewFileUsecase 创建文件用例
func NewFileUsecase(repo FileRepo, storage StorageRepo, config *conf.Storage, logger log.Logger) *FileUsecase {
	return &FileUsecase{
		repo:    repo,
		storage: storage,
		config:  config,
		log:     log.NewHelper(logger),
	}
}

// Upload 上传文件
func (uc *FileUsecase) Upload(ctx context.Context, req *v1.UploadRequest) (*v1.UploadReply, error) {
	// 验证文件大小
	if int64(len(req.File)) > uc.config.MaxFileSize {
		return nil, ErrFileSizeExceeded
	}

	// 检测文件类型
	mimeType := mime.TypeByExtension(filepath.Ext(req.Filename))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// 验证文件类型
	if !uc.isAllowedType(req.Filename) {
		return nil, ErrInvalidFileType
	}

	// 生成文件ID和存储文件名
	fileID := uuid.New().String()
	ext := filepath.Ext(req.Filename)
	storageFilename := fmt.Sprintf("%s%s", fileID, ext)

	// 上传到存储
	reader := strings.NewReader(string(req.File))
	storageURL, err := uc.storage.Upload(ctx, storageFilename, reader)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to upload file: %v", err)
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	// 保存文件信息
	file := &File{
		FileID:       fileID,
		Filename:     storageFilename,
		OriginalName: req.Filename,
		Size:         int64(len(req.File)),
		MimeType:     mimeType,
		URL:          storageURL,
		Status:       "uploaded",
		UserID:       req.UserId,
	}

	savedFile, err := uc.repo.Save(ctx, file)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to save file info: %v", err)
		// 清理已上传的文件
		uc.storage.Delete(ctx, storageFilename)
		return nil, fmt.Errorf("failed to save file info: %w", err)
	}

	return &v1.UploadReply{
		File: uc.toProtoFileInfo(savedFile),
	}, nil
}

// ListFiles 获取文件列表
func (uc *FileUsecase) ListFiles(ctx context.Context, req *v1.ListFilesRequest) (*v1.ListFilesReply, error) {
	bizReq := &ListFilesRequest{
		Page:          int(req.Page),
		PerPage:       int(req.PerPage),
		Type:          req.Type,
		Status:        req.Status,
		UserID:        req.UserId,
		CreatedAfter:  req.CreatedAfter,
		CreatedBefore: req.CreatedBefore,
	}

	files, total, err := uc.repo.List(ctx, bizReq)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to list files: %v", err)
		return nil, err
	}

	protoFiles := make([]*v1.FileInfo, len(files))
	for i, file := range files {
		protoFiles[i] = uc.toProtoFileInfo(file)
	}

	totalPages := int32(math.Ceil(float64(total) / float64(req.PerPage)))

	return &v1.ListFilesReply{
		Files:      protoFiles,
		Total:      int32(total),
		Page:       req.Page,
		PerPage:    req.PerPage,
		TotalPages: totalPages,
	}, nil
}

// GetFile 获取文件详情
func (uc *FileUsecase) GetFile(ctx context.Context, req *v1.GetFileRequest) (*v1.GetFileReply, error) {
	file, err := uc.repo.FindByID(ctx, req.FileId, req.UserId)
	if err != nil {
		if err == ErrFileNotFound {
			return nil, ErrFileNotFound
		}
		uc.log.WithContext(ctx).Errorf("failed to get file: %v", err)
		return nil, err
	}

	return &v1.GetFileReply{
		File: uc.toProtoFileInfo(file),
	}, nil
}

// DownloadFile 下载文件
func (uc *FileUsecase) DownloadFile(ctx context.Context, req *v1.DownloadFileRequest) (*v1.DownloadFileReply, error) {
	// 获取文件信息
	file, err := uc.repo.FindByID(ctx, req.FileId, req.UserId)
	if err != nil {
		if err == ErrFileNotFound {
			return nil, ErrFileNotFound
		}
		uc.log.WithContext(ctx).Errorf("failed to get file: %v", err)
		return nil, err
	}

	// 从存储下载文件
	reader, err := uc.storage.Download(ctx, file.Filename)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to download file: %v", err)
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	defer reader.Close()

	// 读取文件内容
	content, err := io.ReadAll(reader)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to read file content: %v", err)
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	return &v1.DownloadFileReply{
		Content:  content,
		Filename: file.OriginalName,
		MimeType: file.MimeType,
	}, nil
}

// DeleteFile 删除文件
func (uc *FileUsecase) DeleteFile(ctx context.Context, req *v1.DeleteFileRequest) (*v1.DeleteFileReply, error) {
	// 获取文件信息
	file, err := uc.repo.FindByID(ctx, req.FileId, req.UserId)
	if err != nil {
		if err == ErrFileNotFound {
			return nil, ErrFileNotFound
		}
		uc.log.WithContext(ctx).Errorf("failed to get file: %v", err)
		return nil, err
	}

	// 从存储删除文件
	if err := uc.storage.Delete(ctx, file.Filename); err != nil {
		uc.log.WithContext(ctx).Errorf("failed to delete file from storage: %v", err)
	}

	// 从数据库删除记录
	if err := uc.repo.Delete(ctx, req.FileId, req.UserId); err != nil {
		uc.log.WithContext(ctx).Errorf("failed to delete file record: %v", err)
		return nil, err
	}

	return &v1.DeleteFileReply{
		Success: true,
		Message: "file deleted successfully",
	}, nil
}

// isAllowedType 检查文件类型是否允许
func (uc *FileUsecase) isAllowedType(filename string) bool {
	if len(uc.config.AllowedTypes) == 0 {
		return true // 如果没有配置限制，则允许所有类型
	}

	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(filename), "."))
	for _, allowedType := range uc.config.AllowedTypes {
		if strings.ToLower(allowedType) == ext {
			return true
		}
	}
	return false
}

// toProtoFileInfo 转换为proto文件信息
func (uc *FileUsecase) toProtoFileInfo(file *File) *v1.FileInfo {
	return &v1.FileInfo{
		FileId:       file.FileID,
		Filename:     file.Filename,
		OriginalName: file.OriginalName,
		Size:         file.Size,
		MimeType:     file.MimeType,
		Url:          file.URL,
		Status:       file.Status,
		UserId:       file.UserID,
		CreatedAt:    timestamppb.New(file.CreatedAt),
		UpdatedAt:    timestamppb.New(file.UpdatedAt),
	}
}
