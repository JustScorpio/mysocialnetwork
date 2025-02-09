package main

import (
	"country_service/internal/database"
	"country_service/internal/handlers"
	"country_service/internal/repository"
	"country_service/internal/services"
	"log"
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
	countryRepo := repository.NewCountryRepository(db)

	// Инициализация сервисов
	countryCrudService := services.NewService(countryRepo)

	// Инициализация обработчиков
	countryHandler := handlers.NewCountryHandler(countryCrudService)

	// Регистрация маршрутов
	mux := http.NewServeMux()
	mux.HandleFunc("/countries", countryHandler.GetAll)
	mux.HandleFunc("/countries/get", countryHandler.Get)
	mux.HandleFunc("/countries/create", countryHandler.Create)
	mux.HandleFunc("/countries/update", countryHandler.Update)
	mux.HandleFunc("/countries/delete", countryHandler.DeleteCountry)

	// Запуск сервера
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
