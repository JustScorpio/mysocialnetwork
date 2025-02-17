package main

import (
	"auth_service/internal/handler"
	"auth_service/internal/service"
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

var store = sessions.NewCookieStore([]byte("your-secret-key"))

func main() {
	// Подключение к базе данных
	db, err := sql.Open("postgres", "user=youruser dbname=yourdb sslmode=disable password=yourpassword")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Инициализация сервиса и хендлера
	authService := service.NewAuthService(db)
	authHandler := handler.NewAuthHandler(authService)

	// Роутинг
	http.HandleFunc("/register", authHandler.Register)
	http.HandleFunc("/login", authHandler.Login)
	http.HandleFunc("/check-auth", authHandler.CheckAuth)
	http.HandleFunc("/logout", authHandler.Logout)

	// Запуск сервера
	log.Println("Auth Service is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
