# 简历优化系统 - API设计文档

## 1. API概述

### 1.1 设计原则
- **RESTful设计**: 遵循REST架构风格
- **统一格式**: 统一的请求响应格式
- **版本控制**: 支持API版本管理
- **安全性**: 完整的认证授权机制
- **可扩展性**: 支持向后兼容的功能扩展

### 1.2 基础信息
- **Base URL**: `https://api.resumeoptim.com`
- **API版本**: `v1`
- **协议**: HTTPS
- **数据格式**: JSON
- **字符编码**: UTF-8

### 1.3 环境地址
```
开发环境: https://dev-api.resumeoptim.com
测试环境: https://test-api.resumeoptim.com
生产环境: https://api.resumeoptim.com
```

## 2. 认证和授权

### 2.1 认证方式

#### JWT Token认证
```http
Authorization: Bearer <jwt_token>
```

#### API Key认证（企业用户）
```http
X-API-Key: <api_key>
X-API-Secret: <api_secret>
```

### 2.2 Token获取

#### 用户登录获取Token
```http
POST /api/v1/auth/login
Content-Type: application/json

{
    "email": "user@example.com",
    "password": "password123",
    "remember_me": true
}
```

**响应示例**:
```json
{
    "code": 200,
    "message": "登录成功",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "refresh_token": "refresh_token_here",
        "expires_in": 86400,
        "user": {
            "id": 12345,
            "username": "user123",
            "email": "user@example.com",
            "user_type": "premium"
        }
    }
}
```

### 2.3 权限控制

#### 用户类型权限
- **免费用户**: 基础功能，每日5次分析
- **高级用户**: 完整功能，每日50次分析
- **企业用户**: API访问，无限制使用

#### 接口权限标识
```http
X-Required-Permission: resume:read
X-Required-Role: premium
```

## 3. 通用响应格式

### 3.1 成功响应
```json
{
    "code": 200,
    "message": "操作成功",
    "data": {
        // 具体数据内容
    },
    "meta": {
        "page": 1,
        "per_page": 20,
        "total": 100,
        "total_pages": 5
    },
    "timestamp": "2025-01-15T10:30:00Z"
}
```

### 3.2 错误响应
```json
{
    "code": 400,
    "message": "请求参数错误",
    "error": {
        "type": "validation_error",
        "details": [
            {
                "field": "email",
                "message": "邮箱格式不正确"
            }
        ]
    },
    "timestamp": "2025-01-15T10:30:00Z"
}
```

### 3.3 HTTP状态码
| 状态码 | 含义 | 说明 |
|--------|------|------|
| 200 | OK | 请求成功 |
| 201 | Created | 资源创建成功 |
| 204 | No Content | 请求成功，无返回内容 |
| 400 | Bad Request | 请求参数错误 |
| 401 | Unauthorized | 未授权 |
| 403 | Forbidden | 权限不足 |
| 404 | Not Found | 资源不存在 |
| 409 | Conflict | 资源冲突 |
| 422 | Unprocessable Entity | 请求格式正确但语义错误 |
| 429 | Too Many Requests | 请求过于频繁 |
| 500 | Internal Server Error | 服务器内部错误 |

## 4. 用户认证模块

### 4.1 用户注册
```http
POST /api/v1/auth/register
Content-Type: application/json

{
    "username": "user123",
    "email": "user@example.com",
    "password": "password123",
    "confirm_password": "password123",
    "invite_code": "INVITE123" // 可选
}
```

**响应**:
```json
{
    "code": 201,
    "message": "注册成功",
    "data": {
        "user_id": 12345,
        "username": "user123",
        "email": "user@example.com",
        "user_type": "free",
        "email_verified": false
    }
}
```

### 4.2 邮箱验证
```http
POST /api/v1/auth/verify-email
Content-Type: application/json

{
    "email": "user@example.com",
    "verification_code": "123456"
}
```

### 4.3 忘记密码
```http
POST /api/v1/auth/forgot-password
Content-Type: application/json

{
    "email": "user@example.com"
}
```

