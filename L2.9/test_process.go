package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	fmt.Println("Test process started. PID:", os.Getpid())
	for {
		time.Sleep(2 * time.Minute) // Удерживаем процесс в состоянии "живого"
	}
}
