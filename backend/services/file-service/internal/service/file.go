package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"

	v1 "github.com/lyb88999/resume_helper/backend/services/file-service/api/file/v1"
	"github.com/lyb88999/resume_helper/backend/services/file-service/internal/biz"
)

// FileService 文件服务实现
type FileService struct {
	v1.UnimplementedFileServiceServer
	uc  *biz.FileUsecase
	log *log.Helper
}

// NewFileService 创建文件服务实例
func NewFileService(uc *biz.FileUsecase, logger log.Logger) *FileService {
	return &FileService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

// Upload 上传文件
func (s *FileService) Upload(ctx context.Context, req *v1.UploadRequest) (*v1.UploadReply, error) {
	s.log.WithContext(ctx).Infof("文件上传请求: filename=%s, size=%d", req.Filename, len(req.File))

	reply, err := s.uc.Upload(ctx, req)
	if err != nil {
		s.log.WithContext(ctx).Errorf("文件上传失败: %v", err)
		return nil, err
	}

	s.log.WithContext(ctx).Infof("文件上传成功: file_id=%s", reply.File.FileId)
	return reply, nil
}

// ListFiles 获取文件列表
func (s *FileService) ListFiles(ctx context.Context, req *v1.ListFilesRequest) (*v1.ListFilesReply, error) {
	s.log.WithContext(ctx).Infof("获取文件列表请求: user_id=%d, page=%d", req.UserId, req.Page)

	reply, err := s.uc.ListFiles(ctx, req)
	if err != nil {
		s.log.WithContext(ctx).Errorf("获取文件列表失败: %v", err)
		return nil, err
	}

	s.log.WithContext(ctx).Infof("获取文件列表成功: total=%d", reply.Total)
	return reply, nil
}

// GetFile 获取文件详情
func (s *FileService) GetFile(ctx context.Context, req *v1.GetFileRequest) (*v1.GetFileReply, error) {
	s.log.WithContext(ctx).Infof("获取文件详情请求: file_id=%s", req.FileId)

	reply, err := s.uc.GetFile(ctx, req)
	if err != nil {
		s.log.WithContext(ctx).Errorf("获取文件详情失败: %v", err)
		return nil, err
	}

	return reply, nil
}

// DownloadFile 下载文件
func (s *FileService) DownloadFile(ctx context.Context, req *v1.DownloadFileRequest) (*v1.DownloadFileReply, error) {
	s.log.WithContext(ctx).Infof("下载文件请求: file_id=%s", req.FileId)

	reply, err := s.uc.DownloadFile(ctx, req)
	if err != nil {
		s.log.WithContext(ctx).Errorf("下载文件失败: %v", err)
		return nil, err
	}

	s.log.WithContext(ctx).Infof("下载文件成功: filename=%s", reply.Filename)
	return reply, nil
}

// DeleteFile 删除文件
func (s *FileService) DeleteFile(ctx context.Context, req *v1.DeleteFileRequest) (*v1.DeleteFileReply, error) {
	s.log.WithContext(ctx).Infof("删除文件请求: file_id=%s", req.FileId)

	reply, err := s.uc.DeleteFile(ctx, req)
	if err != nil {
		s.log.WithContext(ctx).Errorf("删除文件失败: %v", err)
		return nil, err
	}

	s.log.WithContext(ctx).Infof("删除文件成功: file_id=%s", req.FileId)
	return reply, nil
}