### 4.4 重置密码
```http
POST /api/v1/auth/reset-password
Content-Type: application/json

{
    "reset_token": "reset_token_here",
    "new_password": "new_password123",
    "confirm_password": "new_password123"
}
```

### 4.5 刷新Token
```http
POST /api/v1/auth/refresh
Content-Type: application/json

{
    "refresh_token": "refresh_token_here"
}
```

### 4.6 用户登出
```http
POST /api/v1/auth/logout
Authorization: Bearer <jwt_token>
```

## 5. 用户管理模块

### 5.1 获取用户信息
```http
GET /api/v1/users/profile
Authorization: Bearer <jwt_token>
```

**响应**:
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "id": 12345,
        "username": "user123",
        "email": "user@example.com",
        "user_type": "premium",
        "avatar": "https://cdn.example.com/avatar/12345.jpg",
        "created_at": "2025-01-01T00:00:00Z",
        "last_login": "2025-01-15T10:30:00Z",
        "preferences": {
            "target_industry": "technology",
            "target_position": "软件工程师",
            "experience_level": "mid"
        },
        "quota": {
            "daily_limit": 50,
            "daily_used": 5,
            "monthly_limit": 1000,
            "monthly_used": 150
        }
    }
}
```

### 5.2 更新用户信息
```http
PUT /api/v1/users/profile
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
    "username": "new_username",
    "avatar": "https://cdn.example.com/avatar/new.jpg",
    "preferences": {
        "target_industry": "finance",
        "target_position": "数据分析师",
        "experience_level": "senior"
    }
}
```

### 5.3 修改密码
```http
PUT /api/v1/users/password
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
    "current_password": "old_password",
    "new_password": "new_password123",
    "confirm_password": "new_password123"
}
```

### 5.4 用户统计
```http
GET /api/v1/users/stats
Authorization: Bearer <jwt_token>
```

**响应**:
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "total_resumes": 15,
        "total_analyses": 45,
        "success_rate": 0.95,
        "average_score": 8.2,
        "recent_activity": [
            {
                "type": "resume_upload",
                "timestamp": "2025-01-15T09:30:00Z",
                "description": "上传简历：软件工程师简历.pdf"
            }
        ]
    }
}
```

## 6. 文件管理模块

### 6.1 文件上传
```http
POST /api/v1/files/upload
Authorization: Bearer <jwt_token>
Content-Type: multipart/form-data

file=@resume.pdf
title=我的简历
description=2025年最新版本
```

**响应**:
```json
{
    "code": 201,
    "message": "上传成功",
    "data": {
        "file_id": "file_12345",
        "filename": "resume.pdf",
        "original_name": "我的简历.pdf",
        "size": 1024000,
        "mime_type": "application/pdf",
        "url": "https://cdn.example.com/files/file_12345.pdf",
        "status": "uploaded",
        "created_at": "2025-01-15T10:30:00Z"
    }
}
```

### 6.2 获取文件列表
```http
GET /api/v1/files?page=1&per_page=20&type=pdf
Authorization: Bearer <jwt_token>
```

**查询参数**:
- `page`: 页码（默认1）
- `per_page`: 每页数量（默认20，最大100）
- `type`: 文件类型（pdf, docx, md）
- `status`: 文件状态（uploaded, processing, completed, error）
- `created_after`: 创建时间起始
- `created_before`: 创建时间结束

### 6.3 获取文件详情
```http
GET /api/v1/files/{file_id}
Authorization: Bearer <jwt_token>
```

### 6.4 下载文件
```http
GET /api/v1/files/{file_id}/download
Authorization: Bearer <jwt_token>
```

### 6.5 删除文件
```http
DELETE /api/v1/files/{file_id}
Authorization: Bearer <jwt_token>
```

## 7. 简历管理模块

### 7.1 创建简历记录
```http
POST /api/v1/resumes
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
    "title": "软件工程师简历",
    "file_id": "file_12345",
    "target_industry": "technology",
    "target_position": "软件工程师",
    "description": "申请高级软件工程师职位"
}
```

