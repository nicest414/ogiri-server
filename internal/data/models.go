package data

import (
	"errors"
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
