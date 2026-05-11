package linebot

import (
	"bytes"
	_ "embed"
	"log"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"

	"orlder-car-line/server/internal/infra/config"
)

//go:embed richmenu/main/v1_3types.jpg
var richMenuImage []byte

const (
	richMenuWidth  = 1536
	richMenuHeight = 612
)

// SetupRichMenu はリッチメニューを作成・デフォルト設定する。
// サーバー起動時に goroutine で呼び出す。
//
// 画像は LINE API の上限（1MB）以下である必要がある。
// PNG は圧縮後も超えやすいため JPEG（~130KB）を使用する。
func SetupRichMenu() {
	apiClient, err := messaging_api.NewMessagingApiAPI(config.LineChannelAccessToken)
	if err != nil {
		log.Printf("[richmenu] APIクライアント初期化エラー: %v", err)
		return
	}
	blobClient, err := messaging_api.NewMessagingApiBlobAPI(config.LineChannelAccessToken)
	if err != nil {
		log.Printf("[richmenu] BlobAPIクライアント初期化エラー: %v", err)
		return
	}

	// 既存のリッチメニューをすべて削除してからリセット
	if list, err := apiClient.GetRichMenuList(); err != nil {
		log.Printf("[richmenu] GetRichMenuList エラー: %v", err)
	} else {
		for _, m := range list.Richmenus {
			if _, err := apiClient.DeleteRichMenu(m.RichMenuId); err != nil {
				log.Printf("[richmenu] DeleteRichMenu エラー (%s): %v", m.RichMenuId, err)
			}
		}
		log.Printf("[richmenu] 既存メニュー %d 件を削除", len(list.Richmenus))
	}

	req := &messaging_api.RichMenuRequest{
		Size:        &messaging_api.RichMenuSize{Width: richMenuWidth, Height: richMenuHeight},
		Selected:    true,
		Name:        "メインメニュー",
		ChatBarText: "メニュー",
		Areas: []messaging_api.RichMenuArea{
			{
				Bounds: &messaging_api.RichMenuBounds{X: 0, Y: 0, Width: 512, Height: richMenuHeight},
				Action: &messaging_api.MessageAction{
					Label: "SUVを探す",
					Text:  "SUVを探す",
				},
			},
			{
				Bounds: &messaging_api.RichMenuBounds{X: 512, Y: 0, Width: 512, Height: richMenuHeight},
				Action: &messaging_api.MessageAction{
					Label: "スタッフに相談する",
					Text:  "スタッフに相談する",
				},
			},
			{
				Bounds: &messaging_api.RichMenuBounds{X: 1024, Y: 0, Width: 512, Height: richMenuHeight},
				Action: &messaging_api.MessageAction{
					Label: "自分で探す",
					Text:  "自分で探す",
				},
			},
		},
	}

	res, err := apiClient.CreateRichMenu(req)
	if err != nil {
		log.Printf("[richmenu] CreateRichMenu エラー: %v", err)
		return
	}
	richMenuID := res.RichMenuId
	log.Printf("[richmenu] 作成完了: %s", richMenuID)

	if _, err = blobClient.SetRichMenuImage(richMenuID, "image/jpeg", bytes.NewReader(richMenuImage)); err != nil {
		log.Printf("[richmenu] SetRichMenuImage エラー: %v", err)
		return
	}
	log.Printf("[richmenu] 画像アップロード完了")

	if _, err = apiClient.SetDefaultRichMenu(richMenuID); err != nil {
		log.Printf("[richmenu] SetDefaultRichMenu エラー: %v", err)
		return
	}
	log.Printf("[richmenu] デフォルト設定完了: %s", richMenuID)
}
