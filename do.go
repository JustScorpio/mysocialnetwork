package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

type Microservice struct {
	Name     string
	Cmd      *exec.Cmd
	Running  bool
	Path     string
	BuildDir string
}

type MicroserviceManager struct {
	Services map[string]*Microservice
}

func NewMicroserviceManager() *MicroserviceManager {
	return &MicroserviceManager{
		Services: make(map[string]*Microservice),
	}
}

func (m *MicroserviceManager) FindMicroservices() ([]string, error) {
	var microservices []string

	entries, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Пропускаем исключенные директории
		switch entry.Name() {
		case "auth_service":
			continue
		case "chatting_service":
			continue
		}

		cmdPath := filepath.Join(entry.Name(), "cmd", "main.go")
		if _, err := os.Stat(cmdPath); err == nil {
			microservices = append(microservices, entry.Name())
		}
	}

	return microservices, nil
}

func (m *MicroserviceManager) BuildMicroservice(name string) error {
	servicePath := filepath.Join(".", name)
	binDir := filepath.Join(servicePath, "bin")
	exeName := "main.exe"
	if os.PathSeparator == '/' {
		exeName = "main"
	}
	exePath := filepath.Join(binDir, exeName)

	// Create bin directory if not exists
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %v", err)
	}

	// Run go mod init if needed
	if _, err := os.Stat(filepath.Join(servicePath, "go.mod")); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Printf("Running 'go mod init' for %s...\n", name)
			tidyCmd := exec.Command("go", "mod", "init", name)
			tidyCmd.Dir = servicePath
			tidyCmd.Stdout = os.Stdout
			tidyCmd.Stderr = os.Stderr
			if err := tidyCmd.Run(); err != nil {
				return fmt.Errorf("go mod init failed: %v", err)
			}
		} else {
			return fmt.Errorf("Cannot check if go.mod file exists or not. See details: %v", err)
		}
	}

	// Run go mod tidy
	fmt.Printf("Running 'go mod tidy' for %s...\n", name)
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = servicePath
	tidyCmd.Stdout = os.Stdout
	tidyCmd.Stderr = os.Stderr
	if err := tidyCmd.Run(); err != nil {
		return fmt.Errorf("go mod tidy failed: %v", err)
	}

	// Build microservice
	fmt.Printf("Building %s...\n", name)
	buildCmd := exec.Command("go", "build", "-o", exePath, ".")
	buildCmd.Dir = filepath.Join(servicePath, "cmd")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	fmt.Printf(buildCmd.Dir + " " + buildCmd.Path + "\n")
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("build failed: %v", err)
	}

	return nil
}

func (m *MicroserviceManager) StartMicroservice(name string) error {
	if service, exists := m.Services[name]; exists && service.Running {
		return fmt.Errorf("microservice %s is already running", name)
	}

	servicePath := filepath.Join(".", name)
	binDir := filepath.Join(servicePath, "bin")
	exeName := "main.exe"
	if os.PathSeparator == '/' {
		exeName = "main"
	}
	exePath := filepath.Join(binDir, exeName)

	// Build executable
	if err := m.BuildMicroservice(name); err != nil {
		return fmt.Errorf("build failed: %v", err)
	}

	fmt.Printf("Starting %s...\n", name)
	cmd := exec.Command(exePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start: %v", err)
	}

	m.Services[name] = &Microservice{
		Name:     name,
		Cmd:      cmd,
		Running:  true,
		Path:     servicePath,
		BuildDir: binDir,
	}

	go func() {
		err := cmd.Wait()
		if service, exists := m.Services[name]; exists {
			service.Running = false
			if err != nil {
				log.Printf("Microservice %s exited with error: %v\n", name, err)
			} else {
				log.Printf("Microservice %s exited\n", name)
			}
		}
	}()

	return nil
}

func (m *MicroserviceManager) StopMicroservice(name string) error {
	service, exists := m.Services[name]
	if !exists || !service.Running {
		return fmt.Errorf("microservice %s is not running", name)
	}

	fmt.Printf("Stopping %s...\n", name)
	if err := service.Cmd.Process.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("failed to send SIGTERM: %v", err)
	}

	// Wait for process to exit with timeout
	done := make(chan error, 1)
	go func() {
		done <- service.Cmd.Wait()
	}()

	select {
	case <-time.After(5 * time.Second):
		if err := service.Cmd.Process.Kill(); err != nil {
			return fmt.Errorf("failed to kill process: %v", err)
		}
		<-done // Wait for kill to complete
	case err := <-done:
		if err != nil {
			return fmt.Errorf("process exited with error: %v", err)
		}
	}

	delete(m.Services, name)
	return nil
}

func (m *MicroserviceManager) StartAll() error {
	microservices, err := m.FindMicroservices()
	if err != nil {
		return fmt.Errorf("failed to find microservices: %v", err)
	}

	if len(microservices) == 0 {
		return fmt.Errorf("no microservices found")
	}

	var lastError error
	for _, name := range microservices {
		if err := m.StartMicroservice(name); err != nil {
			log.Printf("Error starting %s: %v\n", name, err)
			lastError = err
		}
	}

	return lastError
}

func (m *MicroserviceManager) StopAll() error {
	if len(m.Services) == 0 {
		return fmt.Errorf("no microservices are running")
	}

	var lastError error
	for name := range m.Services {
		if err := m.StopMicroservice(name); err != nil {
			log.Printf("Error stopping %s: %v\n", name, err)
			lastError = err
		}
	}

	return lastError
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  do all up              - Start all microservices")
	fmt.Println("  do all down            - Stop all microservices")
	fmt.Println("  do <microservice> up   - Start a specific microservice")
	fmt.Println("  do <microservice> down - Stop a specific microservice")
	fmt.Println("\nAvailable microservices:")

	manager := NewMicroserviceManager()
	microservices, err := manager.FindMicroservices()
	if err != nil {
		log.Fatalf("Failed to find microservices: %v", err)
	}

	for _, ms := range microservices {
		fmt.Printf("  - %s\n", ms)
	}
}

func main() {
	if len(os.Args) != 3 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	action := strings.ToLower(os.Args[2])

	manager := NewMicroserviceManager()

	switch command {
	case "all":
		switch action {
		case "up":
			if err := manager.StartAll(); err != nil {
				log.Fatalf("Failed to start all microservices: %v", err)
			}
		case "down":
			if err := manager.StopAll(); err != nil {
				log.Fatalf("Failed to stop all microservices: %v", err)
			}
		default:
			fmt.Printf("Unknown action: %s\n", action)
			printUsage()
			os.Exit(1)
		}
	default:
		// Check if microservice exists
		microservices, err := manager.FindMicroservices()
		if err != nil {
			log.Fatalf("Failed to find microservices: %v", err)
		}

		found := false
		for _, ms := range microservices {
			if ms == command {
				found = true
				break
			}
		}

		if !found {
			fmt.Printf("Unknown microservice: %s\n", command)
			fmt.Println("Available microservices:")
			for _, ms := range microservices {
				fmt.Printf("  - %s\n", ms)
			}
			os.Exit(1)
		}

		switch action {
		case "up":
			if err := manager.StartMicroservice(command); err != nil {
				log.Fatalf("Failed to start microservice %s: %v", command, err)
			}
		case "down":
			if err := manager.StopMicroservice(command); err != nil {
				log.Fatalf("Failed to stop microservice %s: %v", command, err)
			}
		default:
			fmt.Printf("Unknown action: %s\n", action)
			printUsage()
			os.Exit(1)
		}
	}
}
