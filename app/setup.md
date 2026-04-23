# セットアップ手順

## 1. Google Maps API キーの取得

### Google Cloud Console でプロジェクト作成

1. https://console.cloud.google.com/ にアクセス（Google アカウントでログイン）
2. 画面上部のプロジェクト選択 → **「新しいプロジェクト」**
3. プロジェクト名を入力（例: `rental-map`）→ **作成**

### Maps JavaScript API を有効化

1. 左メニュー → **「APIとサービス」→「ライブラリ」**
2. 検索欄に `Maps JavaScript API` と入力 → 選択
3. **「有効にする」** をクリック

> 初回利用時は請求先アカウントの設定が必要です。クレジットカードを登録しますが、月 $200 の無料枠があり個人利用では基本的に無料です。

### API キーを発行

1. 左メニュー → **「APIとサービス」→「認証情報」**
2. **「認証情報を作成」→「API キー」**
3. 作成されたキー（`AIzaSy...`）をコピーして保存

### API キーの制限（推奨）

キーが漏洩しても悪用されないよう制限を設定します。

1. 作成したキーをクリック → **「キーを制限」**
2. **「アプリケーションの制限」→「HTTP リファラー（ウェブサイト）」** を選択
3. 以下を追加:
   - `localhost:5173/*`（ローカル開発用）
   - `https://your-app.vercel.app/*`（Vercel の URL に置き換える）
4. **「API の制限」→「キーを制限」** を選択
5. **「Maps JavaScript API」** を選択 → 保存

---

## 2. ローカル開発のセットアップ

### .env.local にキーを設定

`app/.env.local` ファイルを作成（なければ新規作成）して以下を追記:

```sh
VITE_GOOGLE_MAPS_API_KEY=AIzaSy...ここに取得したキーを貼る
```

> `.env.local` は `.gitignore` に含まれているためリポジトリには含まれません。

### 起動確認

```bash
cd app
npm run dev
# → http://localhost:5173/
```

TOP ページの「地図で見る」ボタン → `/map` ページで地図と全物件ピンが表示されればOK。

---

## 3. Vercel へのデプロイ

### Vercel アカウントとプロジェクト

1. https://vercel.com/ にアクセス → GitHub アカウントでログイン
2. **「Add New → Project」**
3. このリポジトリを選択 → **「Import」**
4. 設定:
   - **Framework Preset**: `Vite`（自動検出されるはず）
   - **Root Directory**: `app`
5. **「Deploy」**

### Vercel に環境変数を登録

デプロイ後、地図を表示するには API キーを Vercel に設定する必要があります。

1. Vercel ダッシュボード → プロジェクトを選択
2. **「Settings」→「Environment Variables」**
3. 以下を追加:

| Name                       | Value       | Environment                        |
| -------------------------- | ----------- | ---------------------------------- |
| `VITE_GOOGLE_MAPS_API_KEY` | `AIzaSy...` | Production / Preview / Development |

1. **「Save」**
1. **「Deployments」→「Redeploy」** で再デプロイして反映させる

> `VITE_` プレフィックスは Vite がビルド時にクライアントへ埋め込むために必要です。`VITE_` がないと地図が表示されません。

### 再デプロイ不要な変更

`賃貸情報.json` の編集や画像追加は `git push` するだけで Vercel が自動的に再ビルド＆デプロイします。

---

## 環境変数まとめ

| 変数名                     | 用途                          | 設定場所                                          |
| -------------------------- | ----------------------------- | ------------------------------------------------- |
| `VITE_GOOGLE_MAPS_API_KEY` | Google Maps 表示              | `.env.local`（ローカル）、Vercel 環境変数（本番） |
| `REDIS_URL`                | お気に入り・NG フラグの永続化 | `.env.local`（ローカル）、Vercel 環境変数（本番） |

`REDIS_URL` が未設定でもアプリは動作します（フラグが localStorage に保存される）。
