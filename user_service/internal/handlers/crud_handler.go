package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"user_service/internal/models"
	"user_service/internal/services"
)

type CrudHandler[T models.Entity] struct {
	service *services.CrudService[T]
}

func NewCrudHandler[T models.Entity](service *services.CrudService[T]) *CrudHandler[T] {
	return &CrudHandler[T]{service: service}
}

func (h *CrudHandler[T]) GetAll(w http.ResponseWriter, r *http.Request) {
	entities, err := h.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entities)
}

func (h *CrudHandler[T]) Get(w http.ResponseWriter, r *http.Request) {
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

	// Получение сущности из сервиса
	entity, err := h.service.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возврат в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entity)
}

func (h *CrudHandler[T]) Create(w http.ResponseWriter, r *http.Request) {
	// Декодирование тела запроса в структуру T
	var entity T
	if err := json.NewDecoder(r.Body).Decode(&entity); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Создание сущности через сервис
	if err := h.service.Create(&entity); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возврат успешного статуса
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entity)
}

func (h *CrudHandler[T]) Update(w http.ResponseWriter, r *http.Request) {
	// Декодирование тела запроса в структуру T
	var entity T
	if err := json.NewDecoder(r.Body).Decode(&entity); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Обновление страны через сервис
	if err := h.service.Update(&entity); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возврат успешного статуса
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Entity updated successfully",
	})
}

func (h *CrudHandler[T]) Delete(w http.ResponseWriter, r *http.Request) {
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
	if err := h.service.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возврат успешного статуса
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Entity deleted successfully",
	})
}
