import type { VercelRequest, VercelResponse } from "@vercel/node"
import { redisDel } from "./_lib/redis.js"

function getToken(cookieHeader: string | undefined): string | null {
  if (!cookieHeader) return null
  const m = cookieHeader.match(/rental_session=([^;]+)/)
  return m ? m[1] : null
}

export default async function handler(req: VercelRequest, res: VercelResponse) {
  if (req.method !== "POST")
    return res.status(405).json({ error: "Method not allowed" })

  const token = getToken(req.headers.cookie)
  if (token) {
    await redisDel(`rental:session:${token}`)
  }

  res.setHeader(
    "Set-Cookie",
    "rental_session=; HttpOnly; Path=/; Max-Age=0; SameSite=Lax",
  )
  return res.json({ ok: true })
}
