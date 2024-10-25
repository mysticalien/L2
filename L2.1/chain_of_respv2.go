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

type Handler interface {
	SetNext(handler Handler)
	Handle(user *User)
}

type User struct {
	isAuthenticated bool
	isAdmin         bool
}

type AuthHandler struct {
	next Handler
}

func (a *AuthHandler) SetNext(handler Handler) {
	a.next = handler
}

func (a *AuthHandler) Handle(user *User) {
	if !user.isAuthenticated {
		fmt.Println("Access denied: User is not authenticated.")
		return
	}
	fmt.Println("User authenticated.")
	if a.next != nil {
		a.next.Handle(user)
	}
}

type AdminHandler struct {
	next Handler
}

func (a *AdminHandler) SetNext(handler Handler) {
	a.next = handler
}

func (a *AdminHandler) Handle(user *User) {
	if user.isAdmin {
		fmt.Println("Full access granted to admin.")
	} else {
		fmt.Println("Limited access granted to user.")
	}
	if a.next != nil {
		a.next.Handle(user)
	}
}

func main() {
	authHandler := &AuthHandler{}
	adminHandler := &AdminHandler{}

	authHandler.SetNext(adminHandler)

	user1 := &User{isAuthenticated: false, isAdmin: false}
	fmt.Println("Processing user 1:")
	authHandler.Handle(user1)

	user2 := &User{isAuthenticated: true, isAdmin: false}
	fmt.Println("\nProcessing user 2:")
	authHandler.Handle(user2)

	user3 := &User{isAuthenticated: true, isAdmin: true}
	fmt.Println("\nProcessing user 3:")
	authHandler.Handle(user3)
}
