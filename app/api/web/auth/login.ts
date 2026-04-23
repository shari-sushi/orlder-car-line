import type { VercelRequest, VercelResponse } from "@vercel/node"
import { createSession } from "../../_lib/session.js"
import { AUTH_USER_PASS } from "../../_lib/env.js"

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

export default async function handler(
  req: VercelRequest,
  res: VercelResponse,
): Promise<VercelResponse> {
  if (req.method === "OPTIONS") return res.status(204).end()
  if (req.method !== "POST")
    return res.status(405).json({ error: "Method not allowed" })

  try {
    const { username, password } = (req.body ?? {}) as {
      username?: string
      password?: string
    }
    if (!username || !password)
      return res.status(400).json({ error: "username と password が必要です" })

    const users = parseUsers()
    if (users[username] !== password)
      return res
        .status(401)
        .json({ error: "ユーザー名またはパスワードが正しくありません" })

    const token = await createSession(username)
    return res.json({ token, username })
  } catch (err) {
    console.error("[web/auth/login]", err)
    return res.status(500).json({ error: String(err) })
  }
}
