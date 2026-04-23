# コーディング規約

## TypeScript / React

### catch 内では必ず console.error を出力する

エラーを握りつぶさない。catch ブロックには最低でも `console.error` でログを出力すること。

### 未使用変数を残さない

使っていない変数・import・関数引数は削除する。シグネチャを変えられない場合（イベントハンドラなど）はプレフィックスに `_` をつける。

```ts
// NG
const unused = 123;

// OK（シグネチャ固定の場合）
function handler(_event: MouseEvent) { ... }
```

### エクスポートする関数には戻り値型を明記する

エクスポートする関数・React コンポーネントには戻り値の型アノテーションをつける。

### `any` 型を使わない

`any` は使わない。型が不明な場合は `unknown` を使い、明示的に型を絞り込む。

---

## React

### イベントハンドラの命名

イベントハンドラは `handle<Event>` 形式で命名する（例: `handleClick`, `handleSubmit`）。

### 条件付きレンダリング

シンプルな場合は `condition && <Component />` を使う。両方描画する場合は三項演算子 `condition ? <A /> : <B />`。ネストした三項演算子は避ける。

### Context のファイル分割

1ファイルにコンポーネント・hook・context を混在させると `react-refresh/only-export-components` 警告が出る。
Context を追加するときは `src/context/<featureName>/` サブフォルダを作り、以下に分割する。

```
src/context/imageGallery/
  ImageGalleryContext.tsx   # createContext + Provider（コンポーネントのみ export）
  useImageGallery.ts        # useContext ラッパー hook（関数のみ export）
  ImageGallery.tsx          # Context を使う UI コンポーネント（あれば）
```

---

## その他

### コメントアウトしたコードを残さない

コメントアウトしたコードをリポジトリに残さない。不要なコードは削除し、必要なら git 履歴から復元する。

### Lint エラーは必ず解消する

ESLint / TypeScript のコンパイルエラーはコミット前に必ず解消する。警告も原則対応し、抑制する場合はインラインコメントで理由を記載する。
