import type { VercelRequest, VercelResponse } from "@vercel/node"
import { randomUUID } from "node:crypto"
import { redisSet } from "./_lib/redis.js"

const SESSION_TTL = 60 * 60 * 24 * 7 // 7日

function parseUsers(): Record<string, string> {
  const raw = process.env.AUTH_USER_PASS ?? ""
  const result: Record<string, string> = {}
  for (const entry of raw.split(",")) {
    const colonIdx = entry.indexOf(":")
    if (colonIdx === -1) continue
    const user = entry.slice(0, colonIdx).trim()
    const pass = entry.slice(colonIdx + 1).trim()
    if (user) result[user] = pass
  }
  return result
}

export default async function handler(req: VercelRequest, res: VercelResponse) {
  res.setHeader("Access-Control-Allow-Methods", "POST, OPTIONS")
  res.setHeader("Access-Control-Allow-Headers", "Content-Type")
  if (req.method === "OPTIONS") return res.status(204).end()
  if (req.method !== "POST")
    return res.status(405).json({ error: "Method not allowed" })

  try {
    const body = req.body as
      | { username?: string; password?: string }
      | undefined
    const { username, password } = body ?? {}
    const users = parseUsers()

    if (!username || !password || users[username] !== password) {
      return res
        .status(401)
        .json({ error: "ユーザー名またはパスワードが正しくありません" })
    }

    const token = randomUUID()
    await redisSet(`rental:session:${token}`, username, SESSION_TTL)

    res.setHeader(
      "Set-Cookie",
      `rental_session=${token}; HttpOnly; Path=/; Max-Age=${SESSION_TTL}; SameSite=Lax`,
    )
    return res.json({ username })
  } catch (err) {
    console.error("[login] error:", err)
    return res.status(500).json({ error: String(err) })
  }
}
