package main

import "fmt"

// Состояние - это поведенческий паттерн, позволяющий объекту изменять свое поведение при изменении внутреннего состояния.
// Это выглядит так, как будто объект изменил свой класс.

// Паттерн Состояние
//
// Применимость:
// - Этот паттерн используется, если объект должен изменять свое поведение в зависимости от внутреннего состояния.
// - Полезен, если требуется обработка одного и того же события по-разному в зависимости от состояния объекта.
//
// Плюсы:
// - Упрощает добавление новых состояний и их поведения.
// - Изолирует поведение, зависящее от состояния, от других частей программы.
//
// Минусы:
// - Увеличивает сложность кода, добавляя отдельные классы для каждого состояния.
// - Усложняет поддержку, если требуется большое количество состояний.
//
// Реальные примеры использования:
// - Банкоматы, где действия зависят от состояния (например, "Без карты", "Авторизация", "Выдача наличных").
// - Процесс заказов, где заказ может находиться в разных состояниях, таких как "Создан", "Оплачен", "Отправлен", "Доставлен".

// Status interface
type Mood interface {
	Act()
}

// A state of hunger
type HungryState struct{}

func (s *HungryState) Act() {
	fmt.Println("Pakkun meows and begs for food")
}

// Satiety state
type FullState struct{}

func (s *FullState) Act() {
	fmt.Println("Pakkun purrs and rests.")
}

// Context
type Pakkun struct {
	state Mood
}

func (c *Pakkun) SetState(state Mood) {
	c.state = state
}

func (c *Pakkun) Behave() {
	c.state.Act()
}

func main() {
	pakkun := &Pakkun{}

	fmt.Println("Enter state ('hungry' or 'full'):")
	var state string
	_, err := fmt.Scan(&state)
	if err != nil {
		fmt.Println("Wrong state, pls choose between 'hungry' or 'full'")
		return
	}

	switch state {
	case "hungry":
		hungry := &HungryState{}
		pakkun.SetState(hungry)
		pakkun.Behave()
	case "full":
		full := &FullState{}
		pakkun.SetState(full)
		pakkun.Behave()
	default:
		fmt.Println("Wrong state, pls choose between 'hungry' or 'full'")
	}
}
