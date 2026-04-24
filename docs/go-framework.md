# Go フレームワーク比較

入社先のコードベースに合わせる方針。確認中。

## 候補

| フレームワーク     | 特徴                       | GitHub Stars |
| ------------------ | -------------------------- | ------------ |
| `net/http`（標準） | 依存なし、Go らしい書き方  | -            |
| Echo               | シンプル、ミドルウェア充実 | 30k+         |
| Gin                | 高速、大規模コミュニティ   | 79k+         |
| Chi                | 軽量、標準ライブラリ準拠   | 18k+         |

## 各フレームワークの概要

### net/http（標準ライブラリ）

依存ゼロで動く。小規模・学習用には十分。

```go
http.HandleFunc("/webhook", handleWebhook)
http.ListenAndServe(":8080", nil)
```

### Echo

```go
e := echo.New()
e.POST("/webhook", handleWebhook)
e.Start(":8080")
```

ルーティング・バリデーション・ミドルウェアが充実。
入社先での採用実績次第で選定。

### Gin

Echo と似た API。GitHub スター数が最多でコミュニティが大きい。

```go
r := gin.Default()
r.POST("/webhook", handleWebhook)
r.Run(":8080")
```

### Chi

標準の `net/http` の `Handler` インターフェースを完全準拠。
フレームワーク乗り換えが容易で、テストしやすい。

```go
r := chi.NewRouter()
r.Post("/webhook", handleWebhook)
http.ListenAndServe(":8080", r)
```

## 選定状況

- [ ] 入社先に確認
- [ ] 選定決定 → このファイルを更新
