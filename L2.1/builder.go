package main

import "fmt"

// Паттерн «Строитель» отделяет конструирование сложных объектов от их представления,
// позволяя создавать разные представления объекта с помощью одного и того же процесса конструирования.

// Паттерн Строитель
//
// Применимость:
// - Этот паттерн используется, когда нужно создать сложный объект с множеством параметров, но при этом контролировать процесс его создания пошагово.
// - Полезен, если требуется создать различные представления объекта, сохраняя код для создания каждого отдельного представления.
//
// Плюсы:
// - Позволяет изменять представление конечного продукта, добавляя или изменяя шаги строительства.
// - Отделяет создание объекта от его представления, упрощая читаемость и поддержку кода.
//
// Минусы:
// - Увеличивает сложность кода, добавляя дополнительные классы.
// - Строитель становится избыточным, если объект прост и не требует много шагов для создания.
//
// Реальные примеры использования:
// - Генераторы для сложных объектов, таких как HTML или XML документы, где требуются последовательные шаги для их построения.
// - Создание объектов с множеством опциональных параметров, таких как настройки подключения к базе данных.

// Product
type CatHouse struct {
	catnipStorage string
	softBed       string
	window        string
}

// Builder
type CatHouseBuilder struct {
	house CatHouse
}

func NewCatHouseBuilder() *CatHouseBuilder {
	return &CatHouseBuilder{}
}

func (b *CatHouseBuilder) AddCatMintStorage() {
	b.house.catnipStorage = "Catnip storage🍬"
}

func (b *CatHouseBuilder) AddSoftBed() {
	b.house.softBed = "Super soft super bed😴"
}

func (b *CatHouseBuilder) AddWindow() {
	b.house.window = "Observation window🪟"
}

func (b *CatHouseBuilder) GetHouse() CatHouse {
	return b.house
}

func main() {
	builder := NewCatHouseBuilder()
	builder.AddCatMintStorage()
	builder.AddSoftBed()
	builder.AddWindow()

	pakkunHouse := builder.GetHouse()
	fmt.Printf("Pakkun's house is ready! It has: %s, %s and %s\n", pakkunHouse.catnipStorage, pakkunHouse.softBed, pakkunHouse.window)
}
