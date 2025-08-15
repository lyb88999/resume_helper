package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"

	v1 "github.com/lyb88999/resume_helper/api/user/v1"
	"github.com/lyb88999/resume_helper/backend/services/user-service/internal/biz"
)

// UserService 用户服务实现
type UserService struct {
	v1.UnimplementedUserServiceServer

	uc  *biz.UserUsecase
	log *log.Helper
}

// NewUserService 创建用户服务实例
func NewUserService(uc *biz.UserUsecase, logger log.Logger) *UserService {
	return &UserService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

// Register 用户注册
func (s *UserService) Register(ctx context.Context, req *v1.RegisterRequest) (*v1.RegisterReply, error) {
	s.log.WithContext(ctx).Infof("用户注册请求: email=%s", req.Email)

	reply, err := s.uc.Register(ctx, req)
	if err != nil {
		s.log.WithContext(ctx).Errorf("用户注册失败: %v", err)
		return nil, err
	}

	s.log.WithContext(ctx).Infof("用户注册成功: id=%d", reply.Id)
	return reply, nil
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginReply, error) {
	s.log.WithContext(ctx).Infof("用户登录请求: email=%s", req.Email)

	reply, err := s.uc.Login(ctx, req)
	if err != nil {
		s.log.WithContext(ctx).Errorf("用户登录失败: %v", err)
		return nil, err
	}

	s.log.WithContext(ctx).Infof("用户登录成功: id=%d", reply.User.Id)
	return reply, nil
}

// GetUserInfo 获取用户信息
func (s *UserService) GetUserInfo(ctx context.Context, req *v1.GetUserInfoRequest) (*v1.GetUserInfoReply, error) {
	s.log.WithContext(ctx).Infof("获取用户信息请求: id=%d", req.Id)

	reply, err := s.uc.GetUserInfo(ctx, req.Id)
	if err != nil {
		s.log.WithContext(ctx).Errorf("获取用户信息失败: %v", err)
		return nil, err
	}

	return reply, nil
}

// UpdateUserInfo 更新用户信息
func (s *UserService) UpdateUserInfo(ctx context.Context, req *v1.UpdateUserInfoRequest) (*emptypb.Empty, error) {
	s.log.WithContext(ctx).Infof("更新用户信息请求: id=%d", req.Id)

	err := s.uc.UpdateUserInfo(ctx, req.Id, req)
	if err != nil {
		s.log.WithContext(ctx).Errorf("更新用户信息失败: %v", err)
		return nil, err
	}

	s.log.WithContext(ctx).Infof("更新用户信息成功: id=%d", req.Id)
	return &emptypb.Empty{}, nil
}
