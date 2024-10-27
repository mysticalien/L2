package main

/*
Task:
	Реализовать HTTP-сервер для работы с календарем. В рамках задания необходимо работать строго со стандартной HTTP-библиотекой.
Требования:
	1. Реализовать вспомогательные функции для сериализации объектов доменной области в JSON.
	2. Реализовать вспомогательные функции для парсинга и валидации параметров методов /create_event и /update_event.
	3. Реализовать HTTP обработчики для каждого из методов API, используя вспомогательные функции и объекты доменной области.
	4. Реализовать middleware для логирования запросов
	Методы API:
		POST /create_event
		POST /update_event
		POST /delete_event
		GET /events_for_day
		GET /events_for_week
		GET /events_for_month
	Параметры передаются в виде www-url-form-encoded (т.е. обычные user_id=3&date=2019-09-09).
	В GET методах параметры передаются через queryString, в POST — через тело запроса.
	В результате каждого запроса должен возвращаться JSON-документ содержащий либо {"result": "..."} в случае успешного выполнения метода,
	либо {"error": "..."} в случае ошибки бизнес-логики.
В рамках задачи необходимо:
	1. Реализовать все методы.
	2. Бизнес логика НЕ должна зависеть от кода HTTP сервера.
	3. В случае ошибки бизнес-логики сервер должен возвращать HTTP 503. В случае ошибки входных данных (невалидный int например) сервер должен возвращать HTTP 400. В случае остальных ошибок сервер должен возвращать HTTP 500.
	4. Web-сервер должен запускаться на порту указанном в конфиге и выводить в лог каждый обработанный запрос.
*/

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port string `json:"port"`
}

type Event struct {
	ID     int       `json:"id"`
	UserID int       `json:"user_id"`
	Title  string    `json:"title"`
	Date   time.Time `json:"date"`
}

var (
	events      = make(map[int]Event)
	nextEventID = 1
	config      Config
)

func main() {
	// Загрузка конфигурации
	if err := loadConfig(); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Настройка маршрутов и запуск сервера
	http.HandleFunc("/create_event", logMiddleware(createEventHandler))
	http.HandleFunc("/update_event", logMiddleware(updateEventHandler))
	http.HandleFunc("/delete_event", logMiddleware(deleteEventHandler))
	http.HandleFunc("/events_for_day", logMiddleware(eventsForDayHandler))
	http.HandleFunc("/events_for_week", logMiddleware(eventsForWeekHandler))
	http.HandleFunc("/events_for_month", logMiddleware(eventsForMonthHandler))

	fmt.Printf("Server starting on %s\n", config.Port)
	log.Fatal(http.ListenAndServe(config.Port, nil))
}

// loadConfig загружает конфигурацию из config.json
func loadConfig() error {
	file, err := os.Open("config.json")
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return err
	}
	return nil
}

// logMiddleware - middleware для логирования запросов
func logMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	}
}

// createEventHandler обрабатывает создание события
func createEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	event, err := parseEvent(r, false)
	if err != nil {
		http.Error(w, `{"error":"Invalid input data"}`, http.StatusBadRequest)
		return
	}
	if err := addEvent(event); err != nil {
		http.Error(w, `{"error":"Service unavailable"}`, http.StatusServiceUnavailable)
		return
	}
	jsonResponse(w, map[string]string{"result": "event created"})
}

// updateEventHandler обрабатывает обновление события
func updateEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	event, err := parseEvent(r, true)
	if err != nil {
		http.Error(w, `{"error":"Invalid input data"}`, http.StatusBadRequest)
		return
	}
	if err := updateEvent(event); err != nil {
		if err.Error() == "event not found" {
			http.Error(w, `{"error":"Event not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error":"Service unavailable"}`, http.StatusServiceUnavailable)
		}
		return
	}
	jsonResponse(w, map[string]string{"result": "event updated"})
}

// deleteEventHandler обрабатывает удаление события
func deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	eventID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || eventID <= 0 {
		http.Error(w, `{"error":"Invalid event ID"}`, http.StatusBadRequest)
		return
	}
	if err := deleteEvent(eventID); err != nil {
		if err.Error() == "event not found" {
			http.Error(w, `{"error":"Event not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error":"Service unavailable"}`, http.StatusServiceUnavailable)
		}
		return
	}
	jsonResponse(w, map[string]string{"result": "event deleted"})
}

// eventsForDayHandler обрабатывает получение событий за день
func eventsForDayHandler(w http.ResponseWriter, r *http.Request) {
	date, err := parseDate(r, "date")
	if err != nil {
		http.Error(w, `{"error":"Invalid date"}`, http.StatusBadRequest)
		return
	}
	result := getEventsForDateRange(date, date)
	jsonResponse(w, result)
}

// eventsForWeekHandler обрабатывает получение событий за неделю
func eventsForWeekHandler(w http.ResponseWriter, r *http.Request) {
	date, err := parseDate(r, "date")
	if err != nil {
		http.Error(w, `{"error":"Invalid date"}`, http.StatusBadRequest)
		return
	}
	endDate := date.AddDate(0, 0, 7)
	result := getEventsForDateRange(date, endDate)
	jsonResponse(w, result)
}

// eventsForMonthHandler обрабатывает получение событий за месяц
func eventsForMonthHandler(w http.ResponseWriter, r *http.Request) {
	date, err := parseDate(r, "date")
	if err != nil {
		http.Error(w, `{"error":"Invalid date"}`, http.StatusBadRequest)
		return
	}
	endDate := date.AddDate(0, 1, 0)
	result := getEventsForDateRange(date, endDate)
	jsonResponse(w, result)
}

// Бизнес-логика

// addEvent добавляет событие в хранилище
func addEvent(event Event) error {
	event.ID = nextEventID
	nextEventID++
	events[event.ID] = event
	return nil
}

// updateEvent обновляет существующее событие
func updateEvent(event Event) error {
	if _, exists := events[event.ID]; !exists {
		return fmt.Errorf("event not found")
	}
	events[event.ID] = event
	return nil
}

// deleteEvent удаляет событие по ID
func deleteEvent(eventID int) error {
	if _, exists := events[eventID]; !exists {
		return fmt.Errorf("event not found")
	}
	delete(events, eventID)
	return nil
}

// getEventsForDateRange возвращает события в заданном диапазоне дат
func getEventsForDateRange(start, end time.Time) []Event {
	var result []Event
	for _, event := range events {
		if (event.Date.Equal(start) || event.Date.After(start)) && event.Date.Before(end.AddDate(0, 0, 1)) {
			result = append(result, event)
		}
	}
	return result
}

// Вспомогательные функции для парсинга и сериализации

// parseEvent разбирает данные события из запроса
func parseEvent(r *http.Request, requireID bool) (Event, error) {
	var event Event
	if err := r.ParseForm(); err != nil {
		return event, err
	}
	if requireID {
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return event, fmt.Errorf("invalid id")
		}
		event.ID = id
	}
	userID, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		return event, fmt.Errorf("invalid user_id")
	}
	event.UserID = userID
	event.Title = r.FormValue("title")
	event.Date, err = time.Parse("2006-01-02", r.FormValue("date"))
	if err != nil {
		return event, fmt.Errorf("invalid date format")
	}
	return event, nil
}

// parseDate извлекает дату из запроса
func parseDate(r *http.Request, param string) (time.Time, error) {
	dateStr := r.URL.Query().Get(param)
	return time.Parse("2006-01-02", dateStr)
}

// jsonResponse формирует JSON-ответ
func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
