import Cookies from 'js-cookie'

const TOKEN_KEY = 'resume_optim_token'
const TOKEN_EXPIRE_KEY = 'resume_optim_token_expire'

/**
 * 获取token
 */
export function getToken(): string | undefined {
  return Cookies.get(TOKEN_KEY)
}

/**
 * 设置token
 */
export function setToken(token: string, expires?: number): void {
  const options: Cookies.CookieAttributes = {
    secure: process.env.NODE_ENV === 'production',
    sameSite: 'strict'
  }
  
  if (expires) {
    options.expires = new Date(expires * 1000)
    Cookies.set(TOKEN_EXPIRE_KEY, expires.toString(), options)
  }
  
  Cookies.set(TOKEN_KEY, token, options)
}

/**
 * 移除token
 */
export function removeToken(): void {
  Cookies.remove(TOKEN_KEY)
  Cookies.remove(TOKEN_EXPIRE_KEY)
}

/**
 * 检查token是否过期
 */
export function isTokenExpired(): boolean {
  const expireTime = Cookies.get(TOKEN_EXPIRE_KEY)
  if (!expireTime) return true
  
  const now = Math.floor(Date.now() / 1000)
  return now >= parseInt(expireTime)
}

/**
 * 格式化Authorization header
 */
export function getAuthHeader(): Record<string, string> {
  const token = getToken()
  return token ? { Authorization: `Bearer ${token}` } : {}
}
