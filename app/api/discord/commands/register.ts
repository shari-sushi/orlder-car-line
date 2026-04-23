/**
 * Discord スラッシュコマンドを登録するスクリプト
 * 実行: npx tsx api/_discord/commands/register.ts
 *
 * - 本番/プレビュー: グローバルコマンド（反映に最大1時間）
 * - DISCORD_COMMAND_GUILD_ID 設定時: ギルドコマンド（即時反映）
 */
import "dotenv/config"
import {
  type RESTPostAPIApplicationCommandsJSONBody,
  ApplicationCommandOptionType,
} from "discord-api-types/v10"
import {
  DISCORD_BOT_TOKEN,
  DISCORD_APPLICATION_ID,
  DISCORD_COMMAND_GUILD_ID,
  ENV,
} from "../../_lib/env.js"
import { COMMANDS } from "../_util/commands.js"

if (ENV !== "production" && ENV !== "preview") {
  console.log("開発環境ではコマンド登録をスキップします:", ENV)
  process.exit(0)
}

const commandDefs: RESTPostAPIApplicationCommandsJSONBody[] = [
  {
    name: COMMANDS.SAMPLE.HELLO,
    description: "サンプル: こんにちはと返します",
    options: [],
  },
  {
    name: COMMANDS.DEV.ECHO,
    description: "入力したテキストをそのまま返します",
    options: [
      {
        name: "text",
        description: "送信するテキスト",
        type: ApplicationCommandOptionType.String,
        required: true,
      },
    ],
  },
]

const base = `https://discord.com/api/v10/applications/${DISCORD_APPLICATION_ID}`
const url = DISCORD_COMMAND_GUILD_ID
  ? `${base}/guilds/${DISCORD_COMMAND_GUILD_ID}/commands`
  : `${base}/commands`

console.log(
  `コマンドを登録します（${DISCORD_COMMAND_GUILD_ID ? "ギルド（即時）" : "グローバル（最大1時間）"}）`,
)

const res = await fetch(url, {
  method: "PUT",
  headers: {
    Authorization: `Bot ${DISCORD_BOT_TOKEN}`,
    "Content-Type": "application/json",
  },
  body: JSON.stringify(commandDefs),
})

const data = await res.json()
if (!res.ok) {
  console.error("登録失敗:", JSON.stringify(data, null, 2))
  process.exit(1)
}
console.log("登録完了:", data)
