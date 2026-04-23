import { createClient } from "redis"

// Vercel が作る STORAGE_URL を優先、ローカルは REDIS_URL or デフォルト
const REDIS_URL =
  process.env.STORAGE_URL ?? process.env.REDIS_URL ?? "redis://localhost:6379"

let client: ReturnType<typeof createClient> | null = null

async function getClient() {
  if (client?.isReady) return client

  client = createClient({
    url: REDIS_URL,
    socket: {
      connectTimeout: 5000,
      reconnectStrategy: (retries) => {
        if (retries >= 3) return new Error("Redis reconnect limit reached")
        return Math.min(retries * 200, 1000)
      },
    },
  })
  client.on("error", (err) => console.error("[Redis]", err))
  await client.connect()
  return client
}

export async function redisGet<T>(key: string): Promise<T | null> {
  const c = await getClient()
  const val = await c.get(key)
  if (val === null) return null
  try {
    return JSON.parse(val.toString()) as T
  } catch {
    return val.toString() as unknown as T
  }
}

export async function redisDel(key: string): Promise<void> {
  const c = await getClient()
  await c.del(key)
}

export async function redisSet(
  key: string,
  value: unknown,
  expiresInSeconds?: number,
): Promise<void> {
  const c = await getClient()
  const serialized = typeof value === "string" ? value : JSON.stringify(value)
  if (expiresInSeconds) {
    await c.setEx(key, expiresInSeconds, serialized)
  } else {
    await c.set(key, serialized)
  }
}
