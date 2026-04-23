import { validateSession } from "./session.js"
import { AUTH_USER_PASS } from "./env.js"

interface AuthResult {
  valid: boolean
  username?: string
  error?: string
}

function parseUsers(): Record<string, string> {
  const users: Record<string, string> = {}
  for (const entry of AUTH_USER_PASS.split(",")) {
    const idx = entry.indexOf(":")
    if (idx === -1) continue
    const user = entry.slice(0, idx).trim()
    const pass = entry.slice(idx + 1).trim()
    if (user) users[user] = pass
  }
  return users
}

/**
 * Authorization ヘッダーを検証する
 * Bearer（セッショントークン）と Basic（GAS/curl 向け）の両方をサポート
 */
export async function validateAuthHeader(
  authHeader: string | null,
): Promise<AuthResult> {
  if (!authHeader) return { valid: false, error: "認証ヘッダーが必要です" }

  if (authHeader.startsWith("Bearer ")) {
    const token = authHeader.slice(7)
    const session = await validateSession(token)
    if (!session) return { valid: false, error: "無効な認証トークンです" }
    return { valid: true, username: session.username }
  }

  if (authHeader.startsWith("Basic ")) {
    let credentials: string
    try {
      credentials = Buffer.from(authHeader.slice(6), "base64").toString("utf-8")
    } catch {
      return { valid: false, error: "不正な Basic 認証形式です" }
    }
    const [username, password] = credentials.split(":")
    if (!username || !password)
      return { valid: false, error: "ユーザー名とパスワードが必要です" }
    const users = parseUsers()
    if (users[username] !== password)
      return { valid: false, error: "無効なユーザー名またはパスワードです" }
    return { valid: true, username }
  }

  return { valid: false, error: "Bearer または Basic を使用してください" }
}
