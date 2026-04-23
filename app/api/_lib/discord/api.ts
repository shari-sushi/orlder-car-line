/**
 * Discord REST API 通信ヘルパー
 */
import { DISCORD_BOT_TOKEN, DISCORD_API_BASE_URL } from "../env.js"

export class DiscordApiError extends Error {
  constructor(
    public status: number,
    public statusText: string,
    public details?: unknown,
  ) {
    super(`Discord API Error: ${status} ${statusText}`)
    this.name = "DiscordApiError"
  }
}

interface MessageBody {
  content: string
  components?: unknown[]
}

/** チャンネルにメッセージを送信する */
export async function sendDiscordMessage(
  channelId: string,
  content: string,
  components?: unknown[],
): Promise<{ id: string }> {
  const body: MessageBody = { content }
  if (components?.length) body.components = components

  const res = await fetch(
    `${DISCORD_API_BASE_URL}/channels/${channelId}/messages`,
    {
      method: "POST",
      headers: {
        Authorization: `Bot ${DISCORD_BOT_TOKEN}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    },
  )
  if (!res.ok)
    throw new DiscordApiError(
      res.status,
      res.statusText,
      await res.json().catch(() => ({})),
    )
  return res.json() as Promise<{ id: string }>
}

/** メッセージを編集する */
export async function editDiscordMessage(
  channelId: string,
  messageId: string,
  content: string,
  components?: unknown[],
): Promise<void> {
  const body: MessageBody = { content }
  if (components?.length) body.components = components

  const res = await fetch(
    `${DISCORD_API_BASE_URL}/channels/${channelId}/messages/${messageId}`,
    {
      method: "PATCH",
      headers: {
        Authorization: `Bot ${DISCORD_BOT_TOKEN}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    },
  )
  if (!res.ok)
    throw new DiscordApiError(
      res.status,
      res.statusText,
      await res.json().catch(() => ({})),
    )
}
