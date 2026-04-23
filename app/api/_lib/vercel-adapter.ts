import type { IncomingMessage, ServerResponse } from "http"
import type { VercelRequest, VercelResponse } from "@vercel/node"

export function adaptVercel(
  handler: (req: VercelRequest, res: VercelResponse) => Promise<unknown>,
) {
  return async (req: IncomingMessage, res: ServerResponse) => {
    // --- body のパース ---
    const chunks: Buffer[] = []
    req.on("data", (chunk: Buffer) => chunks.push(chunk))
    await new Promise<void>((resolve) => req.on("end", resolve))
    const raw = Buffer.concat(chunks).toString()
    const body = raw ? JSON.parse(raw) : {}

    // --- VercelRequest の拡張 ---
    const vReq = Object.assign(req, { body }) as VercelRequest

    // --- VercelResponse の拡張 ---
    const vRes = Object.assign(res, {
      json(data: unknown) {
        res.setHeader("Content-Type", "application/json")
        res.end(JSON.stringify(data))
        return vRes
      },
      status(code: number) {
        res.statusCode = code
        return vRes
      },
    }) as unknown as VercelResponse

    await handler(vReq, vRes)
  }
}