**响应**:
```json
{
    "code": 201,
    "message": "创建成功",
    "data": {
        "resume_id": "resume_12345",
        "title": "软件工程师简历",
        "file_id": "file_12345",
        "status": "created",
        "target_industry": "technology",
        "target_position": "软件工程师",
        "created_at": "2025-01-15T10:30:00Z"
    }
}
```

### 7.2 获取简历列表
```http
GET /api/v1/resumes?page=1&per_page=20&status=analyzed
Authorization: Bearer <jwt_token>
```

**查询参数**:
- `page`: 页码
- `per_page`: 每页数量
- `status`: 状态（created, processing, analyzed, error）
- `industry`: 目标行业
- `position`: 目标职位
- `created_after`: 创建时间起始
- `created_before`: 创建时间结束

**响应**:
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "resumes": [
            {
                "resume_id": "resume_12345",
                "title": "软件工程师简历",
                "status": "analyzed",
                "score": 8.5,
                "target_industry": "technology",
                "target_position": "软件工程师",
                "created_at": "2025-01-15T10:30:00Z",
                "analyzed_at": "2025-01-15T10:32:00Z"
            }
        ]
    },
    "meta": {
        "page": 1,
        "per_page": 20,
        "total": 5,
        "total_pages": 1
    }
}
```

### 7.3 获取简历详情
```http
GET /api/v1/resumes/{resume_id}
Authorization: Bearer <jwt_token>
```

**响应**:
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "resume_id": "resume_12345",
        "title": "软件工程师简历",
        "status": "analyzed",
        "file": {
            "file_id": "file_12345",
            "filename": "resume.pdf",
            "url": "https://cdn.example.com/files/file_12345.pdf"
        },
        "parsed_content": {
            "personal_info": {
                "name": "张三",
                "phone": "138****1234",
                "email": "zhangsan@example.com"
            },
            "education": [...],
            "experience": [...],
            "skills": [...],
            "projects": [...]
        },
        "target_industry": "technology",
        "target_position": "软件工程师",
        "created_at": "2025-01-15T10:30:00Z",
        "updated_at": "2025-01-15T10:32:00Z"
    }
}
```

### 7.4 删除简历
```http
DELETE /api/v1/resumes/{resume_id}
Authorization: Bearer <jwt_token>
```

## 8. 简历分析模块

### 8.1 开始分析
```http
POST /api/v1/resumes/{resume_id}/analyze
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
    "analysis_type": "comprehensive", // basic, comprehensive, detailed
    "target_industry": "technology",
    "target_position": "高级软件工程师",
    "callback_url": "https://your-app.com/webhook/analysis" // 可选
}
```

**响应**:
```json
{
    "code": 202,
    "message": "分析已开始",
    "data": {
        "analysis_id": "analysis_12345",
        "resume_id": "resume_12345",
        "status": "processing",
        "estimated_time": 60,
        "started_at": "2025-01-15T10:30:00Z"
    }
}
```

### 8.2 获取分析状态
```http
GET /api/v1/resumes/{resume_id}/analysis/status
Authorization: Bearer <jwt_token>
```

**响应**:
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "analysis_id": "analysis_12345",
        "status": "processing", // processing, completed, error
        "progress": 75,
        "current_step": "生成优化建议",
        "estimated_remaining": 15,
        "started_at": "2025-01-15T10:30:00Z"
    }
}
```

### 8.3 获取分析结果
```http
GET /api/v1/resumes/{resume_id}/analysis
Authorization: Bearer <jwt_token>
```

**响应**:
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "analysis_id": "analysis_12345",
        "resume_id": "resume_12345",
        "status": "completed",
        "overall_score": 8.5,
        "scores": {
            "content_quality": 8.2,
            "structure": 8.8,
            "keyword_match": 7.9,
            "industry_fit": 8.7
        },
        "analysis_result": {
            "strengths": [
                "技术技能描述清晰明确",
                "项目经验丰富且相关性强"
            ],
            "weaknesses": [
                "缺少量化的工作成果描述",
                "教育背景部分信息不够完整"
            ],
            "suggestions": {
                "content": [...],
                "structure": [...],
                "format": [...],
                "industry": [...]
            }
        },
        "completed_at": "2025-01-15T10:32:00Z"
    }
}
```

