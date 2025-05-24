package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/ketya/ogiri-server/internal/data"

	"crypto/rand"
	"encoding/hex"
)

// Handler はAPIハンドラーを管理する構造体
type Handler struct {
	store data.DataStore
}

// NewHandler は新しいHandlerインスタンスを返す
func NewHandler(store data.DataStore) *Handler {
	return &Handler{store: store}
}

// 一意のIDを生成する
func generateID() string {
	bytes := make([]byte, 8) // 16文字のIDになる
	if _, err := rand.Read(bytes); err != nil {
		// エラーが発生した場合はタイムスタンプを使用
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}

// エラーレスポンスを送信するヘルパー関数
func sendErrorResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// JSONレスポンスを送信するヘルパー関数
func sendJSONResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// ---------- お題関連のハンドラー ----------

// ListThemes は全てのお題をリストアップ
func (h *Handler) ListThemes(w http.ResponseWriter, r *http.Request) {
	themes, err := h.store.ListThemes()
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "お題の取得に失敗しました")
		return
	}
	sendJSONResponse(w, http.StatusOK, themes)
}

// GetTheme は特定のお題を取得
func (h *Handler) GetTheme(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	theme, err := h.store.GetTheme(id)
	if err == data.ErrNotFound {
		sendErrorResponse(w, http.StatusNotFound, "お題が見つかりません")
		return
	}
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "お題の取得に失敗しました")
		return
	}
	sendJSONResponse(w, http.StatusOK, theme)
}

// CreateTheme は新しいお題を作成
func (h *Handler) CreateTheme(w http.ResponseWriter, r *http.Request) {
	var theme data.Theme
	if err := json.NewDecoder(r.Body).Decode(&theme); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "無効なリクエスト形式です")
		return
	}

	// バリデーション
	if theme.Title == "" {
		sendErrorResponse(w, http.StatusBadRequest, "タイトルは必須です")
		return
	}

	// IDと時間の設定
	now := time.Now()
	theme.ID = generateID()
	theme.CreatedAt = now
	theme.UpdatedAt = now
	theme.Active = true

	if err := h.store.CreateTheme(&theme); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "お題の作成に失敗しました")
		return
	}

	sendJSONResponse(w, http.StatusCreated, theme)
}

// UpdateTheme はお題を更新
func (h *Handler) UpdateTheme(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var updatedTheme data.Theme
	if err := json.NewDecoder(r.Body).Decode(&updatedTheme); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "無効なリクエスト形式です")
		return
	}

	// 現在のお題を取得
	currentTheme, err := h.store.GetTheme(id)
	if err == data.ErrNotFound {
		sendErrorResponse(w, http.StatusNotFound, "お題が見つかりません")
		return
	}
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "お題の取得に失敗しました")
		return
	}

	// 更新されたフィールドを適用
	if updatedTheme.Title != "" {
		currentTheme.Title = updatedTheme.Title
	}
	if updatedTheme.Description != "" {
		currentTheme.Description = updatedTheme.Description
	}
	currentTheme.Active = updatedTheme.Active
	currentTheme.UpdatedAt = time.Now()

	if err := h.store.UpdateTheme(currentTheme); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "お題の更新に失敗しました")
		return
	}

	sendJSONResponse(w, http.StatusOK, currentTheme)
}

// DeleteTheme はお題を削除
func (h *Handler) DeleteTheme(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.store.DeleteTheme(id); err == data.ErrNotFound {
		sendErrorResponse(w, http.StatusNotFound, "お題が見つかりません")
		return
	} else if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "お題の削除に失敗しました")
		return
	}

	sendJSONResponse(w, http.StatusNoContent, nil)
}

// ---------- 回答関連のハンドラー ----------

// ListAnswers はテーマに対する回答をリストアップ
func (h *Handler) ListAnswers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	themeID := vars["themeID"]

	// テーマの存在確認
	_, err := h.store.GetTheme(themeID)
	if err == data.ErrNotFound {
		sendErrorResponse(w, http.StatusNotFound, "お題が見つかりません")
		return
	}
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "お題の取得に失敗しました")
		return
	}

	answers, err := h.store.ListAnswers(themeID)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "回答の取得に失敗しました")
		return
	}
	sendJSONResponse(w, http.StatusOK, answers)
}

// GetAnswer は特定の回答を取得
func (h *Handler) GetAnswer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	themeID := vars["themeID"]
	id := vars["id"]

	answer, err := h.store.GetAnswer(id, themeID)
	if err == data.ErrNotFound {
		sendErrorResponse(w, http.StatusNotFound, "回答が見つかりません")
		return
	}
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "回答の取得に失敗しました")
		return
	}
	sendJSONResponse(w, http.StatusOK, answer)
}

// SubmitAnswer は新しい回答を投稿
func (h *Handler) SubmitAnswer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	themeID := vars["themeID"]

	// テーマの存在確認
	theme, err := h.store.GetTheme(themeID)
	if err == data.ErrNotFound {
		sendErrorResponse(w, http.StatusNotFound, "お題が見つかりません")
		return
	}
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "お題の取得に失敗しました")
		return
	}

	// 非アクティブなテーマには回答できない
	if !theme.Active {
		sendErrorResponse(w, http.StatusBadRequest, "このお題は現在受付を停止しています")
		return
	}

	var answer data.Answer
	if err := json.NewDecoder(r.Body).Decode(&answer); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "無効なリクエスト形式です")
		return
	}

	// バリデーション
	if answer.Content == "" {
		sendErrorResponse(w, http.StatusBadRequest, "回答内容は必須です")
		return
	}

	// IDと時間の設定
	now := time.Now()
	answer.ID = generateID()
	answer.ThemeID = themeID
	answer.CreatedAt = now
	answer.UpdatedAt = now
	answer.Likes = 0

	if err := h.store.CreateAnswer(&answer); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "回答の投稿に失敗しました")
		return
	}

	sendJSONResponse(w, http.StatusCreated, answer)
}

// UpdateAnswer は回答を更新
func (h *Handler) UpdateAnswer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	themeID := vars["themeID"]
	id := vars["id"]

	var updatedAnswer data.Answer
	if err := json.NewDecoder(r.Body).Decode(&updatedAnswer); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "無効なリクエスト形式です")
		return
	}

	// 現在の回答を取得
	currentAnswer, err := h.store.GetAnswer(id, themeID)
	if err == data.ErrNotFound {
		sendErrorResponse(w, http.StatusNotFound, "回答が見つかりません")
		return
	}
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "回答の取得に失敗しました")
		return
	}

	// 更新されたフィールドを適用
	if updatedAnswer.Content != "" {
		currentAnswer.Content = updatedAnswer.Content
	}
	if updatedAnswer.Likes > 0 {
		currentAnswer.Likes = updatedAnswer.Likes
	}
	currentAnswer.UpdatedAt = time.Now()

	if err := h.store.UpdateAnswer(currentAnswer); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "回答の更新に失敗しました")
		return
	}

	sendJSONResponse(w, http.StatusOK, currentAnswer)
}

// DeleteAnswer は回答を削除
func (h *Handler) DeleteAnswer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	themeID := vars["themeID"]
	id := vars["id"]

	if err := h.store.DeleteAnswer(id, themeID); err == data.ErrNotFound {
		sendErrorResponse(w, http.StatusNotFound, "回答が見つかりません")
		return
	} else if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "回答の削除に失敗しました")
		return
	}

	sendJSONResponse(w, http.StatusNoContent, nil)
}
