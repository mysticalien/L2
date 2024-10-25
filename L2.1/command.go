package main

import "fmt"

// Паттерн «Команда» инкапсулирует запрос в виде объекта, что позволяет параметризовать клиентов с различными запросами,
// ставить запросы в очередь или протоколировать их.

// Паттерн Команда
//
// Применимость:
// - Этот паттерн используется, когда нужно отделить запрос от его обработчика, создавая объекты, которые инкапсулируют действия и параметры для их выполнения.
// - Полезен, когда требуется откладывать выполнение команд, регистрировать историю действий или поддерживать отмену операций.
//
// Плюсы:
// - Позволяет создавать очереди команд и регистрировать историю выполнения действий.
// - Упрощает реализацию операций отмены и повторного выполнения.
//
// Минусы:
// - Увеличивает количество классов в системе.
// - Усложняет код, особенно если команды с большим количеством параметров.
//
// Реальные примеры использования:
// - Системы управления интерфейсом, где каждый пользовательский ввод (например, нажатие кнопки) представляется как команда.
// - Очереди задач и операции, которые можно отменить в текстовых редакторах.

// Command Interface
type Command interface {
	Execute()
}

// Commands
type PawCommand struct{}

func (c *PawCommand) Execute() {
	fmt.Println("Pakkun gives paw")
}

type SitCommand struct{}

func (c *SitCommand) Execute() {
	fmt.Println("Pakkun sits")
}

// Invoker (command owner)
type PakkunTrainer struct {
	command Command
}

func (t *PakkunTrainer) SetCommand(command Command) {
	t.command = command
}

func (t *PakkunTrainer) TeachCommand() {
	t.command.Execute()
}

func main() {
	trainer := &PakkunTrainer{}

	paw := &PawCommand{}
	sit := &SitCommand{}

	trainer.SetCommand(paw)
	trainer.TeachCommand()

	trainer.SetCommand(sit)
	trainer.TeachCommand()
}
