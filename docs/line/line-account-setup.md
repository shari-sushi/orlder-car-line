# LINE 開発者アカウント作成手順

## 1. LINE ビジネス ID を作成する

1. <https://account.line.biz/signup> にアクセス
2. 「メールアドレスで登録」を選択
3. メールアドレス・パスワードを入力して送信
4. 届いた確認メールのリンクをクリックして認証完了

## 2. LINE Developers コンソールにアクセスする

1. <https://developers.line.biz/ja/> にアクセス
2. 作成した LINE ビジネス ID でログイン

## 3. プロバイダーを作成する

1. [コンソールトップ](https://developers.line.biz/console)で「プロバイダーを作成」をクリック
2. プロバイダー名（サービス名など）を入力して作成

## 4. Messaging API チャネルを作成する

※2026/04時点
> INE DevelopersコンソールからMessaging APIチャネルを直接作成することはできなくなりました。
> Messaging APIチャネルを作成するには、以下の［LINE公式アカウントを作成する］ボタンからLINE公式アカウントを作成した後、LINE Official Account Manager上でMessaging APIの利用を有効にしてください。

Messaging API チャネルを作成すると、LINE 公式アカウントが同時に開設される。
チャネル作成時に以下の 3 つの規約・ポリシーへの同意が必要。

- [LINE 公式アカウント利用規約](https://terms2.line.me/official_account_terms_jp)
- [LINE ビジネスマネージャー利用規約](https://terms2.line.me/businessmanager?lang=ja)
- [LY Corporation プライバシーポリシー](https://www.lycorp.co.jp/ja/company/privacypolicy/)

1. 作成したプロバイダーを開き「チャネルを作成」をクリック
2. 「Messaging API」を選択
3. 以下の項目を入力して作成

   | 項目            | 内容                 |
   | --------------- | -------------------- |
   | チャネル名      | Bot の名前           |
   | チャネル説明    | 簡単な説明文         |
   | 大業種 / 小業種 | 任意のカテゴリを選択 |
   | メールアドレス  | 連絡先メールアドレス |

### LINE ビジネスマネージャーについて

LINE 公式アカウントの管理・広告配信などを行う統合管理ツール。
Messaging API チャネル作成後は [LINE ビジネスマネージャー](https://www.lycbiz.com/jp/service/business-manager/) からも公式アカウントを管理できる。

## 5. 認証情報を取得する

チャネル作成後、以下の値をコピーしておく。

| 項目                 | 場所                                      |
| -------------------- | ----------------------------------------- |
| Channel Secret       | 「チャネル基本設定」タブ                  |
| Channel Access Token | 「Messaging API 設定」タブ → 最下部で発行 |

取得した値は `.env` に設定する。

```sh
LINE_CHANNEL_SECRET=<Channel Secret>
LINE_CHANNEL_ACCESS_TOKEN=<Channel Access Token>
```

## 参考

- [LINE Developers コンソール](https://developers.line.biz/ja/)
- [Messaging API ドキュメント](https://developers.line.biz/ja/docs/messaging-api/)
