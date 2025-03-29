package main

import (
	"log"
	"mysocialnetwork/user_service/internal/database"
	"mysocialnetwork/user_service/internal/handlers"
	"mysocialnetwork/user_service/internal/repository"
	"mysocialnetwork/user_service/internal/services"
	"net/http"
)

func main() {
	// Инициализация базы данных
	db, err := database.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Инициализация репозиториев
	countryCacheRepository := repository.NewCountryRepository(db)
	userRepo := repository.NewUserRepository(db, countryCacheRepository)

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
	log.Println("Server started on :8081")
	log.Fatal(http.ListenAndServe(":8081", mux))
}
