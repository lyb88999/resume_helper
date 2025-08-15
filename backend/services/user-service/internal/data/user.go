package data

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"

	"github.com/lyb88999/resume_helper/backend/services/user-service/internal/biz"
	"github.com/lyb88999/resume_helper/backend/shared/pkg/models"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserRepo creates a new user repository.
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *userRepo) CreateUser(ctx context.Context, user *biz.User) (*biz.User, error) {
	// 转换业务层模型到数据层模型
	u := &models.User{
		Email:    user.Email,
		Password: user.Password,
		Nickname: user.Nickname,
	}

	if err := r.data.db.WithContext(ctx).Create(u).Error; err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	// 转换回业务层模型
	return &biz.User{
		ID:        u.ID,
		Email:     u.Email,
		Password:  u.Password,
		Nickname:  u.Nickname,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}

func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*biz.User, error) {
	u := &models.User{}
	if err := r.data.db.WithContext(ctx).Where("email = ?", email).First(u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, biz.ErrUserNotFound
		}
		return nil, fmt.Errorf("通过邮箱查询用户失败: %w", err)
	}

	return &biz.User{
		ID:        u.ID,
		Email:     u.Email,
		Password:  u.Password,
		Nickname:  u.Nickname,
		Avatar:    u.Avatar,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}

func (r *userRepo) GetUserByID(ctx context.Context, id uint64) (*biz.User, error) {
	u := &models.User{}
	if err := r.data.db.WithContext(ctx).Where("id = ?", id).First(u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, biz.ErrUserNotFound
		}
		return nil, fmt.Errorf("通过ID查询用户失败: %w", err)
	}

	return &biz.User{
		ID:        u.ID,
		Email:     u.Email,
		Password:  u.Password,
		Nickname:  u.Nickname,
		Avatar:    u.Avatar,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}

func (r *userRepo) UpdateUser(ctx context.Context, user *biz.User) (*biz.User, error) {
	u := &models.User{}
	if err := r.data.db.WithContext(ctx).Where("id = ?", user.ID).First(u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, biz.ErrUserNotFound
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	// 更新字段
	if user.Nickname != "" {
		u.Nickname = user.Nickname
	}
	if user.Avatar != "" {
		u.Avatar = user.Avatar
	}
	u.UpdatedAt = time.Now()

	if err := r.data.db.WithContext(ctx).Save(u).Error; err != nil {
		return nil, fmt.Errorf("更新用户失败: %w", err)
	}

	return &biz.User{
		ID:        u.ID,
		Email:     u.Email,
		Password:  u.Password,
		Nickname:  u.Nickname,
		Avatar:    u.Avatar,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}

func (r *userRepo) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.data.db.WithContext(ctx).Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, fmt.Errorf("检查邮箱是否存在失败: %w", err)
	}
	return count > 0, nil
}
