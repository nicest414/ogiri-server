package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

var (
	ErrNotFound = errors.New("項目が見つかりません")
)

// Theme はお題を表す構造体
type Theme struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   string    `json:"created_by"`
	Active      bool      `json:"active"`
}

// Answer は大喜利の回答を表す構造体
type Answer struct {
	ID        string    `json:"id"`
	ThemeID   string    `json:"theme_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy string    `json:"created_by"`
	Likes     int       `json:"likes"`
}

// DataStore はデータ操作のためのインターフェース
type DataStore interface {
	// お題関連
	GetTheme(id string) (*Theme, error)
	ListThemes() ([]*Theme, error)
	CreateTheme(theme *Theme) error
	UpdateTheme(theme *Theme) error
	DeleteTheme(id string) error

	// 回答関連
	GetAnswer(id string, themeID string) (*Answer, error)
	ListAnswers(themeID string) ([]*Answer, error)
	CreateAnswer(answer *Answer) error
	UpdateAnswer(answer *Answer) error
	DeleteAnswer(id string, themeID string) error
}

// InMemoryStore はメモリ内にデータを保持する実装
type InMemoryStore struct {
	themes       map[string]*Theme
	answers      map[string]map[string]*Answer
	themesMutex  sync.RWMutex
	answersMutex sync.RWMutex
}

// NewInMemoryStore は新しいInMemoryStoreインスタンスを返す
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		themes:  make(map[string]*Theme),
		answers: make(map[string]map[string]*Answer),
	}
}

// GetTheme はIDからテーマを取得
func (s *InMemoryStore) GetTheme(id string) (*Theme, error) {
	s.themesMutex.RLock()
	defer s.themesMutex.RUnlock()

	theme, exists := s.themes[id]
	if !exists {
		return nil, ErrNotFound
	}
	return theme, nil
}

// ListThemes は全てのテーマをリストアップ
func (s *InMemoryStore) ListThemes() ([]*Theme, error) {
	s.themesMutex.RLock()
	defer s.themesMutex.RUnlock()

	themes := make([]*Theme, 0, len(s.themes))
	for _, theme := range s.themes {
		themes = append(themes, theme)
	}
	return themes, nil
}

// CreateTheme は新しいテーマを作成
func (s *InMemoryStore) CreateTheme(theme *Theme) error {
	s.themesMutex.Lock()
	defer s.themesMutex.Unlock()

	s.themes[theme.ID] = theme
	return nil
}

// UpdateTheme はテーマを更新
func (s *InMemoryStore) UpdateTheme(theme *Theme) error {
	s.themesMutex.Lock()
	defer s.themesMutex.Unlock()

	_, exists := s.themes[theme.ID]
	if !exists {
		return ErrNotFound
	}

	s.themes[theme.ID] = theme
	return nil
}

// DeleteTheme はテーマを削除
func (s *InMemoryStore) DeleteTheme(id string) error {
	s.themesMutex.Lock()
	defer s.themesMutex.Unlock()

	_, exists := s.themes[id]
	if !exists {
		return ErrNotFound
	}

	delete(s.themes, id)
	return nil
}

// GetAnswer は回答を取得
func (s *InMemoryStore) GetAnswer(id string, themeID string) (*Answer, error) {
	s.answersMutex.RLock()
	defer s.answersMutex.RUnlock()

	themeAnswers, exists := s.answers[themeID]
	if !exists {
		return nil, ErrNotFound
	}

	answer, exists := themeAnswers[id]
	if !exists {
		return nil, ErrNotFound
	}

	return answer, nil
}

// ListAnswers はテーマに対する全ての回答を取得
func (s *InMemoryStore) ListAnswers(themeID string) ([]*Answer, error) {
	s.answersMutex.RLock()
	defer s.answersMutex.RUnlock()

	themeAnswers, exists := s.answers[themeID]
	if !exists {
		return []*Answer{}, nil
	}

	answers := make([]*Answer, 0, len(themeAnswers))
	for _, answer := range themeAnswers {
		answers = append(answers, answer)
	}
	return answers, nil
}

// CreateAnswer は新しい回答を作成
func (s *InMemoryStore) CreateAnswer(answer *Answer) error {
	s.answersMutex.Lock()
	defer s.answersMutex.Unlock()

	if _, exists := s.answers[answer.ThemeID]; !exists {
		s.answers[answer.ThemeID] = make(map[string]*Answer)
	}

	s.answers[answer.ThemeID][answer.ID] = answer
	return nil
}

// UpdateAnswer は回答を更新
func (s *InMemoryStore) UpdateAnswer(answer *Answer) error {
	s.answersMutex.Lock()
	defer s.answersMutex.Unlock()

	themeAnswers, exists := s.answers[answer.ThemeID]
	if !exists {
		return ErrNotFound
	}

	_, exists = themeAnswers[answer.ID]
	if !exists {
		return ErrNotFound
	}

	themeAnswers[answer.ID] = answer
	return nil
}

// DeleteAnswer は回答を削除
func (s *InMemoryStore) DeleteAnswer(id string, themeID string) error {
	s.answersMutex.Lock()
	defer s.answersMutex.Unlock()

	themeAnswers, exists := s.answers[themeID]
	if !exists {
		return ErrNotFound
	}

	_, exists = themeAnswers[id]
	if !exists {
		return ErrNotFound
	}

	delete(themeAnswers, id)
	return nil
}

// JSONファイル用のデータ構造
type JSONData struct {
	Themes       map[string]*Theme  `json:"themes"`
	Answers      map[string]*Answer `json:"answers"`
	NextThemeID  int               `json:"next_theme_id"`
	NextAnswerID int               `json:"next_answer_id"`
}

// JSONファイルベースのデータストア
type JSONStore struct {
	mu           sync.RWMutex
	themes       map[string]*Theme
	answers      map[string]*Answer
	filePath     string
	nextThemeID  int
	nextAnswerID int
}

// 新しいJSONストアを作成
func NewJSONStore(filePath string) *JSONStore {
	store := &JSONStore{
		themes:       make(map[string]*Theme),
		answers:      make(map[string]*Answer),
		filePath:     filePath,
		nextThemeID:  1,
		nextAnswerID: 1,
	}
	
	// ファイルからデータを読み込み
	store.loadFromFile()
	
	return store
}

// ファイルからデータを読み込み
func (s *JSONStore) loadFromFile() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// ファイルが存在しない場合は何もしない
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		return nil
	}
	
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return fmt.Errorf("ファイル読み込みエラー: %w", err)
	}
	
	var jsonData JSONData
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return fmt.Errorf("JSON解析エラー: %w", err)
	}
	
	s.themes = jsonData.Themes
	s.answers = jsonData.Answers
	s.nextThemeID = jsonData.NextThemeID
	s.nextAnswerID = jsonData.NextAnswerID
	
	// nilマップの初期化
	if s.themes == nil {
		s.themes = make(map[string]*Theme)
	}
	if s.answers == nil {
		s.answers = make(map[string]*Answer)
	}
	
	return nil
}

// ファイルにデータを保存
func (s *JSONStore) saveToFile() error {
	jsonData := JSONData{
		Themes:       s.themes,
		Answers:      s.answers,
		NextThemeID:  s.nextThemeID,
		NextAnswerID: s.nextAnswerID,
	}
	
	data, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON変換エラー: %w", err)
	}
	
	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("ファイル書き込みエラー: %w", err)
	}
	
	return nil
}

// GetTheme implements DataStore
func (s *JSONStore) GetTheme(id string) (*Theme, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	theme, exists := s.themes[id]
	if !exists {
		return nil, ErrNotFound
	}
	
	return theme, nil
}

// ListThemes implements DataStore
func (s *JSONStore) ListThemes() ([]*Theme, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	themes := make([]*Theme, 0, len(s.themes))
	for _, theme := range s.themes {
		themes = append(themes, theme)
	}
	
	return themes, nil
}

// CreateTheme implements DataStore
func (s *JSONStore) CreateTheme(theme *Theme) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// IDを自動生成
	theme.ID = fmt.Sprintf("theme_%d", s.nextThemeID)
	theme.CreatedAt = time.Now()
	theme.UpdatedAt = time.Now()
	theme.Active = true
	
	s.themes[theme.ID] = theme
	s.nextThemeID++
	
	// ファイルに保存
	return s.saveToFile()
}

// UpdateTheme implements DataStore
func (s *JSONStore) UpdateTheme(theme *Theme) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if _, exists := s.themes[theme.ID]; !exists {
		return ErrNotFound
	}
	
	theme.UpdatedAt = time.Now()
	s.themes[theme.ID] = theme
	
	// ファイルに保存
	return s.saveToFile()
}

// DeleteTheme implements DataStore
func (s *JSONStore) DeleteTheme(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if _, exists := s.themes[id]; !exists {
		return ErrNotFound
	}
	
	delete(s.themes, id)
	
	// 関連する回答も削除
	for answerID, answer := range s.answers {
		if answer.ThemeID == id {
			delete(s.answers, answerID)
		}
	}
	
	// ファイルに保存
	return s.saveToFile()
}

// GetAnswer implements DataStore
func (s *JSONStore) GetAnswer(id string, themeID string) (*Answer, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	answer, exists := s.answers[id]
	if !exists || answer.ThemeID != themeID {
		return nil, ErrNotFound
	}
	
	return answer, nil
}

// ListAnswers implements DataStore
func (s *JSONStore) ListAnswers(themeID string) ([]*Answer, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	answers := make([]*Answer, 0)
	for _, answer := range s.answers {
		if answer.ThemeID == themeID {
			answers = append(answers, answer)
		}
	}
	
	return answers, nil
}

// CreateAnswer implements DataStore
func (s *JSONStore) CreateAnswer(answer *Answer) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// お題の存在確認
	if _, exists := s.themes[answer.ThemeID]; !exists {
		return ErrNotFound
	}
	
	// IDを自動生成
	answer.ID = fmt.Sprintf("answer_%d", s.nextAnswerID)
	answer.CreatedAt = time.Now()
	answer.UpdatedAt = time.Now()
	
	s.answers[answer.ID] = answer
	s.nextAnswerID++
	
	// ファイルに保存
	return s.saveToFile()
}

// UpdateAnswer implements DataStore
func (s *JSONStore) UpdateAnswer(answer *Answer) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	existing, exists := s.answers[answer.ID]
	if !exists || existing.ThemeID != answer.ThemeID {
		return ErrNotFound
	}
	
	answer.UpdatedAt = time.Now()
	s.answers[answer.ID] = answer
	
	// ファイルに保存
	return s.saveToFile()
}

// DeleteAnswer implements DataStore
func (s *JSONStore) DeleteAnswer(id string, themeID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	answer, exists := s.answers[id]
	if !exists || answer.ThemeID != themeID {
		return ErrNotFound
	}
	
	delete(s.answers, id)
	
	// ファイルに保存
	return s.saveToFile()
}