### 8.4 重新分析
```http
POST /api/v1/resumes/{resume_id}/reanalyze
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
    "reason": "更新了目标职位",
    "target_position": "架构师"
}
```

### 8.5 导出分析报告
```http
GET /api/v1/resumes/{resume_id}/analysis/report?format=pdf
Authorization: Bearer <jwt_token>
```

**查询参数**:
- `format`: 导出格式（pdf, word, json）
- `include_suggestions`: 是否包含建议（true/false）
- `template`: 报告模板（standard, detailed, simple）

## 9. 知识库管理模块

### 9.1 添加知识文档
```http
POST /api/v1/knowledge/docs
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
    "title": "软件工程师简历最佳实践",
    "content": "# 软件工程师简历编写指南\n\n## 技能描述\n...",
    "category": "technology",
    "tags": ["软件工程师", "技能描述", "项目经验"],
    "industry": "technology",
    "position": "软件工程师",
    "language": "zh-CN",
    "source": "internal" // internal, external, uploaded
}
```

**响应**:
```json
{
    "code": 201,
    "message": "添加成功",
    "data": {
        "doc_id": "doc_12345",
        "title": "软件工程师简历最佳实践",
        "category": "technology",
        "status": "indexing",
        "created_at": "2025-01-15T10:30:00Z"
    }
}
```

### 9.2 获取知识文档列表
```http
GET /api/v1/knowledge/docs?category=technology&page=1&per_page=20
Authorization: Bearer <jwt_token>
```

**查询参数**:
- `category`: 分类
- `industry`: 行业
- `position`: 职位
- `tags`: 标签（多个用逗号分隔）
- `language`: 语言
- `status`: 状态（active, archived）
- `created_after`: 创建时间起始
- `created_before`: 创建时间结束

### 9.3 搜索知识文档
```http
GET /api/v1/knowledge/search?q=软件工程师技能&limit=10
Authorization: Bearer <jwt_token>
```

**查询参数**:
- `q`: 搜索关键词
- `limit`: 返回结果数量（默认10，最大50）
- `category`: 限制搜索分类
- `industry`: 限制搜索行业
- `min_score`: 最小相似度分数（0-1）

**响应**:
```json
{
    "code": 200,
    "message": "搜索成功",
    "data": {
        "query": "软件工程师技能",
        "results": [
            {
                "doc_id": "doc_12345",
                "title": "软件工程师简历最佳实践",
                "content_snippet": "软件工程师在简历中应该突出...",
                "score": 0.95,
                "category": "technology",
                "tags": ["软件工程师", "技能描述"]
            }
        ],
        "total": 15,
        "search_time": 0.05
    }
}
```

### 9.4 获取知识文档详情
```http
GET /api/v1/knowledge/docs/{doc_id}
Authorization: Bearer <jwt_token>
```

### 9.5 更新知识文档
```http
PUT /api/v1/knowledge/docs/{doc_id}
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
    "title": "更新后的标题",
    "content": "更新后的内容",
    "tags": ["新标签1", "新标签2"]
}
```

### 9.6 删除知识文档
```http
DELETE /api/v1/knowledge/docs/{doc_id}
Authorization: Bearer <jwt_token>
```

## 10. 通知管理模块

### 10.1 获取通知列表
```http
GET /api/v1/notifications?page=1&per_page=20&type=analysis
Authorization: Bearer <jwt_token>
```

**查询参数**:
- `type`: 通知类型（analysis, system, promotion）
- `status`: 状态（unread, read, archived）
- `created_after`: 创建时间起始

