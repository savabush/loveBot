package entities

type FoodCardKey string

type FoodCard struct {
	Name          string
	Key           FoodCardKey
	Description   string
	Price         []uint
	Currency      []uint8
	TimeCooking   uint
	PhotoFilePath string
}
