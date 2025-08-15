// AI聊天请求
export interface ChatRequest {
  sessionId: string
  message: string
  context?: string
  options?: ChatOptions
}

// 聊天选项
export interface ChatOptions {
  useResumeContext?: boolean
  useKnowledgeBase?: boolean
  language?: string
}

// AI聊天响应
export interface ChatResponse {
  response: string
  sessionId: string
  sources?: string[]
  status: string
  message: string
}

// 知识检索请求
export interface RetrieveKnowledgeRequest {
  query: string
  topK?: number
  similarityThreshold?: number
  filters?: Record<string, string>
}

// 知识检索响应
export interface RetrieveKnowledgeResponse {
  items: KnowledgeItem[]
  status: string
  message: string
}

// 知识项
export interface KnowledgeItem {
  id: string
  title: string
  content: string
  score: number
  metadata?: Record<string, string>
}

// 生成建议请求
export interface GenerateSuggestionsRequest {
  analysisId: string
  targetPosition?: string
  industry?: string
  analysisResult?: any
  options?: SuggestionOptions
}

// 建议选项
export interface SuggestionOptions {
  maxSuggestions?: number
  focusArea?: string
  experienceLevel?: string
}

// 生成建议响应
export interface GenerateSuggestionsResponse {
  suggestions: Suggestion[]
  reasoning?: string
  status: string
  message: string
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

// AI服务健康检查响应
export interface AIHealthResponse {
  status: string
  version: string
  components: Record<string, string>
}