**响应**:
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "notifications": [
            {
                "id": "notification_12345",
                "type": "analysis",
                "title": "简历分析完成",
                "content": "您的简历《软件工程师简历》分析已完成，总分8.5分",
                "status": "unread",
                "data": {
                    "resume_id": "resume_12345",
                    "analysis_id": "analysis_12345"
                },
                "created_at": "2025-01-15T10:32:00Z"
            }
        ]
    },
    "meta": {
        "page": 1,
        "per_page": 20,
        "total": 5,
        "total_pages": 1,
        "unread_count": 3
    }
}
```

### 10.2 标记通知为已读
```http
PUT /api/v1/notifications/{notification_id}/read
Authorization: Bearer <jwt_token>
```

### 10.3 批量标记已读
```http
PUT /api/v1/notifications/batch-read
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
    "notification_ids": ["notification_12345", "notification_12346"],
    "mark_all": false // true表示标记所有未读通知
}
```

### 10.4 删除通知
```http
DELETE /api/v1/notifications/{notification_id}
Authorization: Bearer <jwt_token>
```

## 11. 系统管理模块

### 11.1 获取系统状态
```http
GET /api/v1/system/status
```

**响应**:
```json
{
    "code": 200,
    "message": "系统正常",
    "data": {
        "status": "healthy",
        "version": "1.0.0",
        "uptime": 3600,
        "services": {
            "database": "healthy",
            "redis": "healthy",
            "ai_service": "healthy",
            "file_storage": "healthy"
        },
        "performance": {
            "avg_response_time": 150,
            "requests_per_second": 1000,
            "error_rate": 0.01
        }
    }
}
```

### 11.2 获取配置信息
```http
GET /api/v1/system/config
Authorization: Bearer <jwt_token>
```

**响应**:
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "file_upload": {
            "max_size": 10485760,
            "allowed_types": ["pdf", "docx", "md"],
            "max_files_per_day": 10
        },
        "analysis": {
            "max_concurrent": 5,
            "timeout": 300,
            "retry_limit": 3
        },
        "features": {
            "batch_analysis": true,
            "api_access": false,
            "advanced_templates": true
        }
    }
}
```

## 12. 数据模型定义

### 12.1 用户模型
```json
{
    "User": {
        "id": "integer",
        "username": "string",
        "email": "string",
        "user_type": "enum[free, premium, enterprise]",
        "avatar": "string",
        "preferences": {
            "target_industry": "string",
            "target_position": "string",
            "experience_level": "enum[entry, mid, senior, executive]",
            "notification_settings": {
                "email": "boolean",
                "push": "boolean",
                "sms": "boolean"
            }
        },
        "quota": {
            "daily_limit": "integer",
            "daily_used": "integer",
            "monthly_limit": "integer",
            "monthly_used": "integer"
        },
        "created_at": "datetime",
        "updated_at": "datetime",
        "last_login": "datetime"
    }
}
```

### 12.2 简历模型
```json
{
    "Resume": {
        "resume_id": "string",
        "user_id": "integer",
        "title": "string",
        "file_id": "string",
        "status": "enum[created, processing, analyzed, error]",
        "parsed_content": {
            "personal_info": {
                "name": "string",
                "phone": "string",
                "email": "string",
                "address": "string",
                "linkedin": "string",
                "github": "string"
            },
            "education": [
                {
                    "degree": "string",
                    "major": "string",
                    "school": "string",
                    "graduation_date": "string",
                    "gpa": "string",
                    "achievements": ["string"]
                }
            ],
            "experience": [
                {
                    "company": "string",
                    "position": "string",
                    "duration": "string",
                    "location": "string",
                    "responsibilities": ["string"],
                    "achievements": ["string"]
                }
            ],
            "skills": {
                "technical": ["string"],
                "languages": ["string"],
                "soft_skills": ["string"],
                "certifications": ["string"]
            },
            "projects": [
                {
                    "name": "string",
                    "description": "string",
                    "technologies": ["string"],
                    "duration": "string",
                    "url": "string",
                    "achievements": ["string"]
                }
            ]
        },
        "target_industry": "string",
        "target_position": "string",
        "created_at": "datetime",
        "updated_at": "datetime"
    }
}
```

