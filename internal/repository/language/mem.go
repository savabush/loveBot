package language

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"github.com/savabush/breakfastLoveBot/internal/entities"
)

const autoSaveInterval = time.Minute * 1

type MemoryRepository struct {
	langs    map[entities.UserTelegramID]entities.LanguageCode
	mu       *sync.RWMutex
	filePath string
}

func NewMemoryRepository() *MemoryRepository {
	repo := &MemoryRepository{
		langs:    make(map[entities.UserTelegramID]entities.LanguageCode),
		mu:       &sync.RWMutex{},
		filePath: "data/languages.json",
	}
	repo.load()
	go repo.autoSave()
	return repo
}

func (r *MemoryRepository) Get(userID entities.UserTelegramID) (entities.LanguageCode, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	lang, ok := r.langs[userID]
	return lang, ok
}

func (r *MemoryRepository) Set(userID entities.UserTelegramID, lang entities.LanguageCode) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.langs[userID] = lang
	return nil
}

func (r *MemoryRepository) load() {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return
	}
	var items map[entities.UserTelegramID]entities.LanguageCode
	if err := json.Unmarshal(data, &items); err != nil {
		log.Printf("failed to unmarshal languages: %v", err)
		return
	}
	r.mu.Lock()
	r.langs = items
	r.mu.Unlock()
}

func (r *MemoryRepository) save() {
	r.mu.RLock()
	data, err := json.MarshalIndent(r.langs, "", "  ")
	r.mu.RUnlock()
	if err != nil {
		log.Printf("failed to marshal languages: %v", err)
		return
	}
	if err := os.WriteFile(r.filePath, data, 0644); err != nil {
		log.Printf("failed to write languages file: %v", err)
	}
}

func (r *MemoryRepository) autoSave() {
	ticker := time.NewTicker(autoSaveInterval)
	defer ticker.Stop()
	for range ticker.C {
		r.save()
	}
}

func (r *MemoryRepository) Close() {
	r.save()
}
