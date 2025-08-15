import { defineStore } from 'pinia'
import { ref, computed, readonly } from 'vue'
import { ElMessage } from 'element-plus'
import * as resumeApi from '@/api/resume'
import type { Resume, AnalysisResult, UploadResumeData } from '@/types/resume'

export const useResumeStore = defineStore('resume', () => {
  // State
  const resumes = ref<Resume[]>([])
  const currentResume = ref<Resume | null>(null)
  const analysisResult = ref<AnalysisResult | null>(null)
  const loading = ref(false)
  const uploading = ref(false)
  const analyzing = ref(false)
  const uploadProgress = ref(0)

  // Getters
  const resumeCount = computed(() => resumes.value.length)
  const hasCurrentResume = computed(() => !!currentResume.value)
  const hasAnalysisResult = computed(() => !!analysisResult.value)

  // Actions
  
  // 获取简历列表
  const fetchResumes = async (page = 1, pageSize = 10) => {
    loading.value = true
    try {
      const response = await resumeApi.getResumeList({ page, pageSize })
      resumes.value = response.data
      return response
    } catch (error: any) {
      ElMessage.error(error.message || '获取简历列表失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  // 上传简历
  const uploadResume = async (uploadData: UploadResumeData, onProgress?: (progress: number) => void) => {
    uploading.value = true
    uploadProgress.value = 0
    
    try {
      const response = await resumeApi.uploadResume(uploadData, (progress) => {
        uploadProgress.value = progress
        onProgress?.(progress)
      })
      
      // 添加到简历列表
      resumes.value.unshift(response.data)
      currentResume.value = response.data
      
      ElMessage.success('简历上传成功')
      return response.data
    } catch (error: any) {
      ElMessage.error(error.message || '简历上传失败')
      throw error
    } finally {
      uploading.value = false
      uploadProgress.value = 0
    }
  }

    // 解析简历
  const parseResume = async (resumeId: number) => {
    loading.value = true
    try {
      const response = await resumeApi.parseResume({ resumeId })
      
      // 更新当前简历状态
      if (currentResume.value?.id === resumeId) {
        currentResume.value = { ...currentResume.value, status: 'processing' }
      }
      
      // 更新列表中的简历状态
      const index = resumes.value.findIndex(r => r.id === resumeId)
      if (index !== -1) {
        resumes.value[index] = { ...resumes.value[index], status: 'processing' }
      }
      
      return response.taskId
    } catch (error: any) {
      ElMessage.error(error.message || '简历解析失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  // 获取解析结果
  const getParseResult = async (resumeId: number) => {
    try {
      const response = await resumeApi.getParseResult(resumeId)
      const resume = response.data
      
      // 更新当前简历
      if (currentResume.value?.id === resumeId) {
        currentResume.value = resume
      }
      
      // 更新列表中的简历
      const index = resumes.value.findIndex(r => r.id === resumeId)
      if (index !== -1) {
        resumes.value[index] = resume
      }
      
      return resume
    } catch (error: any) {
      ElMessage.error(error.message || '获取解析结果失败')
      throw error
    }
  }

  // 分析简历
  const analyzeResume = async (resumeId: number, targetPosition: string, industry: string) => {
    analyzing.value = true
    try {
      const response = await resumeApi.analyzeResume({
        resumeId,
        targetPosition,
        industry
      })
      
      ElMessage.success('分析任务已提交，请稍候查看结果')
      return response.analysisId
    } catch (error: any) {
      ElMessage.error(error.message || '简历分析失败')
      throw error
    } finally {
      analyzing.value = false
    }
  }

  // 获取分析结果
  const getAnalysisResult = async (analysisId: number) => {
    try {
      const response = await resumeApi.getAnalysisResult(analysisId)
      analysisResult.value = response.data
      return response.data
    } catch (error: any) {
      ElMessage.error(error.message || '获取分析结果失败')
      throw error
    }
  }

  // 删除简历
  const deleteResume = async (resumeId: number) => {
    try {
      await resumeApi.deleteResume(resumeId)
      
      // 从列表中移除
      const index = resumes.value.findIndex(r => r.id === resumeId)
      if (index !== -1) {
        resumes.value.splice(index, 1)
      }
      
      // 如果删除的是当前简历，清空当前简历
      if (currentResume.value?.id === resumeId) {
        currentResume.value = null
        analysisResult.value = null
      }
      
      ElMessage.success('简历删除成功')
    } catch (error: any) {
      ElMessage.error(error.message || '简历删除失败')
      throw error
    }
  }

  // 设置当前简历
  const setCurrentResume = (resume: Resume | null) => {
    currentResume.value = resume
    // 清空之前的分析结果
    analysisResult.value = null
  }

  // 清空分析结果
  const clearAnalysisResult = () => {
    analysisResult.value = null
  }

  // 重置store状态
  const resetStore = () => {
    resumes.value = []
    currentResume.value = null
    analysisResult.value = null
    loading.value = false
    uploading.value = false
    analyzing.value = false
    uploadProgress.value = 0
  }

  return {
    // State
    resumes: readonly(resumes),
    currentResume: readonly(currentResume),
    analysisResult: readonly(analysisResult),
    loading: readonly(loading),
    uploading: readonly(uploading),
    analyzing: readonly(analyzing),
    uploadProgress: readonly(uploadProgress),

    // Getters
    resumeCount,
    hasCurrentResume,
    hasAnalysisResult,

    // Actions
    fetchResumes,
    uploadResume,
    parseResume,
    getParseResult,
    analyzeResume,
    getAnalysisResult,
    deleteResume,
    setCurrentResume,
    clearAnalysisResult,
    resetStore
  }
})
