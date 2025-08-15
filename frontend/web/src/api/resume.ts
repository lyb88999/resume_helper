import request from '@/utils/request'
import type {
  Resume,
  ResumeListRequest,
  ResumeListResponse,
  UploadResumeData,
  UploadResumeResponse,
  ParseResumeRequest,
  ParseResumeResponse,
  ParseResultResponse,
  AnalyzeResumeRequest,
  AnalyzeResumeResponse,
  AnalysisResultResponse,
  DeleteResumeResponse
} from '@/types/resume'

/**
 * 获取简历列表
 */
export function getResumeList(params: ResumeListRequest): Promise<ResumeListResponse> {
  return request.get('/v1/resume/list', { params })
}

/**
 * 上传简历
 */
export function uploadResume(
  data: UploadResumeData,
  onProgress?: (progress: number) => void
): Promise<UploadResumeResponse> {
  const formData = new FormData()
  formData.append('file', data.file)
  if (data.targetPosition) {
    formData.append('targetPosition', data.targetPosition)
  }
  if (data.industry) {
    formData.append('industry', data.industry)
  }

  return request.post('/v1/resume/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    onUploadProgress: (progressEvent) => {
      if (progressEvent.total) {
        const progress = Math.round((progressEvent.loaded * 100) / progressEvent.total)
        onProgress?.(progress)
      }
    }
  })
}

/**
 * 解析简历
 */
export function parseResume(data: ParseResumeRequest): Promise<ParseResumeResponse> {
  return request.post('/v1/resume/parse', data)
}

/**
 * 获取解析结果
 */
export function getParseResult(resumeId: number): Promise<ParseResultResponse> {
  return request.get(`/v1/resume/${resumeId}/parse-result`)
}

/**
 * 分析简历
 */
export function analyzeResume(data: AnalyzeResumeRequest): Promise<AnalyzeResumeResponse> {
  return request.post('/v1/resume/analyze', data)
}

/**
 * 获取分析结果
 */
export function getAnalysisResult(analysisId: string): Promise<AnalysisResultResponse> {
  return request.get(`/v1/ai/analysis/${analysisId}/result`)
}

/**
 * 删除简历
 */
export function deleteResume(resumeId: number): Promise<DeleteResumeResponse> {
  return request.delete(`/v1/resume/${resumeId}`)
}

/**
 * 获取简历详情
 */
export function getResumeDetail(resumeId: number): Promise<{ data: Resume }> {
  return request.get(`/v1/resume/${resumeId}`)
}

/**
 * 更新简历信息
 */
export function updateResume(resumeId: number, data: Partial<Resume>): Promise<{ data: Resume }> {
  return request.put(`/v1/resume/${resumeId}`, data)
}
