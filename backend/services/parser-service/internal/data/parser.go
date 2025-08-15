package data

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/lyb88999/resume_helper/backend/services/parser-service/internal/biz"
	"gorm.io/gorm"
)

// ParseTaskModel 解析任务数据模型
type ParseTaskModel struct {
	ID          string     `gorm:"primaryKey;size:32" json:"id"`
	ResumeID    string     `gorm:"size:32;not null;index" json:"resume_id"`
	UserID      string     `gorm:"size:32;not null;index" json:"user_id"`
	FilePath    string     `gorm:"size:500;not null" json:"file_path"`
	FileType    string     `gorm:"size:10;not null" json:"file_type"`
	Status      string     `gorm:"size:20;not null;default:'pending'" json:"status"` // pending, processing, completed, failed
	Progress    int        `gorm:"default:0" json:"progress"`                        // 0-100
	Result      string     `gorm:"type:longtext" json:"result"`                      // JSON格式的解析结果
	ErrorMsg    string     `gorm:"size:1000" json:"error_msg"`
	Options     string     `gorm:"type:text" json:"options"` // JSON格式的解析选项
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

func (ParseTaskModel) TableName() string {
	return "parse_tasks"
}

// parseTaskRepo 解析任务仓库实现
type parseTaskRepo struct {
	data *Data
	log  *log.Helper
}

// NewParseTaskRepo 创建解析任务仓库
func NewParseTaskRepo(data *Data, logger log.Logger) biz.ParseTaskRepo {
	return &parseTaskRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *parseTaskRepo) CreateTask(ctx context.Context, task *biz.ParseTask) (*biz.ParseTask, error) {
	// 序列化选项
	optionsBytes, err := json.Marshal(task.Options)
	if err != nil {
		return nil, err
	}

	po := &ParseTaskModel{
		ID:        task.ID,
		ResumeID:  task.ResumeID,
		UserID:    task.UserID,
		FilePath:  task.FilePath,
		FileType:  task.FileType,
		Status:    task.Status,
		Progress:  task.Progress,
		Options:   string(optionsBytes),
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}

	if err := r.data.db.WithContext(ctx).Create(po).Error; err != nil {
		return nil, err
	}

	return r.poToBiz(po), nil
}

func (r *parseTaskRepo) GetTask(ctx context.Context, taskID string) (*biz.ParseTask, error) {
	var po ParseTaskModel
	if err := r.data.db.WithContext(ctx).Where("id = ?", taskID).First(&po).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, biz.ErrTaskNotFound
		}
		return nil, err
	}

	return r.poToBiz(&po), nil
}

func (r *parseTaskRepo) UpdateTask(ctx context.Context, task *biz.ParseTask) error {
	// 序列化结果
	var resultStr string
	if task.Result != nil {
		resultBytes, err := json.Marshal(task.Result)
		if err != nil {
			return err
		}
		resultStr = string(resultBytes)
	}

	updates := map[string]interface{}{
		"status":     task.Status,
		"progress":   task.Progress,
		"result":     resultStr,
		"error_msg":  task.ErrorMsg,
		"updated_at": time.Now(),
	}

	if task.Status == "completed" || task.Status == "failed" {
		now := time.Now()
		updates["completed_at"] = &now
	}

	return r.data.db.WithContext(ctx).
		Model(&ParseTaskModel{}).
		Where("id = ?", task.ID).
		Updates(updates).Error
}

func (r *parseTaskRepo) ListTasksByUser(ctx context.Context, userID string, limit, offset int) ([]*biz.ParseTask, error) {
	var pos []ParseTaskModel
	if err := r.data.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&pos).Error; err != nil {
		return nil, err
	}

	tasks := make([]*biz.ParseTask, len(pos))
	for i, po := range pos {
		tasks[i] = r.poToBiz(&po)
	}

	return tasks, nil
}

func (r *parseTaskRepo) DeleteTask(ctx context.Context, taskID string) error {
	return r.data.db.WithContext(ctx).
		Where("id = ?", taskID).
		Delete(&ParseTaskModel{}).Error
}

// poToBiz 将数据模型转换为业务模型
func (r *parseTaskRepo) poToBiz(po *ParseTaskModel) *biz.ParseTask {
	task := &biz.ParseTask{
		ID:          po.ID,
		ResumeID:    po.ResumeID,
		UserID:      po.UserID,
		FilePath:    po.FilePath,
		FileType:    po.FileType,
		Status:      po.Status,
		Progress:    po.Progress,
		ErrorMsg:    po.ErrorMsg,
		CreatedAt:   po.CreatedAt,
		UpdatedAt:   po.UpdatedAt,
		CompletedAt: po.CompletedAt,
	}

	// 反序列化选项
	if po.Options != "" {
		var options biz.ParseOptions
		if err := json.Unmarshal([]byte(po.Options), &options); err == nil {
			task.Options = &options
		}
	}

	// 反序列化结果
	if po.Result != "" {
		var result biz.ParsedContent
		if err := json.Unmarshal([]byte(po.Result), &result); err == nil {
			task.Result = &result
		}
	}

	return task
}
