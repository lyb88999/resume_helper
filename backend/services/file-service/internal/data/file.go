package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"

	"github.com/lyb88999/resume_helper/backend/services/file-service/internal/biz"
)

// FileModel 文件数据模型
type FileModel struct {
	ID           uint   `gorm:"primarykey"`
	FileID       string `gorm:"uniqueIndex;size:100;not null"`
	Filename     string `gorm:"size:255;not null"`
	OriginalName string `gorm:"size:255;not null"`
	Size         int64  `gorm:"not null"`
	MimeType     string `gorm:"size:100"`
	URL          string `gorm:"size:500"`
	Status       string `gorm:"size:50;default:'uploaded'"`
	UserID       int64  `gorm:"not null;index"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (FileModel) TableName() string {
	return "files"
}

type fileRepo struct {
	data *Data
	log  *log.Helper
}

// NewFileRepo .
func NewFileRepo(data *Data, logger log.Logger) biz.FileRepo {
	return &fileRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *fileRepo) Save(ctx context.Context, file *biz.File) (*biz.File, error) {
	model := &FileModel{
		FileID:       file.FileID,
		Filename:     file.Filename,
		OriginalName: file.OriginalName,
		Size:         file.Size,
		MimeType:     file.MimeType,
		URL:          file.URL,
		Status:       file.Status,
		UserID:       file.UserID,
	}

	if err := r.data.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}

	file.ID = int64(model.ID)
	file.CreatedAt = model.CreatedAt
	file.UpdatedAt = model.UpdatedAt
	return file, nil
}

func (r *fileRepo) Update(ctx context.Context, file *biz.File) (*biz.File, error) {
	model := &FileModel{
		ID:           uint(file.ID),
		FileID:       file.FileID,
		Filename:     file.Filename,
		OriginalName: file.OriginalName,
		Size:         file.Size,
		MimeType:     file.MimeType,
		URL:          file.URL,
		Status:       file.Status,
		UserID:       file.UserID,
	}

	if err := r.data.db.WithContext(ctx).Save(model).Error; err != nil {
		return nil, err
	}

	file.UpdatedAt = model.UpdatedAt
	return file, nil
}

func (r *fileRepo) FindByID(ctx context.Context, fileID string, userID int64) (*biz.File, error) {
	var model FileModel
	if err := r.data.db.WithContext(ctx).Where("file_id = ? AND user_id = ?", fileID, userID).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, biz.ErrFileNotFound
		}
		return nil, err
	}

	return &biz.File{
		ID:           int64(model.ID),
		FileID:       model.FileID,
		Filename:     model.Filename,
		OriginalName: model.OriginalName,
		Size:         model.Size,
		MimeType:     model.MimeType,
		URL:          model.URL,
		Status:       model.Status,
		UserID:       model.UserID,
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
	}, nil
}

func (r *fileRepo) List(ctx context.Context, req *biz.ListFilesRequest) ([]*biz.File, int64, error) {
	var models []FileModel
	var total int64

	query := r.data.db.WithContext(ctx).Model(&FileModel{}).Where("user_id = ?", req.UserID)

	// 添加筛选条件
	if req.Type != "" {
		query = query.Where("mime_type LIKE ?", "%"+req.Type+"%")
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.CreatedAfter != "" {
		query = query.Where("created_at >= ?", req.CreatedAfter)
	}
	if req.CreatedBefore != "" {
		query = query.Where("created_at <= ?", req.CreatedBefore)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PerPage
	if err := query.Offset(offset).Limit(req.PerPage).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, 0, err
	}

	files := make([]*biz.File, len(models))
	for i, model := range models {
		files[i] = &biz.File{
			ID:           int64(model.ID),
			FileID:       model.FileID,
			Filename:     model.Filename,
			OriginalName: model.OriginalName,
			Size:         model.Size,
			MimeType:     model.MimeType,
			URL:          model.URL,
			Status:       model.Status,
			UserID:       model.UserID,
			CreatedAt:    model.CreatedAt,
			UpdatedAt:    model.UpdatedAt,
		}
	}

	return files, total, nil
}

func (r *fileRepo) Delete(ctx context.Context, fileID string, userID int64) error {
	result := r.data.db.WithContext(ctx).Where("file_id = ? AND user_id = ?", fileID, userID).Delete(&FileModel{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return biz.ErrFileNotFound
	}
	return nil
}
