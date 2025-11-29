package cart

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
	foodCart map[entities.UserTelegramID][]entities.FoodCard
	mu       *sync.RWMutex
	filePath string
}

func NewMemoryRepository(userUDs []entities.UserTelegramID) *MemoryRepository {
	memRepo := &MemoryRepository{
		foodCart: make(map[entities.UserTelegramID][]entities.FoodCard, 0),
		mu:       &sync.RWMutex{},
		filePath: "data/cart.json",
	}
	memRepo.loadFoodCart()
	go memRepo.autoSave()
	return memRepo
}

func (r *MemoryRepository) AddFood(userID entities.UserTelegramID, foodCard entities.FoodCard) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.foodCart[userID]; !ok {
		r.foodCart[userID] = []entities.FoodCard{}
	}
	r.foodCart[userID] = append(r.foodCart[userID], foodCard)
}

func (r *MemoryRepository) GetFoodCart(userID entities.UserTelegramID) []entities.FoodCard {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if _, ok := r.foodCart[userID]; !ok {
		return []entities.FoodCard{}
	}

	if len(r.foodCart[userID]) == 0 {
		return []entities.FoodCard{}
	}

	return r.foodCart[userID]
}

func (r *MemoryRepository) AcceptFoodCart(userID entities.UserTelegramID) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.foodCart, userID)
}

func (r *MemoryRepository) CleanFoodCart(userID entities.UserTelegramID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.foodCart[userID] = []entities.FoodCard{}
	return nil
}

func (r *MemoryRepository) loadFoodCart() {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return
	}

	var items map[entities.UserTelegramID][]entities.FoodCard = make(map[entities.UserTelegramID][]entities.FoodCard)
	err = json.Unmarshal(data, &items)
	if err != nil {
		log.Printf("failed to unmarshal stickers: %v", err)
		return
	}
	r.mu.Lock()
	r.foodCart = items
	r.mu.Unlock()
}

func (r *MemoryRepository) saveCart() {
	r.mu.RLock()
	data, err := json.MarshalIndent(r.foodCart, "", "  ")
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
		r.saveCart()
	}
}

func (r *MemoryRepository) Close() {
	r.saveCart()
}