### 12.3 分析结果模型
```json
{
    "AnalysisResult": {
        "analysis_id": "string",
        "resume_id": "string",
        "status": "enum[processing, completed, error]",
        "overall_score": "float",
        "scores": {
            "content_quality": "float",
            "structure": "float",
            "keyword_match": "float",
            "industry_fit": "float",
            "formatting": "float"
        },
        "analysis_result": {
            "strengths": ["string"],
            "weaknesses": ["string"],
            "missing_elements": ["string"],
            "keyword_analysis": {
                "matched_keywords": ["string"],
                "missing_keywords": ["string"],
                "keyword_density": "float"
            }
        },
        "suggestions": {
            "content": [
                {
                    "section": "string",
                    "issue": "string",
                    "suggestion": "string",
                    "priority": "enum[high, medium, low]",
                    "examples": ["string"]
                }
            ],
            "structure": [...],
            "format": [...],
            "industry": [...]
        },
        "started_at": "datetime",
        "completed_at": "datetime"
    }
}
```

## 13. 错误码定义

### 13.1 通用错误码
| 错误码 | HTTP状态码 | 说明 |
|--------|------------|------|
| 10001 | 400 | 请求参数错误 |
| 10002 | 401 | 未授权访问 |
| 10003 | 403 | 权限不足 |
| 10004 | 404 | 资源不存在 |
| 10005 | 409 | 资源冲突 |
| 10006 | 429 | 请求过于频繁 |
| 10007 | 500 | 服务器内部错误 |

### 13.2 业务错误码
| 错误码 | HTTP状态码 | 说明 |
|--------|------------|------|
| 20001 | 400 | 用户已存在 |
| 20002 | 400 | 用户不存在 |
| 20003 | 400 | 密码错误 |
| 20004 | 400 | Token无效 |
| 20005 | 400 | Token已过期 |
| 30001 | 400 | 文件格式不支持 |
| 30002 | 400 | 文件大小超限 |
| 30003 | 400 | 文件上传失败 |
| 40001 | 400 | 简历不存在 |
| 40002 | 400 | 简历正在处理中 |
| 40003 | 400 | 分析配额不足 |
| 40004 | 500 | 分析服务异常 |

## 14. 限流和配额

### 14.1 API限流规则
| 用户类型 | 每分钟请求数 | 每小时请求数 | 每日请求数 |
|----------|--------------|--------------|------------|
| 免费用户 | 60 | 1000 | 5000 |
| 高级用户 | 120 | 3000 | 20000 |
| 企业用户 | 300 | 10000 | 100000 |

### 14.2 分析配额
| 用户类型 | 每日分析次数 | 每月分析次数 | 并发分析数 |
|----------|--------------|--------------|------------|
| 免费用户 | 5 | 50 | 1 |
| 高级用户 | 50 | 1000 | 3 |
| 企业用户 | 无限制 | 无限制 | 10 |

### 14.3 限流响应头
```http
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 45
X-RateLimit-Reset: 1642262400
```

## 15. WebHook回调

### 15.1 分析完成回调
```http
POST <callback_url>
Content-Type: application/json
X-Webhook-Signature: sha256=<signature>

{
    "event": "analysis.completed",
    "data": {
        "analysis_id": "analysis_12345",
        "resume_id": "resume_12345",
        "user_id": 12345,
        "status": "completed",
        "overall_score": 8.5,
        "completed_at": "2025-01-15T10:32:00Z"
    },
    "timestamp": "2025-01-15T10:32:05Z"
}
```

### 15.2 签名验证
```python
import hmac
import hashlib

def verify_webhook_signature(payload, signature, secret):
    expected_signature = hmac.new(
        secret.encode('utf-8'),
        payload.encode('utf-8'),
        hashlib.sha256
    ).hexdigest()
    return f"sha256={expected_signature}" == signature
```

