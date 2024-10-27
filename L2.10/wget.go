package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

/*
Task:
	Реализовать утилиту wget с возможностью скачивать сайты целиком.
*/

func wget(url string, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to establish a connection to the site: %v", err)
	}
	defer resp.Body.Close()

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error when creating a file: %v", err)
	}
	defer file.Close()

	if _, err = io.Copy(file, resp.Body); err != nil {
		return fmt.Errorf("file write error: %v", err)
	}
	return nil
}

func main() {
	urlFlag := flag.String("url", "", "URL веб-страницы для скачивания")
	filenameFlag := flag.String("output", "output.html", "Имя файла для сохранения")

	flag.Parse()

	if *urlFlag == "" {
		fmt.Println("Usage: wget -url <URL> [-output <filename>]")
		return
	}

	err := wget(*urlFlag, *filenameFlag)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Downloaded: %s to file %s\n", *urlFlag, *filenameFlag)
}
