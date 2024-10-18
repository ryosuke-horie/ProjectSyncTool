package main

import (
    "database/sql"
    "log"
    "net/http"

    "ryosuke-horie/project-management-homemade-api/app"
    "github.com/syumai/workers"
    _ "github.com/syumai/workers/cloudflare/d1" // ドライバーを登録
)

func main() {
    // データベース接続の初期化
    db, err := sql.Open("d1", "DB")

	// データベース接続が失敗した場合はエラーログを出力して終了
    if err != nil {
        log.Fatalf("error opening DB: %s", err.Error())
    }
	// DB接続を閉じる
    defer db.Close()
    
    // ハンドラーの設定
    handler := app.NewItemHandler(db)
    
    // CORSミドルウェアを適用
	// 簡易化のため、全てのドメインからのアクセスを許可している TODO: アクセス制御を行う
    corsHandler := CORSMiddleware(handler)
    
    // エンドポイントにハンドラーを登録
    http.Handle("/items", corsHandler)
    http.Handle("/items/for-sync", corsHandler)
    http.Handle("/items/mark-linked", corsHandler)
    http.Handle("/items/issue-closed", corsHandler)

    // Workersのサーバーを起動
    workers.Serve(nil)
}
