/**
 * サーバー専用環境変数の一元管理
 * 本番/プレビュー環境では未設定時に例外、開発環境では警告のみ
 */

function getRequiredEnv(key: string): string {
  const value = process.env[key]
  if (!value) {
    const isProd =
      process.env.VERCEL_ENV === "production" ||
      process.env.VERCEL_ENV === "preview"
    if (isProd) throw new Error(`環境変数 ${key} が設定されていません`)
    console.warn(`[dev] 環境変数 ${key} が未設定です`)
    return ""
  }
  return value
}

function getOptionalEnv(key: string, defaultValue: string): string {
  return process.env[key] ?? defaultValue
}

// デプロイ環境
export const ENV = getOptionalEnv("VERCEL_ENV", "development") as
  | "production"
  | "preview"
  | "development"

// Discord
export const DISCORD_PUBLIC_KEY = getRequiredEnv("DISCORD_PUBLIC_KEY")
export const DISCORD_APPLICATION_ID = getRequiredEnv("DISCORD_APPLICATION_ID")
export const DISCORD_BOT_TOKEN = getRequiredEnv("DISCORD_BOT_TOKEN")
export const DISCORD_COMMAND_GUILD_ID = getOptionalEnv(
  "DISCORD_COMMAND_GUILD_ID",
  "",
)
export const DISCORD_API_BASE_URL = "https://discord.com/api/v10"

// 認証
export const AUTH_USER_PASS = getRequiredEnv("AUTH_USER_PASS") // "user1:pass1,user2:pass2"

// データベース
export const REDIS_URL = getOptionalEnv(
  "STORAGE_URL",
  getOptionalEnv("REDIS_URL", "redis://localhost:6379"),
)
