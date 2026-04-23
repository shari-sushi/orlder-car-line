/**
 * Discord スラッシュコマンド名 & コンポーネント custom_id の定数定義
 * コマンド追加時はここに追記する
 */

export const COMMANDS = {
  // サンプル
  SAMPLE: {
    HELLO: "sample-hello",
  },
  // 開発者用
  DEV: {
    ECHO: "dev-echo",
  },
} as const

export const CLIENT_ACTIONS = {
  SAMPLE: {
    BUTTON: "sample-button",
  },
} as const
