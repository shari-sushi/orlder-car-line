import type { VercelRequest, VercelResponse } from "@vercel/node"
import { redisGet } from "./_lib/redis.js"

function getToken(cookieHeader: string | undefined): string | null {
  if (!cookieHeader) return null
  const m = cookieHeader.match(/rental_session=([^;]+)/)
  return m ? m[1] : null
}

export default async function handler(req: VercelRequest, res: VercelResponse) {
  if (req.method !== "GET")
    return res.status(405).json({ error: "Method not allowed" })

  const token = getToken(req.headers.cookie)
  if (!token) return res.status(401).json({ error: "Not logged in" })

  const username = await redisGet<string>(`rental:session:${token}`)
  if (!username) return res.status(401).json({ error: "Session expired" })

  return res.json({ username })
}
