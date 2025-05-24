package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/nicest414/ogiri-server/internal/data"
	"github.com/nicest414/ogiri-server/internal/handlers"
)

const (
	defaultPort = "8080"
	dataFile    = "ogiri_data.json"  // JSONファイル名
)

// CORSミドルウェアを実装
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// すべてのオリジンを許可
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// OPTIONSリクエストは処理せずに返す
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 次のハンドラーを実行
		next.ServeHTTP(w, r)
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	// JSONファイルベースのデータストアを初期化
	store := data.NewJSONStore(dataFile)
	log.Printf("📁 データファイル: %s", dataFile)

	// ハンドラー初期化
	h := handlers.NewHandler(store)
	// ルーターの設定
	r := mux.NewRouter()
	
	// お題関連のエンドポイント
	r.HandleFunc("/api/themes", h.ListThemes).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/themes", h.CreateTheme).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/themes/{id}", h.GetTheme).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/themes/{id}", h.UpdateTheme).Methods("PUT", "OPTIONS")
	r.HandleFunc("/api/themes/{id}", h.DeleteTheme).Methods("DELETE", "OPTIONS")

	// 回答関連のエンドポイント
	r.HandleFunc("/api/themes/{themeID}/answers", h.ListAnswers).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/themes/{themeID}/answers", h.SubmitAnswer).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/themes/{themeID}/answers/{id}", h.GetAnswer).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/themes/{themeID}/answers/{id}", h.UpdateAnswer).Methods("PUT", "OPTIONS")
	r.HandleFunc("/api/themes/{themeID}/answers/{id}", h.DeleteAnswer).Methods("DELETE", "OPTIONS")	// CORSミドルウェアを適用
	corsRouter := enableCORS(r)

	// 静的ファイルハンドラー（HTMLテスター用）
	// カレントディレクトリからの静的ファイル提供
	fileServer := http.FileServer(http.Dir("."))
	corsFileServer := enableCORS(fileServer)
	
	http.Handle("/", corsFileServer)
	http.Handle("/api/", corsRouter)	// サーバー起動
	log.Printf("--------------------------------------------------------")
	log.Printf("🎉 大喜利サーバーを起動中...ポート: %s", port)
	log.Printf("💾 データ保存方式: JSONファイル (%s)", dataFile)
	log.Printf("🌐 以下のURLでアクセスできます:")
	log.Printf("   - トップページ: http://localhost:%s/", port)
	log.Printf("   - APIテスター: http://localhost:%s/api_tester.html", port)
	log.Printf("   - お題募集: http://localhost:%s/theme_submission.html", port)
	log.Printf("--------------------------------------------------------")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
