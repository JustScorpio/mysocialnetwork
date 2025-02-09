package main

import (
	"log"
	"net/http"
	"user-service/internal/database"
	"user-service/internal/handlers"
	"user-service/internal/repository"
	"user-service/internal/services"
)

func main() {
	// Инициализация базы данных
	db, err := database.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Инициализация репозиториев
	userRepo := repository.NewUserRepository(db)

	// Инициализация сервисов
	userCrudService := services.NewService(userRepo)

	// Инициализация обработчиков
	userHandler := handlers.NewUserHandler(userCrudService)

	// Регистрация маршрутов
	mux := http.NewServeMux()
	mux.HandleFunc("/countries", userHandler.GetAll)
	mux.HandleFunc("/countries/get", userHandler.Get)
	mux.HandleFunc("/countries/create", userHandler.Create)
	mux.HandleFunc("/countries/update", userHandler.Update)
	mux.HandleFunc("/countries/delete", userHandler.DeleteUser)

	// Запуск сервера
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
