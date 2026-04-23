import type { VercelResponse } from "@vercel/node"
import { redisGet, redisSet, redisDel } from "../../_lib/redis.js"

// Redis キーの検証
export function validateKey(key: unknown): { valid: boolean; error?: string } {
  if (typeof key !== "string" || key.length === 0)
    return { valid: false, error: "key は空でない文字列である必要があります" }
  if (key.length > 256)
    return { valid: false, error: "key は 256 文字以内である必要があります" }
  if (!/^[a-zA-Z0-9_:-]+$/.test(key))
    return {
      valid: false,
      error:
        "key に使用できる文字は英数字・ハイフン・アンダースコア・コロンのみです",
    }
  return { valid: true }
}

export function validateValue(value: unknown): {
  valid: boolean
  error?: string
} {
  if (typeof value !== "string")
    return { valid: false, error: "value は文字列である必要があります" }
  if (Buffer.byteLength(value, "utf8") > 10 * 1024 * 1024)
    return { valid: false, error: "value のサイズが上限（10MB）を超えています" }
  return { valid: true }
}

export async function handleCreate(
  key: string,
  value: string,
  res: VercelResponse,
): Promise<VercelResponse> {
  try {
    const existing = await redisGet(key)
    if (existing !== null)
      return res
        .status(409)
        .json({ success: false, error: "すでにキーが存在します" })
    await redisSet(key, value)
    return res.json({ success: true, data: { key, value, created: true } })
  } catch (error) {
    console.error("Redis create error:", error)
    return res
      .status(500)
      .json({ success: false, error: "Internal server error" })
  }
}

export async function handleGet(
  key: string,
  res: VercelResponse,
): Promise<VercelResponse> {
  try {
    const value = await redisGet<string>(key)
    return res.json({
      success: true,
      data: { key, value, exists: value !== null },
    })
  } catch (error) {
    console.error("Redis get error:", error)
    return res
      .status(500)
      .json({ success: false, error: "Internal server error" })
  }
}

export async function handleUpdate(
  key: string,
  value: string,
  res: VercelResponse,
): Promise<VercelResponse> {
  try {
    const existing = await redisGet(key)
    if (existing === null)
      return res
        .status(404)
        .json({ success: false, error: "キーが見つかりません" })
    await redisSet(key, value)
    return res.json({ success: true, data: { key, value, updated: true } })
  } catch (error) {
    console.error("Redis update error:", error)
    return res
      .status(500)
      .json({ success: false, error: "Internal server error" })
  }
}

export async function handleDelete(
  key: string,
  res: VercelResponse,
): Promise<VercelResponse> {
  try {
    await redisDel(key)
    return res.json({ success: true, data: { key, deleted: true } })
  } catch (error) {
    console.error("Redis delete error:", error)
    return res
      .status(500)
      .json({ success: false, error: "Internal server error" })
  }
}
