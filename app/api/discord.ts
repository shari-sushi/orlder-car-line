/**
 * Discord Interactions エンドポイント
 * POST /api/discord
 *
 * Discord Developer Portal の Interactions Endpoint URL に設定する
 */
import type { VercelRequest, VercelResponse } from "@vercel/node"
import { verifyKey } from "discord-interactions"
import { InteractionResponseType, InteractionType } from "discord-interactions"
import { DISCORD_PUBLIC_KEY } from "./_lib/env.js"
import { helloCommand } from "./discord/commands/sample.js"
import { CLIENT_ACTIONS, COMMANDS } from "./discord/_util/commands.js"

export default async function handler(
  req: VercelRequest,
  res: VercelResponse,
): Promise<VercelResponse> {
  if (req.method !== "POST")
    return res.status(405).json({ error: "Method not allowed" })

  // 署名検証
  const signature = req.headers["x-signature-ed25519"] as string | undefined
  const timestamp = req.headers["x-signature-timestamp"] as string | undefined
  const rawBody = JSON.stringify(req.body)

  if (!signature || !timestamp)
    return res.status(401).json({ error: "Missing headers" })

  const isValid = verifyKey(rawBody, signature, timestamp, DISCORD_PUBLIC_KEY)
  if (!isValid) return res.status(401).json({ error: "Invalid signature" })

  const interaction = req.body as {
    type: number
    data?: { name?: string; custom_id?: string }
  }

  // PING（疎通確認）
  if (interaction.type === InteractionType.PING) {
    return res.json({ type: InteractionResponseType.PONG })
  }

  // スラッシュコマンド
  if (interaction.type === InteractionType.APPLICATION_COMMAND) {
    const commandName = interaction.data?.name ?? ""
    console.log("command:", commandName)

    switch (commandName) {
      case COMMANDS.SAMPLE.HELLO:
        return res.json(helloCommand())
      default:
        return res.json({
          type: InteractionResponseType.CHANNEL_MESSAGE_WITH_SOURCE,
          data: { content: `不明なコマンド: ${commandName}` },
        })
    }
  }

  // ボタン・セレクトメニュー
  if (interaction.type === InteractionType.MESSAGE_COMPONENT) {
    const customId = interaction.data?.custom_id ?? ""
    const [actionId] = customId.split("?")
    console.log("component action:", actionId)

    switch (actionId) {
      case CLIENT_ACTIONS.SAMPLE.BUTTON:
        return res.json({
          type: InteractionResponseType.CHANNEL_MESSAGE_WITH_SOURCE,
          data: { content: "ボタンが押されました！" },
        })
      default:
        return res.status(400).json({ error: `Unknown action: ${actionId}` })
    }
  }

  return res.status(400).json({ error: "Unknown interaction type" })
}
