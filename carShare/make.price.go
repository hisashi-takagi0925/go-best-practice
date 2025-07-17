package main

import (
	"fmt"
)

const basePrice = 200

// 価格調整ルール用インターフェース
type PriceRule interface {
	Apply(int) int
}

// タイプによる価格調整
type TypeRule struct {
	CarType string
}

func (t TypeRule) Apply(price int) int {
	priceMap := map[string]int{
		"SUV":   100,
		"セダン": 50,
		"軽":    20,
	}
	return price + priceMap[t.CarType]
}

// モデルによる価格調整
type ModelRule struct {
	Model string
}

func (m ModelRule) Apply(price int) int {
	modelMap := map[string]int{
		"高級": 80,
		"中級": 50,
		"低級": 20,
	}
	return price + modelMap[m.Model]
}

// 走行距離による価格調整
type MileageRule struct {
	Mileage int
}

func (mi MileageRule) Apply(price int) int {
	return price - mi.Mileage/10000*10
}

// 年数による価格調整
type YearRule struct {
	Years int
}

func (y YearRule) Apply(price int) int {
	return price - y.Years*15
}

// UsedCar構造体
type UsedCar struct {
	priceRules []PriceRule
}

// コンストラクタ
func NewUsedCar(rules ...PriceRule) *UsedCar {
	return &UsedCar{priceRules: rules}
}

// 価格計算
func (u *UsedCar) Price() int {
	price := basePrice
	for _, rule := range u.priceRules {
		price = rule.Apply(price)
	}
	if price < 0 {
		price = 0
	}
	return price
}

func main() {
	car := NewUsedCar(
		TypeRule{CarType: "SUV"},
		ModelRule{Model: "高級"},
		MileageRule{Mileage: 30000},
		YearRule{Years: 5},
	)
	fmt.Println(car.Price())
}