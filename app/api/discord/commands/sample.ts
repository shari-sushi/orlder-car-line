import { InteractionResponseType } from "discord-interactions"

/**
 * /sample-hello コマンドのハンドラ
 */
export function helloCommand(): { type: number; data: { content: string } } {
  return {
    type: InteractionResponseType.CHANNEL_MESSAGE_WITH_SOURCE,
    data: { content: "こんにちは！" },
  }
}
