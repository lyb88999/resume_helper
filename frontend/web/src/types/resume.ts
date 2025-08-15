// 简历状态枚举
export type ResumeStatus = 'pending' | 'processing' | 'completed' | 'failed'

// 文件类型枚举
export type FileType = 'pdf' | 'doc' | 'docx' | 'md' | 'txt'

// 简历基础信息
export interface Resume {
  id: number
  name: string
  status: ResumeStatus
  uploadTime: string
  fileType: FileType
  fileSize: number
  filePath: string
  score?: number
  analysisId?: string
  userId: number
  createdAt: string
  updatedAt: string
}

// 上传简历数据
export interface UploadResumeData {
  file: File
  targetPosition?: string
  industry?: string
}

// 简历列表请求参数
export interface ResumeListRequest {
  page: number
  pageSize: number
  status?: ResumeStatus
  userId?: number
}

// 简历列表响应
export interface ResumeListResponse {
  data: Resume[]
  total: number
  page: number
  pageSize: number
}

// 简历上传响应
export interface UploadResumeResponse {
  data: Resume
  message: string
}

// 简历解析请求
export interface ParseResumeRequest {
  resumeId: number
}

// 简历解析响应
export interface ParseResumeResponse {
  taskId: string
  message: string
}

// 获取解析结果响应
export interface ParseResultResponse {
  data: Resume
  message: string
}

// 简历分析请求
export interface AnalyzeResumeRequest {
  resumeId: number
  targetPosition: string
  industry: string
  options?: AnalysisOptions
}

// 分析选项
export interface AnalysisOptions {
  enableCompleteness?: boolean
  enableClarity?: boolean
  enableKeyword?: boolean
  enableFormat?: boolean
  enableQuantification?: boolean
}

// 简历分析响应
export interface AnalyzeResumeResponse {
  analysisId: string
  status: string
  message: string
}

// 获取分析结果响应
export interface AnalysisResultResponse {
  data: AnalysisResult
  message: string
}

// 分析结果
export interface AnalysisResult {
  id: string
  resumeId: number
  status: string
  scores: ScoreBreakdown
  suggestions: Suggestion[]
  improvements: Improvement[]
  summary: string
  analyzedAt: string
}

// 评分详情
export interface ScoreBreakdown {
  overallScore: number
  completenessScore: number
  clarityScore: number
  keywordScore: number
  formatScore: number
  quantificationScore: number
  dimensionScores: Record<string, number>
}

// 建议
export interface Suggestion {
  id: string
  type: string
  title: string
  description: string
  priority: string
  section: string
  action: string
  examples?: string[]
}

// 改进建议
export interface Improvement {
  type: string
  description: string
  priority: string
  section: string
  before: string
  after: string
  examples?: string[]
}

// 删除简历响应
export interface DeleteResumeResponse {
  message: string
}
