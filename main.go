package main

import (
	"fmt"
	"os/exec"
	"time"
)

// Services directories
var serviceDirs = []string{
	"user_service",
	//"auth_service",
	"country_service",
	//"chatting_service",
}

func main() {
	var processes []*exec.Cmd

	// Build and launch microservices
	for _, dir := range serviceDirs {
		// Do go mod tidy for each service
		fmt.Printf("Executing `go mod tidy` for %s...\n", dir)
		tidyCmd := exec.Command("go", "mod", "tidy")
		tidyCmd.Dir = dir // Устанавливаем рабочую директорию
		if err := tidyCmd.Run(); err != nil {
			fmt.Printf("Error executing `go mod tidy` for %s: %v\n", dir, err)
			return
		}
		fmt.Printf("`go mod tidy` for %s executed\n", dir)

		// Path to the microservice source file
		sourcePath := "./cmd/main.go"
		// Path to the microservice binary file
		binaryPath := "./bin"

		// Сборка микросервиса
		fmt.Printf("Building service to %s...\n", dir)
		buildCmd := exec.Command("go", "build", "-o", binaryPath, sourcePath)
		buildCmd.Dir = dir // Устанавливаем рабочую директорию для сборки
		if err := buildCmd.Run(); err != nil {
			fmt.Printf("Error building %s: %v\n", dir, err)
			return
		}
		fmt.Printf("Microservice %s built\n", dir)

		// Запуск микросервиса
		runCmd := exec.Command(binaryPath + "/main.exe")
		runCmd.Dir = dir
		if err := runCmd.Start(); err != nil {
			fmt.Printf("Error starting service %s: %v\n", dir, err)
			stopProcesses(processes)
			return
		}
		processes = append(processes, runCmd)
		fmt.Printf("Service %s launched (PID: %d)\n", dir, runCmd.Process.Pid)

		// Даем время для инициализации (можно настроить)
		time.Sleep(2 * time.Second)

		// Проверяем, что процесс еще работает
		if runCmd.ProcessState != nil && runCmd.ProcessState.Exited() {
			fmt.Printf("Service %s ended with an error\n", dir)
			stopProcesses(processes)
			return
		}
	}

	fmt.Println("All services have been launched successfully.")

	// Ожидание завершения (для демонстрации)
	time.Sleep(10 * time.Second)
	stopProcesses(processes)
}

// Остановка всех запущенных процессов
func stopProcesses(processes []*exec.Cmd) {
	for _, cmd := range processes {
		if cmd.Process != nil {
			fmt.Printf("Останавливаем процесс %d\n", cmd.Process.Pid)
			if err := cmd.Process.Kill(); err != nil {
				fmt.Printf("Ошибка при остановке процесса %d: %v\n", cmd.Process.Pid, err)
			}
		}
	}
}
