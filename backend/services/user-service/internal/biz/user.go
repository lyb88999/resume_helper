package biz

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	v1 "github.com/liyubo06/resumeOptim_claude/api/user/v1"
	"github.com/liyubo06/resumeOptim_claude/backend/services/user-service/internal/conf"
)

// User 业务层用户模型
type User struct {
	ID        uint64    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Nickname  string    `json:"nickname"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRepo 用户仓库接口
type UserRepo interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id uint64) (*User, error)
	UpdateUser(ctx context.Context, user *User) (*User, error)
	CheckEmailExists(ctx context.Context, email string) (bool, error)
}

// 错误定义
var (
	ErrUserNotFound     = errors.NotFound("USER_NOT_FOUND", "用户不存在")
	ErrEmailExists      = errors.BadRequest("EMAIL_EXISTS", "邮箱已存在")
	ErrInvalidPassword  = errors.Unauthorized("INVALID_PASSWORD", "密码错误")
	ErrInvalidEmail     = errors.BadRequest("INVALID_EMAIL", "邮箱格式错误")
	ErrPasswordTooShort = errors.BadRequest("PASSWORD_TOO_SHORT", "密码长度不能少于6位")
)

// UserUsecase 用户业务逻辑
type UserUsecase struct {
	repo UserRepo
	auth *conf.Auth
	log  *log.Helper
}

// NewUserUsecase 创建用户业务逻辑实例
func NewUserUsecase(repo UserRepo, auth *conf.Auth, logger log.Logger) *UserUsecase {
	return &UserUsecase{
		repo: repo,
		auth: auth,
		log:  log.NewHelper(logger),
	}
}

// Register 用户注册
func (uc *UserUsecase) Register(ctx context.Context, req *v1.RegisterRequest) (*v1.RegisterReply, error) {
	// 验证输入
	if err := uc.validateRegisterRequest(req); err != nil {
		return nil, err
	}

	// 检查邮箱是否已存在
	exists, err := uc.repo.CheckEmailExists(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailExists
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %w", err)
	}

	// 创建用户
	user := &User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
	}

	createdUser, err := uc.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &v1.RegisterReply{
		Id:        createdUser.ID,
		Email:     createdUser.Email,
		Nickname:  createdUser.Nickname,
		CreatedAt: createdUser.CreatedAt.Format(time.RFC3339),
	}, nil
}

// Login 用户登录
func (uc *UserUsecase) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginReply, error) {
	// 获取用户
	user, err := uc.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == ErrUserNotFound {
			return nil, ErrInvalidPassword // 为了安全，不暴露用户是否存在
		}
		return nil, err
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidPassword
	}

	// 生成JWT token
	token, expiresAt, err := uc.generateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &v1.LoginReply{
		Token:     token,
		ExpiresAt: uint64(expiresAt),
		User: &v1.UserInfo{
			Id:        user.ID,
			Email:     user.Email,
			Nickname:  user.Nickname,
			Avatar:    user.Avatar,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

// GetUserInfo 获取用户信息
func (uc *UserUsecase) GetUserInfo(ctx context.Context, userID uint64) (*v1.GetUserInfoReply, error) {
	user, err := uc.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &v1.GetUserInfoReply{
		User: &v1.UserInfo{
			Id:        user.ID,
			Email:     user.Email,
			Nickname:  user.Nickname,
			Avatar:    user.Avatar,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

// UpdateUserInfo 更新用户信息
func (uc *UserUsecase) UpdateUserInfo(ctx context.Context, userID uint64, req *v1.UpdateUserInfoRequest) error {
	user := &User{
		ID:       userID,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
	}

	_, err := uc.repo.UpdateUser(ctx, user)
	return err
}

// generateToken 生成JWT token
func (uc *UserUsecase) generateToken(userID uint64) (string, int64, error) {
	expiresAt := time.Now().Add(uc.auth.JwtExpire.AsDuration()).Unix()

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expiresAt,
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(uc.auth.JwtSecret))
	if err != nil {
		return "", 0, fmt.Errorf("生成token失败: %w", err)
	}

	return tokenString, expiresAt, nil
}

// validateRegisterRequest 验证注册请求
func (uc *UserUsecase) validateRegisterRequest(req *v1.RegisterRequest) error {
	if req.Email == "" {
		return ErrInvalidEmail
	}
	if len(req.Password) < 6 {
		return ErrPasswordTooShort
	}
	if req.Nickname == "" {
		req.Nickname = "用户" + fmt.Sprintf("%d", time.Now().Unix())
	}
	return nil
}
