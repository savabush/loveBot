package foodCard

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
	"time"

	"github.com/savabush/breakfastLoveBot/internal/entities"
)

const autoSaveInterval = time.Minute * 1

type MemoryRepository struct {
	foodCards []entities.FoodCard
	mu        *sync.RWMutex
	filePath  string
}

func NewMemoryRepository() *MemoryRepository {
	memRepo := &MemoryRepository{
		foodCards: make([]entities.FoodCard, 0),
		mu:        &sync.RWMutex{},
		filePath:  "data/foodCards.json",
	}
	memRepo.loadFoodCards()
	go memRepo.autoSave()
	return memRepo
}

func (r *MemoryRepository) AddNewFood(foodCard entities.FoodCard) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.foodCards = append(r.foodCards, foodCard)
}

func (r *MemoryRepository) GetFoodCardByID(id entities.FoodCardKey) (entities.FoodCard, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, foodCard := range r.foodCards {
		if foodCard.Key == id {
			return foodCard, nil
		}
	}
	return entities.FoodCard{}, errors.New("food card not found")
}

func (r *MemoryRepository) GetAllFoodCards() []entities.FoodCard {
	r.mu.RLock()
	defer r.mu.RUnlock()

	copiedCards := make([]entities.FoodCard, len(r.foodCards))
	copy(copiedCards, r.foodCards)
	return copiedCards
}

func (r *MemoryRepository) UpdateFoodCard(foodCard entities.FoodCard) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.foodCards {
		if r.foodCards[i].Key == foodCard.Key {
			r.foodCards[i] = foodCard
			return nil
		}
	}

	return errors.New("food card not found")
}

func (r *MemoryRepository) DeleteFoodCard(id entities.FoodCardKey) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.foodCards {
		if r.foodCards[i].Key == id {
			r.foodCards = append(r.foodCards[:i], r.foodCards[i+1:]...)
			return nil
		}
	}

	return errors.New("food card not found")
}

func (r *MemoryRepository) loadFoodCards() {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return
	}

	var items []entities.FoodCard
	err = json.Unmarshal(data, &items)
	if err != nil {
		log.Printf("failed to unmarshal stickers: %v", err)
		return
	}
	r.mu.Lock()
	r.foodCards = items
	r.mu.Unlock()
}

func (r *MemoryRepository) saveFoodCards() {
	r.mu.RLock()
	data, err := json.MarshalIndent(r.foodCards, "", "  ")
	r.mu.RUnlock()
	if err != nil {
		log.Printf("failed to marshal stickers: %v", err)
		return
	}
	err = os.WriteFile(r.filePath, data, 0644)
	if err != nil {
		log.Printf("failed to write stickers to file: %v", err)
	}
}

func (r *MemoryRepository) autoSave() {
	ticker := time.NewTicker(autoSaveInterval)
	defer ticker.Stop()

	for range ticker.C {
		r.saveFoodCards()
	}
}

func (r *MemoryRepository) Close() {
	r.saveFoodCards()
}
