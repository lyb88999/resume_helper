package service

import (
	"context"
	"fmt"
	"time"

	pb "github.com/lyb88999/resume_helper/backend/services/parser-service/api/parser/v1"
	"github.com/lyb88999/resume_helper/backend/services/parser-service/internal/biz"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// ParserService 解析服务实现
type ParserService struct {
	uc *biz.ParserUsecase
}

// NewParserService 创建解析服务
func NewParserService(uc *biz.ParserUsecase) *ParserService {
	// 注册解析器
	uc.RegisterParser("txt", biz.NewTextParser())
	uc.RegisterParser("md", biz.NewMarkdownParser())
	uc.RegisterParser("markdown", biz.NewMarkdownParser())
	uc.RegisterParser("pdf", biz.NewPDFParser())
	uc.RegisterParser("docx", biz.NewDocxParser())
	uc.RegisterParser("doc", biz.NewDocxParser())

	return &ParserService{
		uc: uc,
	}
}

func (s *ParserService) ParseDocument(ctx context.Context, req *pb.ParseDocumentRequest) (*pb.ParseDocumentReply, error) {
	// 转换解析选项
	var options *biz.ParseOptions
	if req.Options != nil {
		options = &biz.ParseOptions{
			ExtractImages:  req.Options.ExtractImages,
			CleanText:      req.Options.CleanText,
			TargetLanguage: req.Options.TargetLanguage,
			SkipSections:   req.Options.SkipSections,
		}
	}

	// 调用业务逻辑
	task, err := s.uc.ParseDocument(ctx, req.FilePath, req.FileType, req.ResumeId, req.UserId, options)
	if err != nil {
		return nil, err
	}

	// 转换响应
	reply := &pb.ParseDocumentReply{
		TaskId:    task.ID,
		Status:    task.Status,
		Message:   "解析任务已创建",
		CreatedAt: timestamppb.New(task.CreatedAt),
	}

	// 如果已经有结果，则包含结果
	if task.Result != nil {
		reply.Content = s.convertParsedContent(task.Result)
	}

	return reply, nil
}

func (s *ParserService) GetParseStatus(ctx context.Context, req *pb.GetParseStatusRequest) (*pb.GetParseStatusReply, error) {
	task, err := s.uc.GetParseStatus(ctx, req.TaskId)
	if err != nil {
		if err == biz.ErrTaskNotFound {
			return nil, fmt.Errorf("task not found: %s", req.TaskId)
		}
		return nil, err
	}

	reply := &pb.GetParseStatusReply{
		TaskId:    task.ID,
		Status:    task.Status,
		Progress:  int32(task.Progress),
		Message:   task.ErrorMsg,
		CreatedAt: timestamppb.New(task.CreatedAt),
		UpdatedAt: timestamppb.New(task.UpdatedAt),
	}

	if task.Result != nil {
		reply.Content = s.convertParsedContent(task.Result)
	}

	return reply, nil
}

func (s *ParserService) Health(ctx context.Context, req *pb.HealthRequest) (*pb.HealthReply, error) {
	return &pb.HealthReply{
		Status:    "ok",
		Service:   "parser-service",
		Timestamp: timestamppb.New(time.Now()),
	}, nil
}

// convertParsedContent 转换解析内容为protobuf格式
func (s *ParserService) convertParsedContent(content *biz.ParsedContent) *pb.ParsedContent {
	result := &pb.ParsedContent{
		RawText: content.RawText,
	}

	// 转换个人信息
	if content.PersonalInfo != nil {
		result.PersonalInfo = &pb.PersonalInfo{
			Name:        content.PersonalInfo.Name,
			Phone:       content.PersonalInfo.Phone,
			Email:       content.PersonalInfo.Email,
			Address:     content.PersonalInfo.Address,
			BirthDate:   content.PersonalInfo.BirthDate,
			Gender:      content.PersonalInfo.Gender,
			Nationality: content.PersonalInfo.Nationality,
			SocialLinks: content.PersonalInfo.SocialLinks,
			AvatarUrl:   content.PersonalInfo.AvatarURL,
		}
	}

	// 转换教育背景
	for _, edu := range content.Education {
		result.Education = append(result.Education, &pb.Education{
			School:      edu.School,
			Degree:      edu.Degree,
			Major:       edu.Major,
			StartDate:   edu.StartDate,
			EndDate:     edu.EndDate,
			Gpa:         edu.GPA,
			Description: edu.Description,
			Courses:     edu.Courses,
		})
	}

	// 转换工作经历
	for _, exp := range content.Experience {
		result.Experience = append(result.Experience, &pb.Experience{
			Company:          exp.Company,
			Position:         exp.Position,
			StartDate:        exp.StartDate,
			EndDate:          exp.EndDate,
			Location:         exp.Location,
			Department:       exp.Department,
			Responsibilities: exp.Responsibilities,
			Achievements:     exp.Achievements,
			Technologies:     exp.Technologies,
		})
	}

	// 转换项目经历
	for _, proj := range content.Projects {
		result.Projects = append(result.Projects, &pb.Project{
			Name:         proj.Name,
			Role:         proj.Role,
			StartDate:    proj.StartDate,
			EndDate:      proj.EndDate,
			Description:  proj.Description,
			Technologies: proj.Technologies,
			Achievements: proj.Achievements,
			Url:          proj.URL,
			Company:      proj.Company,
		})
	}

	// 转换技能
	if content.Skills != nil {
		result.Skills = &pb.Skills{
			Languages:      content.Skills.Languages,
			Certifications: content.Skills.Certifications,
			Awards:         content.Skills.Awards,
		}

		for _, cat := range content.Skills.Categories {
			pbCat := &pb.SkillCategory{
				Category: cat.Category,
			}

			for _, skill := range cat.Skills {
				pbCat.Skills = append(pbCat.Skills, &pb.SkillItem{
					Name:  skill.Name,
					Level: skill.Level,
					Years: skill.Years,
				})
			}

			result.Skills.Categories = append(result.Skills.Categories, pbCat)
		}
	}

	// 转换元数据
	if content.Metadata != nil {
		result.Metadata = &pb.ParseMetadata{
			FileSize:        content.Metadata.FileSize,
			PageCount:       content.Metadata.PageCount,
			ParseDuration:   content.Metadata.ParseDuration,
			ParserVersion:   content.Metadata.ParserVersion,
			Warnings:        content.Metadata.Warnings,
			ConfidenceScore: content.Metadata.ConfidenceScore,
		}
	}

	return result
}
