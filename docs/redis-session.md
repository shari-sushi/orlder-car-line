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

```json
{
  "user_id": "Uxxxxxxxx",
  "step": "waiting_year",
  "data": {
    "car_name": "アルファード",
    "selected_car_id": "MA_0001"
  }
}
```

Redis キー: `session:{user_id}`、TTL: 36時間

## ステートマシン

```txt
（初期状態 / セッションなし）
  ↓ 車種名を入力
waiting_car_select    クイックリプライで選択肢提示
  ↓ ユーザーが選択
waiting_year          年式入力待ち
  ↓ 年式を入力
waiting_distance      走行距離入力待ち
  ↓ 走行距離を入力
（完了）              Flex Message で URL 返却 → セッション削除
```

タイムアウト後にメッセージが来た場合は「最初からやり直し」メッセージを返す。

## Go での実装イメージ

```go
type Session struct {
    UserID string            `json:"user_id"`
    Step   string            `json:"step"`
    Data   map[string]string `json:"data"`
}

// 保存（TTL: 30 分）
func SaveSession(ctx context.Context, rdb *redis.Client, s Session) error {
    b, _ := json.Marshal(s)
    return rdb.Set(ctx, "session:"+s.UserID, b, 30*time.Minute).Err()
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
```