## 16. SDK和代码示例

### 16.1 JavaScript SDK
```javascript
// 安装：npm install resumeoptim-sdk

import ResumeOptim from 'resumeoptim-sdk';

const client = new ResumeOptim({
    apiKey: 'your-api-key',
    baseURL: 'https://api.resumeoptim.com'
});

// 上传简历
const result = await client.resumes.upload({
    file: fileData,
    title: '我的简历',
    targetIndustry: 'technology'
});

// 开始分析
const analysis = await client.resumes.analyze(result.resume_id, {
    analysisType: 'comprehensive',
    targetPosition: '软件工程师'
});

// 获取分析结果
const analysisResult = await client.resumes.getAnalysis(result.resume_id);
console.log('分析得分:', analysisResult.overall_score);
```

### 16.2 Python SDK
```python
# 安装：pip install resumeoptim-python

from resumeoptim import ResumeOptimClient

client = ResumeOptimClient(
    api_key='your-api-key',
    base_url='https://api.resumeoptim.com'
)

# 上传简历
with open('resume.pdf', 'rb') as f:
    result = client.resumes.upload(
        file=f,
        title='我的简历',
        target_industry='technology'
    )

# 开始分析
analysis = client.resumes.analyze(
    resume_id=result['resume_id'],
    analysis_type='comprehensive',
    target_position='软件工程师'
)

# 等待分析完成
analysis_result = client.resumes.wait_for_analysis(result['resume_id'])
print(f"分析得分: {analysis_result['overall_score']}")
```

### 16.3 cURL示例
```bash
# 用户登录
curl -X POST https://api.resumeoptim.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'

# 上传文件
curl -X POST https://api.resumeoptim.com/api/v1/files/upload \
  -H "Authorization: Bearer <jwt_token>" \
  -F "file=@resume.pdf" \
  -F "title=我的简历"

# 创建简历分析
curl -X POST https://api.resumeoptim.com/api/v1/resumes \
  -H "Authorization: Bearer <jwt_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "软件工程师简历",
    "file_id": "file_12345",
    "target_industry": "technology",
    "target_position": "软件工程师"
  }'

# 开始分析
curl -X POST https://api.resumeoptim.com/api/v1/resumes/resume_12345/analyze \
  -H "Authorization: Bearer <jwt_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "analysis_type": "comprehensive",
    "target_position": "高级软件工程师"
  }'
```

## 17. 测试接口

### 17.1 测试数据生成
```http
POST /api/v1/test/generate-sample-data
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
    "type": "resume", // resume, analysis, knowledge
    "count": 5,
    "industry": "technology"
}
```

### 17.2 性能测试
```http
GET /api/v1/test/performance?endpoint=/api/v1/resumes&concurrent=10&duration=60
Authorization: Bearer <jwt_token>
```

### 17.3 清理测试数据
```http
DELETE /api/v1/test/cleanup
Authorization: Bearer <jwt_token>
```

## 18. API版本控制

### 18.1 版本策略
- URL版本控制：`/api/v1/`, `/api/v2/`
- 向后兼容：旧版本API继续可用
- 废弃通知：提前3个月通知API废弃
- 迁移指南：提供详细的版本迁移文档

### 18.2 版本差异
| 功能 | v1.0 | v1.1 | v2.0 |
|------|------|------|------|
| 基础功能 | ✅ | ✅ | ✅ |
| 批量分析 | ❌ | ✅ | ✅ |
| 实时分析 | ❌ | ❌ | ✅ |
| 多语言支持 | ❌ | ❌ | ✅ |

### 18.3 废弃API列表
| API | 废弃版本 | 替代API | 移除时间 |
|-----|----------|---------|----------|
| `/api/v1/old-upload` | v1.1 | `/api/v1/files/upload` | 2025-06-01 |

---

**文档版本**: v1.0  
**创建时间**: 2025年  
**最后更新**: 2025年  
**负责人**: API开发团队
