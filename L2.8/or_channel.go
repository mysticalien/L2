package main

import (
	"fmt"
	"time"
)

/*
Task:
	Реализовать функцию, которая будет объединять один или более done-каналов в single-канал, если один из его составляющих каналов закроется.
	Очевидным вариантом решения могло бы стать выражение при использовании select, которое бы реализовывало эту связь, однако иногда неизвестно общее число done-каналов, с которыми вы работаете в рантайме.
	В этом случае удобнее использовать вызов единственной функции, которая, приняв на вход один или более or-каналов, реализовывала бы весь функционал.
	Определение функции:
		var or func(channels ...<- chan interface{}) <- chan interface{}
*/

func or(channels ...<-chan interface{}) <-chan interface{} {
	done := make(chan interface{})
	for _, ch := range channels {
		go func(c <-chan interface{}) {
			select {
			case <-c:
				close(done)
			case <-done:
			}
		}(ch)
	}
	return done
}

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("Done after %v\n", time.Since(start))
}
