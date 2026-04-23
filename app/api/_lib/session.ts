import { randomBytes } from "node:crypto"
import { redisGet, redisSet, redisDel } from "./redis.js"

const SESSION_PREFIX = "session:"
const SESSION_TTL = 60 * 60 * 24 * 7 // 7日

export interface SessionData {
  username: string
  createdAt: number
  lastAccessedAt: number
}

export function generateToken(): string {
  return randomBytes(32).toString("hex")
}

export async function createSession(username: string): Promise<string> {
  const token = generateToken()
  const data: SessionData = {
    username,
    createdAt: Date.now(),
    lastAccessedAt: Date.now(),
  }
  await redisSet(`${SESSION_PREFIX}${token}`, data, SESSION_TTL)
  return token
}

export async function validateSession(
  token: string,
): Promise<SessionData | null> {
  if (!token) return null
  const data = await redisGet<SessionData>(`${SESSION_PREFIX}${token}`)
  if (!data) return null

  // 最終アクセス時刻を更新してセッションを延長
  data.lastAccessedAt = Date.now()
  await redisSet(`${SESSION_PREFIX}${token}`, data, SESSION_TTL)
  return data
}

export async function deleteSession(token: string): Promise<void> {
  await redisDel(`${SESSION_PREFIX}${token}`)
}
