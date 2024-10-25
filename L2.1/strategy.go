package main

// Стратегия - это поведенческий паттерн проектирования, который позволяет определить семейство алгоритмов,
// поместить каждый из них в отдельный класс и сделать их объекты взаимозаменяемыми.

// Паттерн Стратегия
//
// Применимость:
// - Этот паттерн используется, когда необходимо выбрать алгоритм на этапе выполнения, сохраняя при этом единый интерфейс для всех алгоритмов.
// - Удобен, если часто требуется менять алгоритмы или политики во время работы приложения.
//
// Плюсы:
// - Упрощает замену и добавление новых алгоритмов без изменения клиента.
// - Делает алгоритмы взаимозаменяемыми и облегчает их тестирование.
//
// Минусы:
// - Увеличивает количество классов в программе.
// - Стратегии зависят от интерфейсов, которые должны быть согласованы с клиентом.
//
// Реальные примеры использования:
// - Алгоритмы сортировки, поиска или фильтрации, которые можно менять в зависимости от условий.
// - Оптимизация работы приложения, например, различные методы сжатия данных.

import "fmt"

// Hunting strategy interface
type HuntingStrategy interface {
	Hunt()
}

// Pakkun's strategy is silent hunting.
type SilentHunt struct{}

func (s *SilentHunt) Hunt() {
	fmt.Println("Pakkun sneaks in quietly")
}

// Moe's strategy is to storm
type AttackHunt struct{}

func (a *AttackHunt) Hunt() {
	fmt.Println("Moe thundering forward!")
}

type Hunter struct {
	strategy HuntingStrategy
}

func (h *Hunter) SetStrategy(strategy HuntingStrategy) {
	h.strategy = strategy
}

func (h *Hunter) StartHunting() {
	h.strategy.Hunt()
}

func main() {
	pakkunHunter := &Hunter{}
	pakkunHunter.SetStrategy(&SilentHunt{})
	pakkunHunter.StartHunting()

	moHunter := &Hunter{}
	moHunter.SetStrategy(&AttackHunt{})
	moHunter.StartHunting()
}
