package sticker

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/savabush/breakfastLoveBot/internal/entities"
)

const autoSaveInterval = time.Minute * 1

var (
	Sticker *MemoryRepository
	once    sync.Once
)

type MemoryRepository struct {
	// cycle buffer
	stickers []entities.Sticker
	// use for cyclic access
	curIndex uint
	mu       *sync.RWMutex
	filePath string
}

func NewMemoryRepository() *MemoryRepository {
	once.Do(func() {
		Sticker = &MemoryRepository{
			stickers: make([]entities.Sticker, 0),
			mu:       &sync.RWMutex{},
			filePath: "data/stickers.json",
		}

		Sticker.loadStickers()

		go Sticker.autoSave()
	})

	return Sticker
}

func (r *MemoryRepository) Add(sticker entities.Sticker) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.stickers = append(r.stickers, sticker)
	return nil
}

// GetNext returns the next sticker cyclically
func (r *MemoryRepository) GetNext() entities.Sticker {
	r.mu.RLock()
	if len(r.stickers) == 0 {
		return entities.Sticker{}
	}
	sticker := r.stickers[r.curIndex]
	stickersCount := len(r.stickers)
	r.mu.RUnlock()

	r.mu.Lock()
	defer r.mu.Unlock()
	// move to the next index cyclically
	r.curIndex = (r.curIndex + 1) % uint(stickersCount)
	return sticker
}

func (r *MemoryRepository) HasStickers() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.stickers) > 0
}

func (r *MemoryRepository) loadStickers() {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return
	}

	var items []entities.Sticker
	err = json.Unmarshal(data, &items)
	if err != nil {
		log.Printf("failed to unmarshal stickers: %v", err)
		return
	}
	r.mu.Lock()
	r.stickers = items
	r.mu.Unlock()
}

func (r *MemoryRepository) saveStickers() {
	r.mu.RLock()
	data, err := json.MarshalIndent(r.stickers, "", "  ")
	r.mu.RUnlock()
	if err != nil {
		log.Printf("failed to marshal stickers: %v", err)
		return
	}
	err = os.WriteFile(r.filePath, data, 0644)
	if err != nil {
		log.Printf("failed to write stickers to file: %v", err)

		if errors.Is(err, os.ErrNotExist) {
			dirName := strings.Split(r.filePath, "/")[0]
			err = os.Mkdir(dirName, os.FileMode(0755))
			if errors.Is(err, os.ErrExist) {
				log.Println("dir already created ", dirName)
			}
		}
	}
}

func (r *MemoryRepository) autoSave() {
	ticker := time.NewTicker(autoSaveInterval)
	defer ticker.Stop()

	for range ticker.C {
		r.saveStickers()
	}
}

func (r *MemoryRepository) Close() {
	r.saveStickers()
}
