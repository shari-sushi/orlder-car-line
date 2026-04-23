import type { VercelRequest, VercelResponse } from "@vercel/node"
import { deleteSession } from "../../_lib/session.js"

export default async function handler(
  req: VercelRequest,
  res: VercelResponse,
): Promise<VercelResponse> {
  if (req.method === "OPTIONS") return res.status(204).end()
  if (req.method !== "POST")
    return res.status(405).json({ error: "Method not allowed" })

  try {
    const token = (req.headers["authorization"] ?? "").replace("Bearer ", "")
    if (token) await deleteSession(token)
    return res.json({ success: true })
  } catch (err) {
    console.error("[web/auth/logout]", err)
    return res.status(500).json({ error: String(err) })
  }
}
