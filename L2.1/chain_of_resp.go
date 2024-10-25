package main

import "fmt"

// Цепочка ответственности - это поведенческий шаблон, позволяющий передавать запросы по цепочке обработчиков.
// Получив запрос, каждый обработчик решает, обработать его или передать следующему обработчику в цепочке.

// Паттерн Цепочка вызовов
//
// Применимость:
// - Этот паттерн используется, когда несколько объектов могут обработать запрос, и вы хотите определить получателя запроса на этапе выполнения.
// - Полезен, если требуется выполнение проверки или обработки запроса последовательно через различные уровни проверки.
//
// Плюсы:
// - Ослабляет связь между отправителем и получателем запроса, что упрощает изменение логики обработки.
// - Позволяет гибко настраивать цепочку вызовов в разных конфигурациях.
//
// Минусы:
// - Запрос может остаться необработанным, если ни один из объектов цепочки его не примет.
// - Увеличивает накладные расходы, если цепочка слишком длинная.
//
// Реальные примеры использования:
// - Middleware в веб-фреймворках, например, проверка аутентификации и логирование запросов на серверах.
// - Библиотеки логирования, где несколько логгеров обрабатывают сообщения разного уровня.

// Chain link interface
type AnimalHandler interface {
	SetNext(handler AnimalHandler)
	Handle(request string)
}

// Pakkun
type CatPakkun struct {
	next AnimalHandler
}

func (c *CatPakkun) SetNext(handler AnimalHandler) {
	c.next = handler
}

func (c *CatPakkun) Handle(request string) {
	if c.next != nil {
		c.next.Handle(request)
	} else if request == "help" {
		fmt.Println("Pakkun replies, 'If Mo is busy, I'll help.'")
	}
}

// Mo
type DogMo struct {
	next AnimalHandler
}

func (d *DogMo) SetNext(handler AnimalHandler) {
	d.next = handler
}

func (d *DogMo) Handle(request string) {
	if d.next != nil {
		d.next.Handle(request)
	} else if request == "help" {
		fmt.Println("Moe replies, 'I'll help first!'")
	}
}

func main() {
	mo := &DogMo{}
	pakkun := &CatPakkun{}

	pakkun.SetNext(mo)

	pakkun.Handle("help")
}
