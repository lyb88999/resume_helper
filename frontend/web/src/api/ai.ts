import request from '@/utils/request'
import type {
  ChatRequest,
  ChatResponse,
  RetrieveKnowledgeRequest,
  RetrieveKnowledgeResponse,
  GenerateSuggestionsRequest,
  GenerateSuggestionsResponse,
  AIHealthResponse
} from '@/types/ai'

/**
 * AI智能问答
 */
export function chatWithAI(data: ChatRequest): Promise<ChatResponse> {
  return request.post('/v1/ai/chat', data)
}

/**
 * 知识检索
 */
export function retrieveKnowledge(data: RetrieveKnowledgeRequest): Promise<RetrieveKnowledgeResponse> {
  return request.post('/v1/ai/knowledge/retrieve', data)
}

/**
 * 生成优化建议
 */
export function generateSuggestions(data: GenerateSuggestionsRequest): Promise<GenerateSuggestionsResponse> {
  return request.post('/v1/ai/suggestions/generate', data)
}

/**
 * AI服务健康检查
 */
export function checkAIHealth(): Promise<AIHealthResponse> {
  return request.get('/v1/ai/health')
}

/**
 * 简历分析
 */
export function analyzeResumeWithAI(data: {
  resumeId: string
  content?: string
  fileType?: string
  targetPosition?: string
  options?: {
    enableCompleteness?: boolean
    enableClarity?: boolean
    enableKeyword?: boolean
    enableFormat?: boolean
    enableQuantification?: boolean
  }
}): Promise<{
  analysisId: string
  status: string
  message: string
}> {
  return request.post('/v1/ai/resume/analyze', data)
}

/**
 * 获取简历分析结果
 */
export function getResumeAnalysisResult(analysisId: string): Promise<{
  data: {
    sections: Record<string, any>
    suggestions: any[]
    scores: {
      overallScore: number
      completenessScore: number
      clarityScore: number
      keywordScore: number
      formatScore: number
      quantificationScore: number
      dimensionScores: Record<string, number>
    }
    improvements: any[]
    summary: string
  }
  status: string
  message: string
}> {
  return request.get(`/v1/ai/resume/analysis/${analysisId}/result`)
}
