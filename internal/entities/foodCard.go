package entities

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"unicode/utf8"
)

type FoodCardKey string

type CurrencyList []string

func (c *CurrencyList) UnmarshalJSON(data []byte) error {
	var list []string
	if err := json.Unmarshal(data, &list); err == nil {
		*c = CurrencyList(list)
		return nil
	}

	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	raw = strings.TrimSpace(raw)
	if raw == "" {
		*c = CurrencyList{}
		return nil
	}
	if decoded, err := base64.StdEncoding.DecodeString(raw); err == nil && utf8.Valid(decoded) {
		*c = CurrencyList{string(decoded)}
		return nil
	}
	*c = CurrencyList{raw}
	return nil
}

type FoodCard struct {
	Name          string
	Key           FoodCardKey
	Description   string
	Price         []uint
	Currency      CurrencyList
	TimeCooking   uint
	PhotoFilePath string
}
