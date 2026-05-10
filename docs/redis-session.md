# Redis セッション設計

## なぜ Redis か

LINE の Webhook はリクエストが来るたびに単発で叩かれる。
HTTP はステートレスなので「このユーザーはどこまで答えたか」を外部に保存する必要がある。

|                        | Redis | MySQL |
| ---------------------- | ----- | ----- |
| 消えても困らないデータ | ✅    |       |
| TTL をつけたい         | ✅    |       |
| 読み書きが超高頻度     | ✅    |       |
| 消えたら困るデータ     |       | ✅    |
| JOIN や複雑な集計      |       | ✅    |
| 監査ログ的な記録       |       | ✅    |

会話状態は「消えても困らない + TTL をつけたい」→ **Redis 確定**

## データ構造

Redis キー: `session:{user_id}`、TTL: 30分

```json
{
  "step": "waiting_price",
  "data": {
    "car_name": "アルファード",
    "year": "2020",
    "distance_min": "0",
    "distance_max": "50000"
  }
}
```

## ステートマシン

```txt
[セッションなし]
  ↓ 車種・年式・距離を3行で入力
  → カーセンサーをスクレイピングして価格帯取得
  → URL + 価格帯をユーザーに返信
  → Redis にセッション保存

waiting_price        価格帯の入力待ち
  ↓ 価格帯を入力（例: 100万円〜200万円）
  → 絞り込みURL を返信
  → セッション削除

[セッションなし に戻る]
```

タイムアウト（TTL切れ）後にメッセージが来た場合は「最初からやり直し」メッセージを返す。

## Go での実装イメージ

```go
type Session struct {
    Step string            `json:"step"`
    Data map[string]string `json:"data"`
}

// 保存（TTL: 30 分）
func SaveSession(ctx context.Context, rdb *redis.Client, userID string, s Session) error {
    b, _ := json.Marshal(s)
    return rdb.Set(ctx, "session:"+userID, b, 30*time.Minute).Err()
}

// 取得
func GetSession(ctx context.Context, rdb *redis.Client, userID string) (*Session, error) {
    val, err := rdb.Get(ctx, "session:"+userID).Result()
    if errors.Is(err, redis.Nil) {
        return nil, nil // セッションなし
    }
    if err != nil {
        return nil, err
    }
    var s Session
    if err := json.Unmarshal([]byte(val), &s); err != nil {
        return nil, err
    }
    return &s, nil
}

// 削除
func DeleteSession(ctx context.Context, rdb *redis.Client, userID string) error {
    return rdb.Del(ctx, "session:"+userID).Err()
}
```
