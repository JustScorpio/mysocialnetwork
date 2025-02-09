package main

import (
	"log"
	"net/http"
	"user_service/internal/database"
	"user_service/internal/handlers"
	"user_service/internal/repository"
	"user_service/internal/services"
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
	userCrudService := services.NewCrudService(userRepo)

	// Инициализация обработчиков
	userHandler := handlers.NewCrudHandler(userCrudService)

	// Регистрация маршрутов
	mux := http.NewServeMux()
	mux.HandleFunc("/countries", userHandler.GetAll)
	mux.HandleFunc("/countries/get", userHandler.Get)
	mux.HandleFunc("/countries/create", userHandler.Create)
	mux.HandleFunc("/countries/update", userHandler.Update)
	mux.HandleFunc("/countries/delete", userHandler.Delete)

	// Запуск сервера
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
