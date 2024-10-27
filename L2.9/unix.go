package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

/*
Task:
	Необходимо реализовать свой собственный UNIX-шелл-утилиту с поддержкой ряда простейших команд:
		- cd <args> — смена директории (в качестве аргумента могут быть то-то и то);
		- pwd — показать путь до текущего каталога;
		- echo <args> — вывод аргумента в STDOUT;
		- kill <args> — «убить» процесс, переданный в качесте аргумента (пример: такой-то пример);
		- ps — выводит общую информацию по запущенным процессам в формате такой-то формат.
	Так же требуется поддерживать функционал fork/exec-команд.
	Дополнительно необходимо поддерживать конвейер на пайпах (linux pipes, пример cmd1 | cmd2 | .... | cmdN).
	*Шелл — это обычная консольная программа, которая будучи запущенной, в интерактивном сеансе выводит некое приглашение в STDOUT и ожидает ввода пользователя через STDIN.
	Дождавшись ввода, обрабатывает команду согласно своей логике и при необходимости выводит результат на экран.
	Интерактивный сеанс поддерживается до тех пор, пока не будет введена команда выхода (например \quit).
*/

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "\\quit" {
			break
		}
		executeCommand(input)
	}
}

func executeCommand(input string) {
	// Проверка на пустой ввод
	if input == "" {
		fmt.Println("Error: empty command")
		return
	}

	// Обработка пайпов
	if strings.Contains(input, "|") {
		commands := strings.Split(input, "|")
		var cmdList []*exec.Cmd

		for _, command := range commands {
			cmdArgs := parseArguments(command)
			if len(cmdArgs) > 0 {
				cmdList = append(cmdList, exec.Command(cmdArgs[0], cmdArgs[1:]...))
			}
		}

		for i := 0; i < len(cmdList)-1; i++ {
			pipe, err := cmdList[i].StdoutPipe()
			if err != nil {
				fmt.Println("Error creating pipe:", err)
				return
			}
			if err := cmdList[i].Start(); err != nil {
				fmt.Println("Error starting command:", err)
				return
			}
			cmdList[i+1].Stdin = pipe
		}

		// Для последней команды в цепочке
		cmdList[len(cmdList)-1].Stdout = os.Stdout
		if err := cmdList[len(cmdList)-1].Run(); err != nil {
			fmt.Println("Error running command:", err)
		}
		return
	}

	// Логика для одиночных команд
	args := parseArguments(input)
	if len(args) == 0 {
		return
	}

	switch args[0] {
	case "cd":
		if len(args) < 2 {
			fmt.Println("cd: missing argument")
			return
		}
		if err := os.Chdir(args[1]); err != nil {
			fmt.Println("cd error:", err)
		}
	case "pwd":
		if dir, err := os.Getwd(); err == nil {
			fmt.Println(dir)
		} else {
			fmt.Println("pwd error:", err)
		}
	case "echo":
		fmt.Println(strings.Join(args[1:], " "))
	case "kill":
		if len(args) < 2 {
			fmt.Println("kill: missing PID")
			return
		}
		pid, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("kill: invalid PID")
			return
		}
		if err := syscall.Kill(pid, syscall.SIGKILL); err != nil {
			fmt.Println("kill error:", err)
		}
	case "ps":
		cmd := exec.Command("ps")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	default:
		runExternalCommand(args)
	}
}

func parseArguments(input string) []string {
	var args []string
	var currentArg strings.Builder
	inQuotes := false

	for _, char := range input {
		switch char {
		case '"':
			inQuotes = !inQuotes
		case ' ':
			if inQuotes {
				currentArg.WriteRune(char)
			} else {
				if currentArg.Len() > 0 {
					args = append(args, currentArg.String())
					currentArg.Reset()
				}
			}
		default:
			currentArg.WriteRune(char)
		}
	}
	if currentArg.Len() > 0 {
		args = append(args, currentArg.String())
	}

	return args
}

func runExternalCommand(args []string) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Command execution error:", err)
	}
}
