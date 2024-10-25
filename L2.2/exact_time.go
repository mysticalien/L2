package main

import (
	"fmt"
	"github.com/beevik/ntp"
	"log"
	"os"
)

/*
Task:
	Создать программу печатающую точное время с использованием NTP -библиотеки.
	Инициализировать как go module.
	Использовать библиотеку github.com/beevik/ntp.
	Написать программу печатающую текущее время / точное время с использованием этой библиотеки.
Требования:
	1. Программа должна быть оформлена как go module.
	2. Программа должна корректно обрабатывать ошибки библиотеки: выводить их в STDERR и возвращать ненулевой код выхода в OS.
	3. Создаем новый модуль: go mod init ntp_time
*/

func main() {
	address := "0.beevik-ntp.pool.ntp.org"
	time, err := ntp.Time(address)
	if err != nil {
		log.Println("Error getting the exact time:", err)
		os.Exit(1)
	}

	fmt.Println("Exact time:", time)
}
