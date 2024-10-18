package main

import (
    "net/http"
)

// CORSMiddlewareはCORSヘッダーを全許可するミドルウェア
// TODO: アクセス制御を行う
func CORSMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // すべてのオリジンを許可
        w.Header().Set("Access-Control-Allow-Origin", "*")
        // キャッシュの最適化
        w.Header().Set("Vary", "Origin")

        // 許可するHTTPメソッドを設定
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        // 許可するヘッダーを設定
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        // プレフライトリクエストに対する処理
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusNoContent)
            return
        }

        // 次のハンドラーにリクエストを渡す
        next.ServeHTTP(w, r)
    })
}
