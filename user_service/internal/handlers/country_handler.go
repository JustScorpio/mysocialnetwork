package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"user-service/internal/models"
	"user-service/internal/services"
)

type UserHandler struct {
	userService *services.Service[models.User]
}

func NewUserHandler(userService *services.Service[models.User]) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	countries, err := h.userService.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(countries)
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	// Извлечение Id из URL-параметра
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Id is required", http.StatusBadRequest)
		return
	}

	// Преобразование ID в число
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Id", http.StatusBadRequest)
		return
	}

	// Получение пользователя из сервиса
	user, err := h.userService.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возврат пользователя в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Декодирование тела запроса в структуру User
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Создание страны через сервис
	if err := h.userService.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возврат успешного статуса
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":   user.Id,
		"name": user.Name,
		"user": user.Mail,
	})
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Декодирование тела запроса в структуру User
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Обновление страны через сервис
	if err := h.userService.Update(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возврат успешного статуса
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User updated successfully",
	})
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Извлечение ID из URL-параметра
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	// Преобразование ID в число
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Удаление страны через сервис
	if err := h.userService.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возврат успешного статуса
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User deleted successfully",
	})
}
