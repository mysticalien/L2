package main

import "fmt"

// Паттерн «Фасад» предоставляет простой интерфейс к сложной системе классов, библиотек или фреймворков.
// Он скрывает сложность системы, облегчая использование.

// Паттерн Фасад
//
// Применимость:
// - Этот паттерн используется, когда необходимо упростить интерфейс сложной системы, предоставив единый вход для работы с ней.
// - Удобен, когда нужно скрыть сложную логику или детали реализации, оставив доступ только к основным функциям.
//
// Плюсы:
// - Скрывает сложность системы, упрощая работу с ней.
// - Уменьшает количество зависимостей между клиентом и подсистемой.
//
// Минусы:
// - Может приводить к чрезмерному упрощению, скрывая важные возможности подсистемы.
// - Фасад может стать точкой отказа при изменении внутренней логики системы.
//
// Реальные примеры использования:
// - API библиотек, предоставляющие простые методы для работы с комплексной функциональностью (например, библиотеки для работы с сетями).
// - Интерфейсы для работы с внешними службами и модулями, такими как базы данных, где фасад упрощает взаимодействие с их внутренней структурой.

// Subsystems
type Massage struct{}

func (m *Massage) Start() {
	fmt.Println("Pakkun gets a massage...")
}

type Brushing struct{}

func (b *Brushing) Start() {
	fmt.Println("Pakkun is being brushed...")
}

type Treat struct{}

func (t *Treat) Serve() {
	fmt.Println("Pakkun was served a treat...")
}

// Facade
type CatSpaFacade struct {
	massage  *Massage
	brushing *Brushing
	treat    *Treat
}

func NewCatSpaFacade() *CatSpaFacade {
	return &CatSpaFacade{
		massage:  &Massage{},
		brushing: &Brushing{},
		treat:    &Treat{},
	}
}

func (f *CatSpaFacade) StartSpaDay() {
	f.massage.Start()
	f.brushing.Start()
	f.treat.Serve()
}

func main() {
	spa := NewCatSpaFacade()
	spa.StartSpaDay()
}
