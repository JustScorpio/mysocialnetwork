package handlers

import (
	"encoding/json"
	"net/http"
	"network/internal/models"
	"network/internal/services"
	"strconv"
)

type CountryHandler struct {
	countryService *services.Service[models.Country]
}

func NewCountryHandler(countryService *services.Service[models.Country]) *CountryHandler {
	return &CountryHandler{countryService: countryService}
}

func (h *CountryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	countries, err := h.countryService.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(countries)
}

func (h *CountryHandler) Get(w http.ResponseWriter, r *http.Request) {
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

	// Получение страны из сервиса
	country, err := h.countryService.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возврат страны в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(country)
}

func (h *CountryHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Декодирование тела запроса в структуру Country
	var country models.Country
	if err := json.NewDecoder(r.Body).Decode(&country); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Создание страны через сервис
	if err := h.countryService.Create(&country); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возврат успешного статуса
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":   country.Id,
		"name": country.Name,
	})
}

func (h *CountryHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Декодирование тела запроса в структуру Country
	var country models.Country
	if err := json.NewDecoder(r.Body).Decode(&country); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Обновление страны через сервис
	if err := h.countryService.Update(&country); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возврат успешного статуса
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Country updated successfully",
	})
}

func (h *CountryHandler) DeleteCountry(w http.ResponseWriter, r *http.Request) {
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
	if err := h.countryService.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возврат успешного статуса
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Country deleted successfully",
	})
}
